package request

import (
	"product-service/internal/domain/enum"
	"product-service/internal/services/category"
)

type CreateCategoryRequest struct {
	ParentSlug  *string             `json:"parent_slug"`
	Name        string              `json:"name" validate:"required,notblank"`
	Status      enum.CategoryStatus `json:"status" validate:"required,oneof=有効 無効"`
	Slug        string              `json:"slug" validate:"required,lowercase_hyphen"`
	Description *string             `json:"description"`
}

func ToCreateCategoryInput(req CreateCategoryRequest) category.CreateCategoryInput {
	return category.CreateCategoryInput{
		ParentSlug:  req.ParentSlug,
		Name:        req.Name,
		Status:      req.Status,
		Slug:        req.Slug,
		Description: req.Description,
	}
}
