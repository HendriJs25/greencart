package product

import (
	"context"
	"errors"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/enum"
	"product-service/internal/domain/model"
	categoryrepository "product-service/internal/repository/category"
	productrepository "product-service/internal/repository/product"
	"product-service/internal/storage"
)

type service struct {
	productRepository  productrepository.Repository
	categoryRepository categoryrepository.Repository
	s3Storage          storage.Storage
}

type Service interface {
	Create(ctx context.Context, input CreateProductInput) error
	GetAll(ctx context.Context) ([]GetAllProductsResult, error)
}

func NewService(productRepository productrepository.Repository, categoryRepository categoryrepository.Repository, s3Storage storage.Storage) Service {
	return &service{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
		s3Storage:          s3Storage,
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
			ImageUrl:           s.s3Storage.GetFileURL(product.Image),
			ProductName:        *product.Name,
			CategoryName:       product.Category.Name,
			ProductDescription: product.Description,
			Status:             enum.ConvertFromDBValue(*product.Status),
		})
	}
	return result, nil
}
