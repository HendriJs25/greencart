package response

import "user-service/internal/services/customer"

type GetAllCustomerResponse struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Phone *string `json:"phone"`
	Photo *string `json:"photo"`
}

func ToGetAllCustomerResponse(customers []customer.GetAllCustomerResult) []GetAllCustomerResponse {
	result := make([]GetAllCustomerResponse, 0, len(customers))
	for _, c := range customers {
		result = append(result, GetAllCustomerResponse{
			ID:    c.ID,
			Name:  c.Name,
			Email: c.Email,
			Phone: c.Phone,
			Photo: c.Photo,
		})
	}
	return result
}
