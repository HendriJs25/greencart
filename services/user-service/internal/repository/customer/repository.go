package customer

import (
	"context"
	"errors"
	"fmt"
	"user-service/internal/constants"
	errConstant "user-service/internal/constants/error"
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
	FindByID(ctx context.Context, id int64) (*model.User, error)
	UpdateData(ctx context.Context, id int64, updates map[string]any) error
	DeleteByID(ctx context.Context, id int64) error
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

func (r *repository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var customer model.User

	err := r.customerQuery(ctx).Where("users.id = ?", id).Take(&customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to find customer by id: %w", err)
	}
	return &customer, nil
}

func (r *repository) UpdateData(ctx context.Context, id int64, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
		Where(`EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = users.id AND r.name = ?)`, constants.Customer).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update data customer: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errConstant.ErrUserNotFound
	}
	return nil
}

func (r *repository) DeleteByID(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Where("users.id = ?", id).
		Where(`EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = users.id AND r.name = ?)`, constants.Customer).
		Delete(&model.User{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete customer by id: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errConstant.ErrUserNotFound
	}
	return nil
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
