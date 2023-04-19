package storage

import (
	"context"
	"io"
)

type Client interface {
	MustCreateBucket(ctx context.Context, bucket string)
	GetURLs(bucket, objectName string, count int) []string
	Upload(ctx context.Context, url, bucket, objectName string, data io.Reader, size int64, contentType string) (string, string, error)
	Get(ctx context.Context, url, bucket, objectName string) (io.ReadCloser, int64, error)
	Check(ctx context.Context, url, bucket, objectName string, size int64, etag string) (bool, error)
}
