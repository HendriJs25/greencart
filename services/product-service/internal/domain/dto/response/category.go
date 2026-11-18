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

type GetDetailCategoryResponse struct {
	ID          int64               `json:"id"`
	ParentSlug  *string             `json:"parent_slug"`
	Name        string              `json:"name"`
	Icon        string              `json:"icon"`
	Status      enum.CategoryStatus `json:"status"`
	Slug        string              `json:"slug"`
	Description *string             `json:"description"`
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

func ToGetDetailCategoryResponse(result *category.GetDetailCategoryResult) GetDetailCategoryResponse {
	return GetDetailCategoryResponse{
		ID:          result.ID,
		ParentSlug:  result.ParentSlug,
		Name:        result.Name,
		Icon:        result.Icon,
		Status:      result.Status,
		Slug:        result.Slug,
		Description: result.Description,
	}
}
