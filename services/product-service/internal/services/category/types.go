package category

import (
	"product-service/internal/domain/enum"
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

type UpdateCategoryInput struct {
	ParentSlug  *string
	Name        string
	Status      enum.CategoryStatus
	Slug        string
	Description *string
}
