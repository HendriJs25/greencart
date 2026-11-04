package storage

import (
	"context"
	"mime/multipart"
)

type Storage interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, path string) (string, error)
	DeleteFile(ctx context.Context, path string) error
}
