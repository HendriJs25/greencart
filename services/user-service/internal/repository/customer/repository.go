package customer

import (
	"context"
	"fmt"
	"user-service/internal/constants"
	"user-service/internal/domain/model"
	queryinput "user-service/internal/domain/query"
	"user-service/internal/query"

	"gorm.io/gorm"
)

var customerSortColumns = map[string]string{
	"id":         "users.id",
	"name":       "users.name",
	"created_at": "users.created_at",
}

var customerSearchColumns = []string{
	"users.name", "users.email", "users.address",
}

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindAll(ctx context.Context, params queryinput.GetAllCustomerQuery) ([]model.User, int64, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context, params queryinput.GetAllCustomerQuery) ([]model.User, int64, error) {
	var customers []model.User
	var count int64

	q := r.customerQuery(ctx).Scopes(query.Search(params.Search, customerSearchColumns...))

	err := q.Session(&gorm.Session{}).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count all customers: %w", err)
	}

	err = q.Session(&gorm.Session{}).Select("users.*").Scopes(query.Sort(params.SortBy, params.SortOrder, customerSortColumns), query.Paginate(params.Page, params.Size)).Find(&customers).Error

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
