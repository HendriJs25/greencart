package customer

import (
	"context"
	"fmt"
	"user-service/internal/constants"
	"user-service/internal/domain/model"
	"user-service/internal/query"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindAll(ctx context.Context, params query.Pagination) ([]model.User, int64, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context, params query.Pagination) ([]model.User, int64, error) {
	var customers []model.User
	var count int64

	q := r.customerQuery(ctx)

	err := q.Session(&gorm.Session{}).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count all customers: %w", err)
	}

	err = q.Session(&gorm.Session{}).Scopes(query.Paginate(params.Page, params.Size)).Find(&customers).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find all customers: %w", err)
	}
	return customers, count, nil
}

func (r *repository) customerQuery(ctx context.Context) *gorm.DB {
	q := r.db.WithContext(ctx).Table("users").
		Joins(`JOIN user_roles ur
					ON ur.user_id = users.id`).
		Joins(`JOIN roles r
					ON r.id = ur.role_id
					AND r.deleted_at IS NULL`).
		Where("users.deleted_at IS NULL").
		Where("r.name = ?", constants.Customer)
	return q
}
