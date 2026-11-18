package response

import (
	"product-service/internal/domain/enum"
	"product-service/internal/services/category"
)

type GetAllCategoriesResponse struct {
	ID           int64               `json:"id"`
	Name         string              `json:"name"`
	Icon         string              `json:"icon"`
	Slug         string              `json:"slug"`
	TotalProduct int64               `json:"total_product"`
	Status       enum.CategoryStatus `json:"status"`
}

func ToGetAllCategoriesResponse(categories []category.GetAllCategoriesResult) []GetAllCategoriesResponse {
	result := make([]GetAllCategoriesResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, GetAllCategoriesResponse{
			ID:           c.ID,
			Name:         c.Name,
			Icon:         c.Icon,
			Slug:         c.Slug,
			TotalProduct: c.TotalProduct,
			Status:       c.Status,
		})
	}
	return result
}
