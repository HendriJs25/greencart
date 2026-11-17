package category

import (
	"context"
	"fmt"
	"product-service/internal/domain/model"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindAll(ctx context.Context) ([]model.Category, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category

	subQuery := r.db.Model(&model.Product{}).Select("COUNT(*)").Where("products.category_slug = categories.slug")

	err := r.db.WithContext(ctx).Model(&model.Category{}).Select("categories.*, (?) AS total_product", subQuery).Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all categories: %w", err)
	}
	return categories, nil
}
