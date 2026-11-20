package product

import (
	"context"
	"errors"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/model"
	categoryrepository "product-service/internal/repository/category"
	productrepository "product-service/internal/repository/product"
)

type service struct {
	ProductRepository  productrepository.Repository
	CategoryRepository categoryrepository.Repository
}

type Service interface {
	Create(ctx context.Context, input CreateProductInput) error
}

func NewService(productRepository productrepository.Repository, categoryRepository categoryrepository.Repository) Service {
	return &service{
		ProductRepository:  productRepository,
		CategoryRepository: categoryRepository,
	}
}

func (s *service) Create(ctx context.Context, input CreateProductInput) error {
	category, err := s.CategoryRepository.FindBySlug(ctx, &input.CategorySlug)
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

	err = s.ProductRepository.Create(ctx, product)
	if err != nil {
		return err
	}
	return nil
}
