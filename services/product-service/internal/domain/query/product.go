package query

import "product-service/internal/query"

type GetAllProductsQueryInput struct {
	query.Pagination
	query.Sorting

	Search string
}
