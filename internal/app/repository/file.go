package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/anonymousMoonPrince/file-service/internal/app/client/database"
	"github.com/anonymousMoonPrince/file-service/internal/app/entity"
)

type FileRepository struct {
	client  database.Client
	builder sq.StatementBuilderType
}

func NewFileRepository(client database.Client) *FileRepository {
	return &FileRepository{
		client:  client,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *FileRepository) Create(ctx context.Context, name, contentType string) (string, error) {
	query, args, err := r.builder.
		Insert(entity.TableNameFiles).
		Columns("name", "content_type").
		Values(name, contentType).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("build create file query failed: %w", err)
	}

	var id string
	if err = r.client.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return "", fmt.Errorf("create file failed: %w", err)
	}

	return id, nil
}

func (r *FileRepository) CreateChunk(ctx context.Context, chunk entity.Chunk) error {
	query, args, err := r.builder.
		Insert(entity.TableNameChunks).
		Columns("file_id", "number", "size", "check_sum", "url", "bucket").
		Values(chunk.FileID, chunk.Number, chunk.Size, chunk.CheckSum, chunk.URL, chunk.Bucket).
		ToSql()
	if err != nil {
		return fmt.Errorf("build create chunk query failed: %w", err)
	}

	if _, err = r.client.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("create chunk failed: %w", err)
	}

	return nil
}

func (r *FileRepository) MarkFileAsUploaded(ctx context.Context, id string) error {
	query, args, err := r.builder.
		Update(entity.TableNameFiles).
		SetMap(map[string]interface{}{
			"is_uploaded": true,
			"updated_at":  time.Now(),
		}).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build mark file as uploaded query failed: %w", err)
	}

	if _, err = r.client.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("mark file as uploaded failed: %w", err)
	}

	return nil
}

func (r *FileRepository) GetByID(ctx context.Context, id string) (*entity.File, error) {
	query, args, err := r.builder.
		Select("id", "name", "content_type", "is_uploaded", "created_at", "updated_at").
		From(entity.TableNameFiles).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get file query failed: %w", err)
	}

	var file entity.File
	if err := r.client.QueryRow(ctx, query, args...).
		Scan(&file.ID, &file.Name, &file.ContentType, &file.IsUploaded, &file.CreatedAt, &file.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get file failed: %w", err)
	}
	return &file, nil
}

func (r *FileRepository) GetChunksByFileID(ctx context.Context, fileID string) ([]entity.Chunk, error) {
	query, args, err := r.builder.
		Select("id", "file_id", "number", "size", "check_sum", "url", "bucket", "created_at", "updated_at").
		From(entity.TableNameChunks).
		Where(sq.Eq{"file_id": fileID}).
		OrderBy("number").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get chunks query failed: %w", err)
	}

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get chunks failed: %w", err)
	}
	defer rows.Close()

	var chunks []entity.Chunk
	for rows.Next() {
		var chunk entity.Chunk
		if err := rows.Scan(&chunk.ID, &chunk.FileID, &chunk.Number, &chunk.Size, &chunk.CheckSum, &chunk.URL, &chunk.Bucket, &chunk.CreatedAt, &chunk.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan chunks failed: %w", err)
		}

		chunks = append(chunks, chunk)
	}
	return chunks, nil
}

func (r *FileRepository) GetSizeByURL(ctx context.Context) (map[string]int, error) {
	query, args, err := r.builder.
		Select("sum(size)", "url").
		From(entity.TableNameChunks).
		GroupBy("url").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get size by url query failed: %w", err)
	}

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get size by url failed: %w", err)
	}
	defer rows.Close()

	sizeByURL := make(map[string]int)
	for rows.Next() {
		var (
			size int
			url  string
		)
		if err := rows.Scan(&size, &url); err != nil {
			return nil, fmt.Errorf("scan size by url failed: %w", err)
		}

		sizeByURL[url] = size
	}
	return sizeByURL, nil
}
