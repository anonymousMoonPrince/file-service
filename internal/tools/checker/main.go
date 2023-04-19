package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type response struct {
	Data struct {
		FileID string `json:"file_id"`
	} `json:"data"`
}

var filepath = flag.String("f", "internal/tools/checker/file.txt", "file for check uploading/downloading")

func init() {
	flag.Parse()
}

func main() {
	logrus.Info("start checker")

	file, err := os.Open(*filepath)
	if err != nil {
		logrus.WithError(err).Fatal("open check file failed")
	}
	defer file.Close()
	logrus.Info("open file")

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile("file", "file")
	if err != nil {
		logrus.WithError(err).Fatal("create form file failed")
	}

	if _, err = io.Copy(formFile, file); err != nil {
		logrus.WithError(err).Fatal("read file failed")
	}
	logrus.Info("read file")

	if err = writer.Close(); err != nil {
		logrus.WithError(err).Fatal("close form file failed")
	}
	logrus.Info("create form file")

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Fatal("upload file failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		logrus.WithError(err).WithField("status", resp.Status).Fatal("wrong status code")
	}

	response := new(response)
	if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
		logrus.WithError(err).Fatal("unmarshal response failed")
	}
	logrus.WithField("file_id", response.Data.FileID).Info("upload file")

	downloadResp, err := http.Get("http://localhost:8080/" + response.Data.FileID)
	if err != nil {
		logrus.WithError(err).Fatal("download file failed")
	}
	defer downloadResp.Body.Close()
	logrus.Info("download file")

	downloadedFile, err := os.Create(response.Data.FileID + ".output")
	if err != nil {
		logrus.WithError(err).WithField("file_id", response.Data.FileID).Fatal("create downloaded file failed")
	}
	defer downloadedFile.Close()
	logrus.Info("create downloaded file")

	size, err := io.Copy(downloadedFile, downloadResp.Body)
	if err != nil {
		logrus.WithError(err).Fatal("write downloaded file failed")
	}
	logrus.Info("write downloaded file")

	info, _ := file.Stat()
	if info.Size() != size {
		logrus.WithError(err).WithFields(logrus.Fields{
			"file_size":            info.Size(),
			"downloaded_file_size": size,
		}).Fatal("file size mismatched")
	}

	logrus.Info("finish checker")
}
