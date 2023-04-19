package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/anonymousMoonPrince/file-service/internal/app/response"
)

type FileController struct {
	fileService FileService
}

func NewFileController(fileService FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

func (c *FileController) Upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		response.BadRequest(w, err)
		return
	}

	fileID, err := c.fileService.Upload(r.Context(), header.Filename, r.Header.Get("Content-Type"), file, header.Size)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Ok(w, UploadResponse{FileID: fileID})
}

//easyjson:json
type UploadResponse struct {
	FileID string `json:"file_id"`
}

func (c *FileController) Download(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "uuid")
	if fileID == "" {
		response.BadRequest(w, errors.New("fileID is not provided"))
		return
	}

	if err := c.fileService.Download(r.Context(), w, fileID); err != nil {
		response.InternalServerError(w, err)
	}
}
