package minio

import (
	"context"
	"io"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/serialx/hashring"
	"github.com/sirupsen/logrus"
)

type Client struct {
	ptr         map[string]*minio.Core
	weightByURL map[string]int
	ring        *hashring.HashRing
	mux         sync.RWMutex
}

const MB = 1 << 20

func NewClient(configs []Config, sizeByURL map[string]int) *Client {
	weightByURL := getWeightByURL(configs, sizeByURL)
	return &Client{
		ptr:         getPtr(configs),
		ring:        hashring.NewWithWeights(weightByURL),
		weightByURL: weightByURL,
	}
}

func (c *Client) UpdateClient(configs []Config, sizeByURL map[string]int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	weightByURL := getWeightByURL(configs, sizeByURL)
	c.ptr = getPtr(configs)
	c.ring = hashring.NewWithWeights(weightByURL)
	c.weightByURL = weightByURL
}

func getWeightByURL(configs []Config, sizeByURL map[string]int) map[string]int {
	weightByURL := make(map[string]int)
	for _, config := range configs {
		weight := 1
		if size := sizeByURL[config.URL]; size > 0 {
			weight += size / MB
		}
		weightByURL[config.URL] = weight
	}
	return weightByURL
}

func getPtr(configs []Config) map[string]*minio.Core {
	ptr := make(map[string]*minio.Core)
	for _, config := range configs {
		minioClient, err := minio.NewCore(config.URL, &minio.Options{
			Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, config.Token),
			Secure: config.Secure,
			Region: config.Region,
		})
		if err != nil {
			logrus.WithError(err).Fatal("create minio client failed")
		}

		ptr[config.URL] = minioClient
		ptr[minioClient.EndpointURL().String()] = minioClient
	}
	return ptr
}

func (c *Client) MustCreateBucket(ctx context.Context, bucket string) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	for _, ptr := range c.ptr {
		isExist, err := ptr.BucketExists(ctx, bucket)
		if err != nil {
			logrus.WithError(err).Fatal("check bucket failed")
		}
		if !isExist {
			if err = ptr.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				logrus.WithError(err).Fatal("create bucket failed")
			}
		}
	}
}

func (c *Client) GetURLs(bucket, objectName string, count int) []string {
	c.mux.RLock()
	defer c.mux.RUnlock()

	urls, _ := c.ring.GetNodes(bucket+"/"+objectName, count)
	return urls
}

func (c *Client) Upload(ctx context.Context, url, bucket, objectName string, data io.Reader, size int64, contentType string) (string, string, error) {
	ptr := c.getPtrByURL(url)

	info, err := ptr.PutObject(ctx, bucket, objectName, data, size, "", "", minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", "", err
	}

	c.updateWeight(ptr, size)
	return ptr.EndpointURL().String(), info.ETag, nil
}

func (c *Client) Get(ctx context.Context, url, bucket, objectName string) (io.ReadCloser, int64, error) {
	var opts minio.GetObjectOptions

	reader, info, _, err := c.getPtrByURL(url).GetObject(ctx, bucket, objectName, opts)
	if err != nil {
		return nil, 0, err
	}
	return reader, info.Size, nil
}

func (c *Client) Check(ctx context.Context, url, bucket, objectName string, size int64, etag string) (bool, error) {
	info, err := c.getPtrByURL(url).StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return info.Size == size && info.ETag == etag, nil
}

func (c *Client) getPtrByURL(url string) *minio.Core {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.ptr[url]
}

func (c *Client) updateWeight(ptr *minio.Core, size int64) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	url := ptr.EndpointURL().String()

	ptrSize := c.weightByURL[url]
	ptrSize += int(size)
	c.weightByURL[url] = ptrSize

	c.ring.UpdateWeightedNode(url, ptrSize)
}
