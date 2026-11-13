package customer

import (
	"time"
	"user-service/internal/domain/model"
	"uuid"
)

type GetDetailCustomerResult struct {
	ID         int64
	UUID       uuid.UUID
	Name       string
	Email      string
	Phone      *string
	Photo      *string
	Address    *string
	Lat        *float64
	Lng        *float64
	IsVerified bool
}

type CreateCustomerInput struct {
	Name     string
	Email    string
	Password string
	Phone    *int64
	Address  *string
	Lat      *float64
	Lng      *float64
}

type UpdateCustomerInput struct {
	Name       string
	Phone      *int64
	Address    *string
	Lat        *float64
	Lng        *float64
	IsVerified *bool
}

type GetAllCustomerResult struct {
	ID        int64
	Name      string
	Email     string
	Phone     *string
	Photo     *string
	CreatedAt time.Time
}

func ToGetAllCustomerResult(customers []model.User) []GetAllCustomerResult {
	result := make([]GetAllCustomerResult, 0, len(customers))
	for _, customer := range customers {
		result = append(result, GetAllCustomerResult{
			ID:        customer.ID,
			Name:      customer.Name,
			Email:     customer.Email,
			Phone:     customer.Phone,
			Photo:     customer.Photo,
			CreatedAt: customer.CreatedAt,
		})
	}
	return result
}

func ToGetDetailCustomerResult(customer *model.User) *GetDetailCustomerResult {
	return &GetDetailCustomerResult{
		ID:         customer.ID,
		UUID:       customer.UUID,
		Name:       customer.Name,
		Email:      customer.Email,
		Phone:      customer.Phone,
		Photo:      customer.Photo,
		Address:    customer.Address,
		Lat:        customer.Lat,
		Lng:        customer.Lng,
		IsVerified: customer.IsVerified,
	}
}
