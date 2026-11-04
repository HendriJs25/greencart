package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type localStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) Storage {
	return &localStorage{
		basePath: basePath,
	}
}

func (s *localStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, path string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	fullPath := filepath.Join(s.basePath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create dir: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	_, copyErr := dst.ReadFrom(src)
	closeErr := dst.Close()

	if copyErr != nil {
		_ = os.Remove(fullPath)
		return "", fmt.Errorf("failed to copy file: %w", copyErr)
	}

	if closeErr != nil {
		_ = os.Remove(fullPath)
		return "", fmt.Errorf("failed to close file: %w", closeErr)
	}

	if err := ctx.Err(); err != nil {
		_ = os.Remove(fullPath)
		return "", err
	}

	return
}

func (s *localStorage) DeleteFile(ctx context.Context, path string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}
