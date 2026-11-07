package query

import "user-service/internal/query"

type GetAllCustomerQuery struct {
	query.Pagination
	query.Sorting
}
