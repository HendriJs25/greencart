package response

import (
	"product-service/internal/domain/enum"
	"product-service/internal/services/product"
)

type GetAllProductsResponse struct {
	ID                 int64              `json:"id"`
	ImageUrl           *string            `json:"image_url"`
	ProductName        string             `json:"name"`
	CategoryName       string             `json:"category_name"`
	ProductDescription *string            `json:"description"`
	MinPrice           *int64             `json:"min_price"`
	MaxPrice           *int64             `json:"max_price"`
	Status             enum.ProductStatus `json:"status"`
}

type GetProductDetailResponse struct {
	ID                 int64                  `json:"id"`
	CategorySlug       string                 `json:"category_slug"`
	CategoryName       string                 `json:"category_name"`
	ProductName        string                 `json:"name"`
	ImageUrl           *string                `json:"image_url"`
	ProductDescription *string                `json:"description"`
	Unit               string                 `json:"unit"`
	Variant            int                    `json:"variant"`
	Status             enum.ProductStatus     `json:"status"`
	VariantDetail      []DetailProductVariant `json:"variant_detail"`
}

type DetailProductVariant struct {
	ID           int64              `json:"id"`
	Status       enum.ProductStatus `json:"status"`
	Stock        int64              `json:"stock"`
	Weight       int64              `json:"weight"`
	SalePrice    int64              `json:"sale_price"`
	RegularPrice int64              `json:"regular_price"`
}

type GetProductHomeResponse struct {
	ID           int64   `json:"id"`
	ImageUrl     *string `json:"image_url"`
	ProductName  string  `json:"name"`
	CategoryName string  `json:"category_name"`
	MinPrice     *int64  `json:"min_price"`
	MaxPrice     *int64  `json:"max_price"`
}

func ToGetAllProductsResponse(result []product.GetAllProductsResult) []GetAllProductsResponse {
	response := make([]GetAllProductsResponse, 0, len(result))
	for _, p := range result {
		response = append(response, GetAllProductsResponse{
			ID:                 p.ID,
			ImageUrl:           p.ImageUrl,
			ProductName:        p.ProductName,
			CategoryName:       p.CategoryName,
			ProductDescription: p.ProductDescription,
			MinPrice:           p.MinPrice,
			MaxPrice:           p.MaxPrice,
			Status:             p.Status,
		})
	}
	return response
}

func ToGetProductDetailResponse(result *product.DetailProductResult) GetProductDetailResponse {
	var detailVariant []DetailProductVariant
	for _, variant := range result.VariantDetail {
		detailVariant = append(detailVariant, DetailProductVariant{
			ID:           variant.ID,
			Status:       variant.Status,
			Stock:        variant.Stock,
			Weight:       variant.Weight,
			SalePrice:    variant.SalePrice,
			RegularPrice: variant.RegularPrice,
		})
	}
	return GetProductDetailResponse{
		ID:                 result.ID,
		CategorySlug:       result.CategorySlug,
		CategoryName:       result.CategoryName,
		ProductName:        result.ProductName,
		ImageUrl:           result.ImageUrl,
		ProductDescription: result.ProductDescription,
		Unit:               result.Unit,
		Variant:            result.Variant,
		Status:             result.Status,
		VariantDetail:      detailVariant,
	}
}

func ToGetProductHomeResponse(products []product.GetAllProductsResult) []GetProductHomeResponse {
	result := make([]GetProductHomeResponse, 0, len(products))
	for _, p := range products {
		result = append(result, GetProductHomeResponse{
			ID:           p.ID,
			ImageUrl:     p.ImageUrl,
			ProductName:  p.ProductName,
			CategoryName: p.CategoryName,
			MinPrice:     p.MinPrice,
			MaxPrice:     p.MaxPrice,
		})
	}
	return result
}
