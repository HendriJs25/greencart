package request

import "product-service/internal/domain/enum"

type CreateCategoryRequest struct {
	ParentID    *int64              `json:"parent_id"`
	Name        string              `json:"name" validate:"required,notblank"`
	Icon        string              `json:"icon" validate:"required,notblank"`
	Status      enum.CategoryStatus `json:"status" validate:"required, oneof=有効 無効"`
	Slug        string              `json:"slug" validate:"required,lowercase_hyphen"`
	Description *string             `json:"description"`
}
