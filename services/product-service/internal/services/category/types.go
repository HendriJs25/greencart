package category

import (
	"product-service/internal/domain/enum"
	"product-service/internal/domain/model"
)

type GetAllCategoriesResult struct {
	ID           int64
	Name         string
	Icon         *string
	Slug         string
	TotalProduct int64
	Status       enum.CategoryStatus
}

type GetDetailCategoryResult struct {
	ID          int64
	ParentSlug  *string
	Name        string
	Icon        *string
	Status      enum.CategoryStatus
	Slug        string
	Description *string
}

type CreateCategoryInput struct {
	ParentSlug  *string
	Name        string
	Status      enum.CategoryStatus
	Slug        string
	Description *string
}

func ToGetAllCategoriesResult(categories []model.Category, iconUrl *string) []GetAllCategoriesResult {
	result := make([]GetAllCategoriesResult, 0, len(categories))
	for _, category := range categories {
		result = append(result, GetAllCategoriesResult{
			ID:           category.ID,
			Name:         category.Name,
			Icon:         iconUrl,
			Slug:         category.Slug,
			TotalProduct: category.TotalProduct,
			Status:       enum.CategoryStatusFromBool(category.Status),
		})
	}
	return result
}
