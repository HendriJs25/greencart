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
	queryinput "product-service/internal/domain/query"
	"product-service/internal/logger"
	"product-service/internal/query"
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
	GetAll(ctx context.Context, params queryinput.GetAllCategoriesQueryInput) ([]GetAllCategoriesResult, *query.Meta, error)
	GetByID(ctx context.Context, id int64) (*GetDetailCategoryResult, error)
	GetAllSlugForOptions(ctx context.Context) ([]string, error)
	Create(ctx context.Context, input CreateCategoryInput) error
	Update(ctx context.Context, id int64, update UpdateCategoryInput) error
	UploadIconImage(ctx context.Context, id int64, file *multipart.FileHeader) (*string, error)
	Delete(ctx context.Context, id int64) error
}

func NewService(categoryRepository categoryrepository.Repository, s3Storage storage.Storage) Service {
	return &service{
		categoryRepository: categoryRepository,
		s3Storage:          s3Storage,
	}
}

func (s *service) GetAll(ctx context.Context, params queryinput.GetAllCategoriesQueryInput) ([]GetAllCategoriesResult, *query.Meta, error) {
	categories, total, err := s.categoryRepository.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	result := make([]GetAllCategoriesResult, 0, len(categories))
	for _, category := range categories {
		result = append(result, GetAllCategoriesResult{
			ID:           category.ID,
			Name:         category.Name,
			Icon:         s.s3Storage.GetFileURL(category.Icon),
			Slug:         category.Slug,
			TotalProduct: category.TotalProduct,
			Status:       enum.CategoryStatusFromBool(category.Status),
		})
	}

	metadata := query.Metadata(params.Page, params.Size, total)
	return result, metadata, nil
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

	result := &GetDetailCategoryResult{
		ID:          category.ID,
		ParentSlug:  parentSlug,
		Name:        category.Name,
		Icon:        s.s3Storage.GetFileURL(category.Icon),
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

func (s *service) Update(ctx context.Context, id int64, update UpdateCategoryInput) error {
	var parentCategoryID *int64
	if update.ParentSlug != nil {
		parentCategory, err := s.categoryRepository.FindBySlug(ctx, update.ParentSlug)
		if err != nil {
			return errConstant.ErrParentCategoryMissing
		}
		parentCategoryID = &parentCategory.ID
	}

	updates := map[string]any{
		"parent_id":   parentCategoryID,
		"name":        update.Name,
		"status":      update.Status.Bool(),
		"slug":        update.Slug,
		"description": update.Description,
	}

	err := s.categoryRepository.Update(ctx, id, updates)
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

	return s.s3Storage.GetFileURL(&uploadedKey), nil
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
