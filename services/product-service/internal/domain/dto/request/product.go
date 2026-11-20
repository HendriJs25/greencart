package request

import (
	"product-service/internal/domain/enum"
	"product-service/internal/services/product"
)

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
