package request

import (
	queryinput "user-service/internal/domain/query"
	"user-service/internal/query"
	"user-service/internal/services/customer"
)

type GetAllCustomersQuery struct {
	PaginationQuery
	SortingQuery

	Search string `form:"search"`
}

type CreateCustomerRequest struct {
	Name                 string   `json:"name" validate:"required,notblank"`
	Email                string   `json:"email" validate:"required,email"`
	Password             string   `json:"password" validate:"required,validpassword"`
	PasswordConfirmation string   `json:"password_confirmation" validate:"required,eqfield=Password"`
	Phone                *string  `json:"phone" validate:"omitempty,number"`
	Address              *string  `json:"address" validate:"omitempty,notblank"`
	Lat                  *float64 `json:"lat" validate:"omitempty,gte=-90,lte=90"`
	Lng                  *float64 `json:"lng" validate:"omitempty,gte=-180,lte=180"`
}

func ToCreateCustomerInput(req CreateCustomerRequest) customer.CreateCustomerInput {
	return customer.CreateCustomerInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Address:  req.Address,
		Lat:      req.Lat,
		Lng:      req.Lng,
	}
}

func ToGetAllCustomerQueryInput(params GetAllCustomersQuery) queryinput.GetAllCustomerQuery {
	return queryinput.GetAllCustomerQuery{
		Pagination: query.Pagination{
			Page: params.Page,
			Size: params.Size,
		},
		Sorting: query.Sorting{
			SortBy:    params.SortBy,
			SortOrder: params.SortOrder,
		},
		Search: params.Search,
	}
}
