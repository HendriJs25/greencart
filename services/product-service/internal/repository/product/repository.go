package product

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

var productSortColumns = map[string]string{
	"id":            "products.id",
	"created_at":    "products.created_at",
	"name":          "products.name",
	"category_slug": "products.category_slug",
}

var productSearchColumns = []string{
	"products.name", "categories.name",
}

type repository struct {
	db *gorm.DB
}

type Repository interface {
	Create(ctx context.Context, product *model.Product) error
	FindAll(ctx context.Context, params queryinput.GetAllProductsQueryInput) ([]model.Product, int64, error)
	FindAllActive(ctx context.Context, params queryinput.GetAllProductsQueryInput) ([]model.Product, int64, error)
	FindParent(ctx context.Context, id int64) (*model.Product, error)
	UpdateProductImage(ctx context.Context, id int64, imageKey string) error
	UpdateProductParent(ctx context.Context, id int64, updates map[string]any) error
	UpdateProductChild(ctx context.Context, id int64, parentID int64, updates map[string]any) error
	DeleteProductChild(ctx context.Context, id []int64, parentID int64) error
	Delete(ctx context.Context, id int64) error
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

func (r *repository) FindAll(ctx context.Context, params queryinput.GetAllProductsQueryInput) ([]model.Product, int64, error) {
	var products []model.Product
	var count int64
	priceSubQuery := r.db.Table("products").
		Select(`parent_id, MIN(sale_price) AS min_price, MAX(sale_price) AS max_price`).
		Where("parent_id IS NOT NULL").Where("deleted_at IS NULL").Group("parent_id")

	q := r.db.WithContext(ctx).Model(&model.Product{}).
		Joins(`LEFT JOIN categories ON categories.slug = products.category_slug`).
		Joins(`LEFT JOIN (?) AS prices ON prices.parent_id = products.id`, priceSubQuery).
		Preload("Category").Where("products.parent_id IS NULL").
		Scopes(query.Search(params.Search, productSearchColumns...), query.ProductPriceRange(params.MinPrice, params.MaxPrice))

	err := q.Session(&gorm.Session{}).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	err = q.Session(&gorm.Session{}).
		Select(`products.*,prices.min_price,prices.max_price`).
		Scopes(query.Sort(params.SortBy, params.SortOrder, productSortColumns), query.Paginate(params.Page, params.Size)).
		Find(&products).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}

	return products, count, nil
}

func (r *repository) FindAllActive(ctx context.Context, params queryinput.GetAllProductsQueryInput) ([]model.Product, int64, error) {
	var products []model.Product
	var count int64
	priceSubQuery := r.db.Table("products").
		Select(`parent_id, MIN(sale_price) AS min_price, MAX(sale_price) AS max_price`).
		Where("parent_id IS NOT NULL").Where("deleted_at IS NULL").Group("parent_id")

	q := r.db.WithContext(ctx).Model(&model.Product{}).
		Joins(`LEFT JOIN categories ON categories.slug = products.category_slug`).
		Joins(`LEFT JOIN (?) AS prices ON prices.parent_id = products.id`, priceSubQuery).
		Preload("Category").Where("products.parent_id IS NULL AND products.status = 'active'").
		Scopes(query.Search(params.Search, productSearchColumns...), query.ProductPriceRange(params.MinPrice, params.MaxPrice))

	err := q.Session(&gorm.Session{}).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	err = q.Session(&gorm.Session{}).
		Select(`products.*,prices.min_price,prices.max_price`).
		Scopes(query.Sort(params.SortBy, params.SortOrder, productSortColumns), query.Paginate(params.Page, params.Size)).
		Find(&products).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}

	return products, count, nil
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

func (r *repository) UpdateProductImage(ctx context.Context, id int64, imageKey string) error {
	result := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ? AND parent_id IS NULL", id).Update("image", imageKey)
	if result.Error != nil {
		return fmt.Errorf("failed to update product image: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrProductNotFound
	}
	return nil
}

func (r *repository) UpdateProductParent(ctx context.Context, id int64, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ? AND parent_id IS NULL", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update product parent: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrProductNotFound
	}
	return nil
}

func (r *repository) UpdateProductChild(ctx context.Context, id int64, parentID int64, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ? AND parent_id = ?", id, parentID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update product child: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrProductNotFound
	}
	return nil
}

func (r *repository) DeleteProductChild(ctx context.Context, id []int64, parentID int64) error {
	err := r.db.WithContext(ctx).Where("parent_id = ? AND id NOT IN ?", parentID, id).Delete(&model.Product{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete product child: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Where("id = ? AND parent_id IS NULL", id).Delete(&model.Product{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrProductNotFound
	}
	return nil
}
