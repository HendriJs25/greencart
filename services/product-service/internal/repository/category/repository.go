package category

import (
	"context"
	"errors"
	"fmt"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/model"
	queryinput "product-service/internal/domain/query"
	"product-service/internal/query"

	"gorm.io/gorm"
)

var categorySortColumns = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"name":       "name",
}

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindAll(ctx context.Context, params queryinput.GetAllCategoriesQueryInput) ([]model.Category, int64, error)
	FindAllSlugForOptions(ctx context.Context) ([]string, error)
	FindBySlug(ctx context.Context, slug *string) (*model.Category, error)
	FindByID(ctx context.Context, id int64) (*model.Category, error)
	IsSlugExist(ctx context.Context, slug string) (bool, error)
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, id int64, updates map[string]any) error
	UpdateIconKeyByID(ctx context.Context, id int64, iconKey string) error
	Delete(ctx context.Context, id int64) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context, params queryinput.GetAllCategoriesQueryInput) ([]model.Category, int64, error) {
	var categories []model.Category
	var count int64

	subQuery := r.db.Model(&model.Product{}).Select("COUNT(*)").Where("products.category_slug = categories.slug")
	q := r.db.WithContext(ctx).Model(&model.Category{}).Select("categories.*, (?) AS total_product", subQuery)

	err := q.Session(&gorm.Session{}).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count all categories: %w", err)
	}

	err = q.Session(&gorm.Session{}).Scopes(query.Sort(params.SortBy, params.SortOrder, categorySortColumns), query.Paginate(params.Page, params.Size)).Find(&categories).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find all categories: %w", err)
	}
	return categories, count, nil
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

func (r *repository) Update(ctx context.Context, id int64, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errConstant.ErrCategoryAlreadyExists
		}
		return fmt.Errorf("failed to update category: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrNotFound
	}
	return nil
}

func (r *repository) UpdateIconKeyByID(ctx context.Context, id int64, iconKey string) error {
	result := r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", id).Update("icon", iconKey)
	if result.Error != nil {
		return fmt.Errorf("failed to update icon key: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrNotFound
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).
		Where(`NOT EXISTS (SELECT 1 FROM categories WHERE parent_id = ?)`, id).Where("id = ?", id).Delete(&model.Category{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete category: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrCategoryHasChildren
	}
	return nil
}
