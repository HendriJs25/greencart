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
	"user-service/internal/common/helper"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/logger"
	sessionrepository "user-service/internal/repository/session"
	userrepository "user-service/internal/repository/user"
	"user-service/internal/storage"
	"uuid"
)

type service struct {
	userRepository    userrepository.Repository
	sessionRepository sessionrepository.Repository
	s3Storage         storage.Storage
}

type Service interface {
	GetMe(ctx context.Context, userUUID uuid.UUID) (*GetMeResult, error)
	GetProfile(ctx context.Context, userUUID uuid.UUID) (*GetProfileResult, error)
	UploadProfileImage(ctx context.Context, userUUID uuid.UUID, file *multipart.FileHeader) (*string, error)
	UpdateUserProfile(ctx context.Context, userUUID uuid.UUID, input UpdateUserDataInput) error
	ChangePassword(ctx context.Context, userUUID uuid.UUID, currentPassword string, newPassword string, currentSessionID *string) error
}

func NewService(userRepository userrepository.Repository, sessionRepository sessionrepository.Repository, s3Storage storage.Storage) Service {
	return &service{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		s3Storage:         s3Storage,
	}
}

func (s *service) GetMe(ctx context.Context, userUUID uuid.UUID) (*GetMeResult, error) {
	user, err := s.userRepository.FindByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	result := ToGetMeResult(user)
	return result, nil
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

func (s *service) UpdateUserProfile(ctx context.Context, userUUID uuid.UUID, input UpdateUserDataInput) error {
	formatedPhone := fmt.Sprintf("+81%d", input.Phone)

	updates := map[string]any{
		"name":    input.Name,
		"phone":   formatedPhone,
		"lat":     input.Lat,
		"lng":     input.Lng,
		"address": input.Address,
	}

	if err := s.userRepository.UpdateUserData(ctx, userUUID, updates); err != nil {
		return err
	}
	return nil
}

func (s *service) ChangePassword(ctx context.Context, userUUID uuid.UUID, currentPassword string, newPassword string, currentSessionID *string) error {
	user, err := s.userRepository.FindByUUID(ctx, userUUID)
	if err != nil {
		return err
	}

	isVerify := helper.VerifyPassword(currentPassword, user.PasswordHash)
	if !isVerify {
		return errConstant.ErrPasswordIncorrect
	}

	hashedPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return err
	}

	err = s.sessionRepository.DeleteByUserUUID(ctx, user.UUID, false, currentSessionID)
	if err != nil {
		return err
	}

	err = s.userRepository.UpdateUserData(ctx, user.UUID, map[string]any{
		"password_hash": hashedPassword,
	})

	if err != nil {
		return err
	}
	return nil
}
