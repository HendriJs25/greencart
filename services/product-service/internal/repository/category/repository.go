package category

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
	FindAll(ctx context.Context) ([]model.Category, error)
	FindAllSlugForOptions(ctx context.Context) ([]string, error)
	FindBySlug(ctx context.Context, slug *string) (*model.Category, error)
	FindByID(ctx context.Context, id int64) (*model.Category, error)
	IsSlugExist(ctx context.Context, slug string) (bool, error)
	Create(ctx context.Context, category *model.Category) error
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

func (r *repository) FindAllSlugForOptions(ctx context.Context) ([]string, error) {
	var slugs []string
	err := r.db.WithContext(ctx).Model(&model.Category{}).Pluck("slug", &slugs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all slug: %w", err)
	}
	return slugs, nil
}

func (r *repository) FindBySlug(ctx context.Context, slug *string) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).Where("categories.slug = ?", slug).Take(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find category by slug: %w", err)
	}
	return &category, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).Preload("Parent").Where("id = ?", id).Take(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find category by id: %w", err)
	}
	return &category, nil
}

func (r *repository) IsSlugExist(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Category{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if slug exists: %w", err)
	}
	return count > 0, nil
}

func (r *repository) Create(ctx context.Context, category *model.Category) error {
	err := r.db.WithContext(ctx).Create(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errConstant.ErrCategoryAlreadyExists
		}
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}
