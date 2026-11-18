package category

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/enum"
	"product-service/internal/domain/model"
	"product-service/internal/logger"
	categoryrepository "product-service/internal/repository/category"
	"product-service/internal/storage"
	"strings"
	"time"
	"uuid"
)

type service struct {
	categoryRepository categoryrepository.Repository
	s3Storage          storage.Storage
}

type Service interface {
	GetAll(ctx context.Context) ([]GetAllCategoriesResult, error)
	GetByID(ctx context.Context, id int64) (*GetDetailCategoryResult, error)
	GetAllSlugForOptions(ctx context.Context) ([]string, error)
	Create(ctx context.Context, input CreateCategoryInput) error
	UploadIconImage(ctx context.Context, id int64, file *multipart.FileHeader) (*string, error)
	Delete(ctx context.Context, id int64) error
}

func NewService(categoryRepository categoryrepository.Repository, s3Storage storage.Storage) Service {
	return &service{
		categoryRepository: categoryRepository,
		s3Storage:          s3Storage,
	}
}

func (s *service) GetAll(ctx context.Context) ([]GetAllCategoriesResult, error) {
	categories, err := s.categoryRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]GetAllCategoriesResult, 0, len(categories))
	var iconURL *string
	for _, category := range categories {
		if category.Icon != nil {
			iconURL = s.s3Storage.GetFileURL(*category.Icon)
		} else {
			iconURL = nil
		}
		result = append(result, GetAllCategoriesResult{
			ID:           category.ID,
			Name:         category.Name,
			Icon:         iconURL,
			Slug:         category.Slug,
			TotalProduct: category.TotalProduct,
			Status:       enum.CategoryStatusFromBool(category.Status),
		})
	}
	return result, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*GetDetailCategoryResult, error) {
	category, err := s.categoryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var parentSlug *string
	if category.Parent != nil {
		parentSlug = &category.Parent.Slug
	}

	var iconURL *string
	if category.Icon != nil {
		iconURL = s.s3Storage.GetFileURL(*category.Icon)
	}

	result := &GetDetailCategoryResult{
		ID:          category.ID,
		ParentSlug:  parentSlug,
		Name:        category.Name,
		Icon:        iconURL,
		Status:      enum.CategoryStatusFromBool(category.Status),
		Slug:        category.Slug,
		Description: category.Description,
	}

	return result, nil
}

func (s *service) GetAllSlugForOptions(ctx context.Context) ([]string, error) {
	return s.categoryRepository.FindAllSlugForOptions(ctx)
}

func (s *service) Create(ctx context.Context, input CreateCategoryInput) error {
	exist, err := s.categoryRepository.IsSlugExist(ctx, input.Slug)
	if err != nil {
		return err
	}
	if exist {
		return errConstant.ErrCategoryAlreadyExists
	}

	var parentCategoryID *int64
	if input.ParentSlug != nil {
		parentCategory, err := s.categoryRepository.FindBySlug(ctx, input.ParentSlug)
		if err != nil {
			return err
		}
		parentCategoryID = &parentCategory.ID
	}

	inputModel := &model.Category{
		ParentID:    parentCategoryID,
		Name:        input.Name,
		Status:      input.Status.Bool(),
		Slug:        input.Slug,
		Description: input.Description,
	}

	err = s.categoryRepository.Create(ctx, inputModel)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UploadIconImage(ctx context.Context, id int64, file *multipart.FileHeader) (*string, error) {
	contentType, err := validateProfileImage(file)
	if err != nil {
		return nil, err
	}

	category, err := s.categoryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	oldKey := category.Icon

	path := fmt.Sprintf("categories/%d/icon/%s%s", category.ID, uuid.NewV4().String(), strings.ToLower(filepath.Ext(file.Filename)))

	uploadedKey, err := s.s3Storage.UploadFile(ctx, file, path, contentType)
	if err != nil {
		return nil, err
	}

	if err := s.categoryRepository.UpdateIconKeyByID(ctx, category.ID, uploadedKey); err != nil {
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

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.categoryRepository.Delete(ctx, id)
}
