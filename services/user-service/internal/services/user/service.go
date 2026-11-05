package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/logger"
	userrepository "user-service/internal/repository/user"
	"user-service/internal/storage"
	"uuid"
)

type service struct {
	userRepository userrepository.Repository
	s3Storage      storage.Storage
}

type Service interface {
	GetProfile(ctx context.Context, userUUID uuid.UUID) (*GetProfileResult, error)
	UploadProfileImage(ctx context.Context, userUUID uuid.UUID, file *multipart.FileHeader) (*string, error)
}

func NewService(userRepository userrepository.Repository, s3Storage storage.Storage) Service {
	return &service{
		userRepository: userRepository,
		s3Storage:      s3Storage,
	}
}

func (s *service) GetProfile(ctx context.Context, userUUID uuid.UUID) (*GetProfileResult, error) {
	user, err := s.userRepository.FindByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	var profileImageURL *string
	if user.Photo != nil {
		profileImageURL = s.s3Storage.GetFileURL(*user.Photo)
	}

	result := ToGetProfileResult(user, profileImageURL)
	return result, nil
}

func (s *service) UploadProfileImage(ctx context.Context, userUUID uuid.UUID, file *multipart.FileHeader) (*string, error) {
	contentType, err := validateProfileImage(file)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	oldKey := user.Photo

	path := fmt.Sprintf("users/%s/profile/%s%s", user.UUID, uuid.NewV4().String(), strings.ToLower(filepath.Ext(file.Filename)))

	uploadedKey, err := s.s3Storage.UploadFile(ctx, file, path, contentType)
	if err != nil {
		return nil, err
	}

	if err := s.userRepository.UpdatePhotoKeyByUUID(ctx, user.UUID, uploadedKey); err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		_ = s.s3Storage.DeleteFile(cleanupCtx, uploadedKey)
		cancel()
		return nil, err
	}

	if oldKey != nil && *oldKey != "" {
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		if err := s.s3Storage.DeleteFile(cleanupCtx, *oldKey); err != nil {
			logger.WithContext(ctx).WithError(err).Error("failed to delete old file")
		}
		cancel()
	}

	return s.s3Storage.GetFileURL(uploadedKey), nil
}

func validateProfileImage(file *multipart.FileHeader) (string, error) {
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
