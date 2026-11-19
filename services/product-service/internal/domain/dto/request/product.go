package request

type CreateProductRequest struct {
	ProductName        string                 `json:"name"`
	CategorySlug       string                 `json:"category_slug"`
	Unit               string                 `json:"unit"`
	Variant            int                    `json:"variant"`
	ProductDescription string                 `json:"description"`
	Status             string                 `json:"status"`
	VariantDetail      []ProductDetailRequest `json:"variant_detail"`
}

type ProductDetailRequest struct {
	Stock        int   `json:"stock"`
	Weight       int   `json:"weight"`
	SalePrice    int64 `json:"sale_price"`
	RegularPrice int64 `json:"regular_price"`
}
