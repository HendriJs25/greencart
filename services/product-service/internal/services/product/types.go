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
	ImageUrl           *string
	ProductName        string
	CategoryName       string
	ProductDescription *string
	Status             enum.ProductStatus
}
