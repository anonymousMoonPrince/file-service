package entity

import "time"

const TableNameChunks = "chunks"

type Chunk struct {
	ID        string    `db:"id"`
	FileID    string    `db:"file_id"`
	Number    int       `db:"number"`
	Size      int64     `db:"size"`
	CheckSum  string    `db:"check_sum"`
	URL       string    `db:"url"`
	Bucket    string    `db:"bucket"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
