package helper

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	errConstant "product-service/internal/constants/error"
	"strings"
)

func ValidateImage(file *multipart.FileHeader) (string, error) {
	const maxSize = 2 << 20 // 2 MiB

	if file.Size > maxSize {
		return "", errConstant.ErrFileTooBig
	}

	allowedTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	expectedType, ok := allowedTypes[ext]
	if !ok {
		return "", errConstant.ErrUnsupportedFileExtension
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	header := make([]byte, 512)

	n, err := io.ReadFull(src, header)
	if err != nil && errors.Is(err, io.ErrUnexpectedEOF) {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	contentType := http.DetectContentType(header[:n])

	if contentType != expectedType {
		return "", errConstant.ErrUnsupportedFileExtension
	}

	return contentType, nil
}
