package product

import "product-service/internal/domain/enum"

type CreateProductInput struct {
	ProductName        string
	CategorySlug       string
	Unit               string
	ProductDescription *string
	Status             enum.ProductStatus
	VariantDetail      []DetailProductInput
}

type DetailProductInput struct {
	Status       enum.ProductStatus
	Stock        int64
	Weight       int64
	SalePrice    int64
	RegularPrice int64
}

type GetAllProductsResult struct {
	ID                 int64
	ImageUrl           *string
	ProductName        string
	CategoryName       string
	ProductDescription *string
	MinPrice           *int64
	MaxPrice           *int64
	Status             enum.ProductStatus
}

type DetailVariantResult struct {
	ID           int64
	Status       enum.ProductStatus
	Stock        int64
	Weight       int64
	SalePrice    int64
	RegularPrice int64
}

type DetailProductResult struct {
	ID                 int64
	CategorySlug       string
	CategoryName       string
	ProductName        string
	ImageUrl           *string
	ProductDescription *string
	Unit               string
	Variant            int
	Status             enum.ProductStatus
	VariantDetail      []DetailVariantResult
}

type UpdateDetailProductInput struct {
	ID           *int64
	Status       enum.ProductStatus
	Stock        int64
	Weight       int64
	SalePrice    int64
	RegularPrice int64
}

type UpdateProductInput struct {
	CategorySlug       string
	ProductName        string
	ProductDescription *string
	Unit               string
	Status             enum.ProductStatus
	VariantDetail      []UpdateDetailProductInput
}
