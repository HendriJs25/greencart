package product

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"product-service/internal/common/helper"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/enum"
	"product-service/internal/domain/model"
	"product-service/internal/logger"
	"product-service/internal/repository"
	categoryrepository "product-service/internal/repository/category"
	productrepository "product-service/internal/repository/product"
	"product-service/internal/storage"
	"strings"
	"time"
	"uuid"
)

type service struct {
	productRepository  productrepository.Repository
	categoryRepository categoryrepository.Repository
	s3Storage          storage.Storage
	transactionManager repository.TransactionManager
}

type Service interface {
	Create(ctx context.Context, input CreateProductInput) error
	GetAll(ctx context.Context) ([]GetAllProductsResult, error)
	GetParent(ctx context.Context, id int64) (*DetailProductResult, error)
	UploadImage(ctx context.Context, id int64, file *multipart.FileHeader) (*string, error)
	UpdateProduct(ctx context.Context, id int64, input UpdateProductInput) error
}

func NewService(productRepository productrepository.Repository, categoryRepository categoryrepository.Repository, s3Storage storage.Storage, transactionManager repository.TransactionManager) Service {
	return &service{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
		s3Storage:          s3Storage,
		transactionManager: transactionManager,
	}
}

func (s *service) Create(ctx context.Context, input CreateProductInput) error {
	category, err := s.categoryRepository.FindBySlug(ctx, &input.CategorySlug)
	if err != nil {
		if errors.Is(err, errConstant.ErrNotFound) {
			return errConstant.ErrCategoryNotFound
		}
		return err
	}

	variant := len(input.VariantDetail)
	product := &model.Product{
		CategorySlug: &category.Slug,
		Name:         &input.ProductName,
		Description:  input.ProductDescription,
		Unit:         &input.Unit,
		Variant:      &variant,
		Status:       input.Status.ConvertToDBValue(),
	}

	for _, detail := range input.VariantDetail {
		product.VariantDetail = append(product.VariantDetail, model.Product{
			Status:       detail.Status.ConvertToDBValue(),
			Stock:        &detail.Stock,
			Weight:       &detail.Weight,
			SalePrice:    &detail.SalePrice,
			RegularPrice: &detail.RegularPrice,
		})
	}

	err = s.productRepository.Create(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAll(ctx context.Context) ([]GetAllProductsResult, error) {
	products, err := s.productRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]GetAllProductsResult, 0, len(products))
	for _, product := range products {
		result = append(result, GetAllProductsResult{
			ID:                 product.ID,
			ImageUrl:           s.s3Storage.GetFileURL(product.Image),
			ProductName:        *product.Name,
			CategoryName:       product.Category.Name,
			ProductDescription: product.Description,
			Status:             enum.ConvertFromDBValue(*product.Status),
		})
	}
	return result, nil
}

func (s *service) GetParent(ctx context.Context, id int64) (*DetailProductResult, error) {
	product, err := s.productRepository.FindParent(ctx, id)
	if err != nil {
		return nil, err
	}

	var variantDetail []DetailVariantResult
	for _, detail := range product.VariantDetail {
		variantDetail = append(variantDetail, DetailVariantResult{
			ID:           detail.ID,
			Status:       enum.ConvertFromDBValue(*detail.Status),
			Stock:        *detail.Stock,
			Weight:       *detail.Weight,
			SalePrice:    *detail.SalePrice,
			RegularPrice: *detail.RegularPrice,
		})
	}

	result := &DetailProductResult{
		ID:                 product.ID,
		CategorySlug:       product.Category.Slug,
		CategoryName:       product.Category.Name,
		ProductName:        *product.Name,
		ImageUrl:           s.s3Storage.GetFileURL(product.Image),
		ProductDescription: product.Description,
		Unit:               *product.Unit,
		Variant:            *product.Variant,
		Status:             enum.ConvertFromDBValue(*product.Status),
		VariantDetail:      variantDetail,
	}

	return result, nil
}

func (s *service) UploadImage(ctx context.Context, id int64, file *multipart.FileHeader) (*string, error) {
	contentType, err := helper.ValidateImage(file)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepository.FindParent(ctx, id)
	if err != nil {
		return nil, err
	}

	oldKey := product.Image

	path := fmt.Sprintf("products/%d/images/%s%s", product.ID, uuid.NewV4().String(), strings.ToLower(filepath.Ext(file.Filename)))

	uploadedKey, err := s.s3Storage.UploadFile(ctx, file, path, contentType)
	if err != nil {
		return nil, err
	}

	if err := s.productRepository.UpdateProductImage(ctx, product.ID, uploadedKey); err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		_ = s.s3Storage.DeleteFile(cleanupCtx, uploadedKey)
		cancel()
		return nil, err
	}

	if oldKey != nil && *oldKey != "" {
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		if err := s.s3Storage.DeleteFile(cleanupCtx, *oldKey); err != nil {
			logger.WithContext(ctx).WithError(err).Error("failed to delete old image")
		}
		cancel()
	}

	return s.s3Storage.GetFileURL(&uploadedKey), nil
}

func (s *service) UpdateProduct(ctx context.Context, id int64, input UpdateProductInput) error {
	if len(input.VariantDetail) == 0 {
		return errConstant.ErrAtLeastOneVariantRequired
	}

	err := s.transactionManager.WithinTransaction(ctx, func(repositories *repository.Registry) error {
		category, err := repositories.CategoryRepository.FindBySlug(ctx, &input.CategorySlug)
		if err != nil {
			if errors.Is(err, errConstant.ErrNotFound) {
				return errConstant.ErrCategoryNotFound
			}
			return err
		}

		parentUpdates := map[string]any{
			"category_slug": category.Slug,
			"name":          input.ProductName,
			"description":   input.ProductDescription,
			"unit":          input.Unit,
			"status":        input.Status.ConvertToDBValue(),
		}

		if err := repositories.ProductRepository.UpdateProductParent(ctx, id, parentUpdates); err != nil {
			return err
		}

		keepIDs := make([]int64, 0, len(input.VariantDetail))
		for _, variant := range input.VariantDetail {
			if variant.ID == nil {
				createVariantInput := model.Product{
					ParentID:     &id,
					Status:       variant.Status.ConvertToDBValue(),
					Stock:        &variant.Stock,
					Weight:       &variant.Weight,
					SalePrice:    &variant.SalePrice,
					RegularPrice: &variant.RegularPrice,
				}
				if err := repositories.ProductRepository.Create(ctx, &createVariantInput); err != nil {
					return err
				}
				keepIDs = append(keepIDs, createVariantInput.ID)
				continue
			}

			variantUpdate := map[string]any{
				"status":        variant.Status.ConvertToDBValue(),
				"stock":         variant.Stock,
				"weight":        variant.Weight,
				"sale_price":    variant.SalePrice,
				"regular_price": variant.RegularPrice,
			}

			if err := repositories.ProductRepository.UpdateProductChild(ctx, *variant.ID, id, variantUpdate); err != nil {
				return err
			}

			keepIDs = append(keepIDs, *variant.ID)
		}

		if err := repositories.ProductRepository.DeleteProductChild(ctx, keepIDs, id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
