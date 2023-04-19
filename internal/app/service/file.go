package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/anonymousMoonPrince/file-service/internal/app/client/storage"
	"github.com/anonymousMoonPrince/file-service/internal/app/config"
	"github.com/anonymousMoonPrince/file-service/internal/app/entity"
	"github.com/anonymousMoonPrince/file-service/internal/app/utils"
)

type FileService struct {
	storageClient storage.Client
	repo          FileRepository
	bucket        string
}

func NewFileService(
	storageClient storage.Client,
	repo FileRepository,
	bucket string,
) *FileService {
	storageClient.MustCreateBucket(context.Background(), bucket)
	return &FileService{
		storageClient: storageClient,
		repo:          repo,
		bucket:        bucket,
	}
}

func (s *FileService) Upload(ctx context.Context, name, contentType string, file io.Reader, size int64) (string, error) {
	fileID, err := s.repo.Create(ctx, name, contentType)
	if err != nil {
		return "", err
	}

	urls := s.storageClient.GetURLs(s.bucket, fileID, config.Get().BusinessConfig.ChunkCount)

	chunks := utils.Split(size, len(urls))
	for i, chunk := range chunks {
		url, etag, err := s.storageClient.Upload(ctx, urls[i], s.bucket, fileID, io.LimitReader(file, chunk), chunk, contentType)
		if err != nil {
			return "", err
		}

		if err = s.repo.CreateChunk(ctx, entity.Chunk{
			FileID:   fileID,
			Number:   i + 1,
			Size:     chunk,
			CheckSum: etag,
			URL:      url,
			Bucket:   s.bucket,
		}); err != nil {
			return "", err
		}
	}

	if err = s.repo.MarkFileAsUploaded(ctx, fileID); err != nil {
		return "", err
	}

	return fileID, nil
}

func (s *FileService) Download(ctx context.Context, w http.ResponseWriter, id string) error {
	file, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if file == nil || !file.IsUploaded {
		return errors.New("file is not uploaded")
	}

	chunks, err := s.repo.GetChunksByFileID(ctx, id)
	if err != nil {
		return err
	}

	var totalSize int64
	for _, chunk := range chunks {
		chunkIsOk, err := s.storageClient.Check(ctx, chunk.URL, chunk.Bucket, chunk.FileID, chunk.Size, chunk.CheckSum)
		if err != nil {
			return err
		}

		if !chunkIsOk {
			return errors.New("some chunk is broken")
		}

		totalSize += chunk.Size
	}

	w.Header().Add("Content-Type", file.ContentType)
	w.Header().Add("Content-Length", fmt.Sprintf("%d", totalSize))
	for _, chunk := range chunks {
		if err = s.uploadChunk(ctx, w, chunk); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileService) uploadChunk(ctx context.Context, w http.ResponseWriter, chunk entity.Chunk) error {
	reader, size, err := s.storageClient.Get(ctx, chunk.URL, chunk.Bucket, chunk.FileID)
	if err != nil {
		return err
	}

	if _, err = io.CopyN(w, reader, size); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}
