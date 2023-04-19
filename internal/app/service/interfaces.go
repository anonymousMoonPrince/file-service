package service

import (
	"context"

	"github.com/anonymousMoonPrince/file-service/internal/app/entity"
)

type FileRepository interface {
	Create(ctx context.Context, name, contentType string) (string, error)
	CreateChunk(ctx context.Context, chunk entity.Chunk) error
	MarkFileAsUploaded(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.File, error)
	GetChunksByFileID(ctx context.Context, fileID string) ([]entity.Chunk, error)
}
