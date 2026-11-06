package customer

import (
	"context"
	"fmt"
	"user-service/internal/constants"
	"user-service/internal/domain/model"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindAll(ctx context.Context) ([]model.User, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]model.User, error) {
	var customers []model.User

	err := r.customerQuery(ctx).Find(&customers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all customers: %w", err)
	}
	return customers, nil
}

func (r *repository) customerQuery(ctx context.Context) *gorm.DB {
	query := r.db.WithContext(ctx).Table("users").
		Joins(`JOIN user_roles ur
					ON ur.user_id = users.id`).
		Joins(`JOIN roles r
					ON r.id = ur.role_id
					AND r.deleted_at IS NULL`).
		Where("users.deleted_at IS NULL").
		Where("r.name = ?", constants.Customer)
	return query
}
