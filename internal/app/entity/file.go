package entity

import "time"

const TableNameFiles = "files"

type File struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	ContentType string    `db:"content_type"`
	IsUploaded  bool      `db:"is_uploaded"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
