package query

import "product-service/internal/query"

type GetAllCategoriesQueryInput struct {
	query.Pagination
	query.Sorting
}
