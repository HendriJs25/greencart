package request

import (
	"product-service/internal/domain/enum"
	queryinput "product-service/internal/domain/query"
	"product-service/internal/services/product"
)

type GetAllProductsQuery struct {
	PaginationQuery
	SortingQuery

	Search string `form:"search"`
}

type CreateProductRequest struct {
	ProductName        string                 `json:"name" validate:"required,notblank"`
	CategorySlug       string                 `json:"category_slug" validate:"required"`
	Unit               string                 `json:"unit" validate:"required,notblank"`
	ProductDescription *string                `json:"description"`
	Status             enum.ProductStatus     `json:"status" validate:"required,oneof=有効 無効 下書き"`
	VariantDetail      []DetailProductRequest `json:"variant_detail" validate:"required"`
}

type DetailProductRequest struct {
	Status       enum.ProductStatus `json:"status" validate:"required,oneof=有効 無効 下書き"`
	Stock        int64              `json:"stock" validate:"required,number"`
	Weight       int64              `json:"weight" validate:"required,number"`
	SalePrice    int64              `json:"sale_price" validate:"required,number"`
	RegularPrice int64              `json:"regular_price" validate:"required,number"`
}

type UpdateDetailProductRequest struct {
	ID           *int64             `json:"id"`
	Status       enum.ProductStatus `json:"status" validate:"required,oneof=有効 無効 下書き"`
	Stock        int64              `json:"stock" validate:"required,number"`
	Weight       int64              `json:"weight" validate:"required,number"`
	SalePrice    int64              `json:"sale_price" validate:"required,number"`
	RegularPrice int64              `json:"regular_price" validate:"required,number"`
}

type UpdateProductRequest struct {
	CategorySlug       string                       `json:"category_slug" validate:"required"`
	ProductName        string                       `json:"name" validate:"required,notblank"`
	ProductDescription *string                      `json:"description"`
	Unit               string                       `json:"unit" validate:"required,notblank"`
	Status             enum.ProductStatus           `json:"status" validate:"required,oneof=有効 無効 下書き"`
	VariantDetail      []UpdateDetailProductRequest `json:"variant_detail" validate:"required"`
}

func ToGetAllProductsQueryInput(params GetAllProductsQuery) queryinput.GetAllProductsQueryInput {
	return queryinput.GetAllProductsQueryInput{
		Page:      params.Page,
		Size:      params.Size,
		SortOrder: params.SortOrder,
		SortBy:    params.SortBy,
		Search:    params.Search,
	}
}

func ToCreateProductInput(req CreateProductRequest) product.CreateProductInput {
	var variantDetailInput []product.DetailProductInput
	for _, v := range req.VariantDetail {
		variantDetailInput = append(variantDetailInput, product.DetailProductInput{
			Status:       v.Status,
			Stock:        v.Stock,
			Weight:       v.Weight,
			SalePrice:    v.SalePrice,
			RegularPrice: v.RegularPrice,
		})
	}

	return product.CreateProductInput{
		ProductName:        req.ProductName,
		CategorySlug:       req.CategorySlug,
		Unit:               req.Unit,
		ProductDescription: req.ProductDescription,
		Status:             req.Status,
		VariantDetail:      variantDetailInput,
	}
}

func ToUpdateProductInput(req UpdateProductRequest) product.UpdateProductInput {
	var variantUpdateInput []product.UpdateDetailProductInput
	for _, v := range req.VariantDetail {
		variantUpdateInput = append(variantUpdateInput, product.UpdateDetailProductInput{
			ID:           v.ID,
			Status:       v.Status,
			Stock:        v.Stock,
			Weight:       v.Weight,
			SalePrice:    v.SalePrice,
			RegularPrice: v.RegularPrice,
		})
	}

	return product.UpdateProductInput{
		CategorySlug:       req.CategorySlug,
		ProductName:        req.ProductName,
		ProductDescription: req.ProductDescription,
		Unit:               req.Unit,
		Status:             req.Status,
		VariantDetail:      variantUpdateInput,
	}
}
