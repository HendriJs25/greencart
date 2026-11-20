package response

import (
	"product-service/internal/domain/enum"
	"product-service/internal/services/product"
)

type GetAllProductsResponse struct {
	ImageUrl           *string            `json:"image_url"`
	ProductName        string             `json:"name"`
	CategoryName       string             `json:"category_name"`
	ProductDescription *string            `json:"description"`
	Status             enum.ProductStatus `json:"status"`
}

func ToGetAllProductsResponse(result []product.GetAllProductsResult) []GetAllProductsResponse {
	response := make([]GetAllProductsResponse, 0, len(result))
	for _, p := range result {
		response = append(response, GetAllProductsResponse{
			ImageUrl:           p.ImageUrl,
			ProductName:        p.ProductName,
			CategoryName:       p.CategoryName,
			ProductDescription: p.ProductDescription,
			Status:             p.Status,
		})
	}
	return response
}
