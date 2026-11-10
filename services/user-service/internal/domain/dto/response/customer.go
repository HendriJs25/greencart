package response

import (
	"time"
	"user-service/internal/services/customer"
	"uuid"
)

type GetAllCustomerResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone"`
	Photo     *string   `json:"photo"`
	CreatedAt time.Time `json:"created_at"`
}

type GetDetailCustomerResponse struct {
	ID         int64     `json:"id"`
	UUID       uuid.UUID `json:"uuid"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      *string   `json:"phone"`
	Photo      *string   `json:"photo"`
	Address    *string   `json:"address"`
	Lat        *float64  `json:"lat"`
	Lng        *float64  `json:"lng"`
	IsVerified bool      `json:"is_verified"`
}

func ToGetAllCustomerResponse(customers []customer.GetAllCustomerResult) []GetAllCustomerResponse {
	result := make([]GetAllCustomerResponse, 0, len(customers))
	for _, c := range customers {
		result = append(result, GetAllCustomerResponse{
			ID:        c.ID,
			Name:      c.Name,
			Email:     c.Email,
			Phone:     c.Phone,
			Photo:     c.Photo,
			CreatedAt: c.CreatedAt,
		})
	}
	return result
}

func ToGetDetailCustomerResponse(customer *customer.GetDetailCustomerResult) GetDetailCustomerResponse {
	return GetDetailCustomerResponse{
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
