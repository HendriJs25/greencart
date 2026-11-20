package product

import (
	"context"
	"errors"
	"fmt"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/model"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	Create(ctx context.Context, product *model.Product) error
	FindAll(ctx context.Context) ([]model.Product, error)
	FindParent(ctx context.Context, id int64) (*model.Product, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, product *model.Product) error {
	err := r.db.WithContext(ctx).Create(&product).Error
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	err := r.db.WithContext(ctx).Preload("Category").Where("parent_id IS NULL").Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}
	return products, nil
}

func (r *repository) FindParent(ctx context.Context, id int64) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Preload("Category").Preload("VariantDetail").Where("id = ? AND parent_id IS NULL", id).Take(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to find parent product: %w", err)
	}
	return &product, nil
}
