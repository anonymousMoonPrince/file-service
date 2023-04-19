package controller

import (
	"context"
	"io"
	"net/http"
)

type FileService interface {
	Upload(ctx context.Context, name, contentType string, file io.Reader, size int64) (string, error)
	Download(ctx context.Context, w http.ResponseWriter, id string) error
}
