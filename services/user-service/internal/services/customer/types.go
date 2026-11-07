package customer

import (
	"user-service/internal/domain/model"
)

type GetAllCustomerResult struct {
	ID    int64
	Name  string
	Email string
	Phone *string
	Photo *string
}

func ToGetAllCustomerResult(customers []model.User) []GetAllCustomerResult {
	result := make([]GetAllCustomerResult, 0, len(customers))
	for _, customer := range customers {
		result = append(result, GetAllCustomerResult{
			ID:    customer.ID,
			Name:  customer.Name,
			Email: customer.Email,
			Phone: customer.Phone,
			Photo: customer.Photo,
		})
	}
	return result
}
