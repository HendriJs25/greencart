package request

import (
	queryinput "user-service/internal/domain/query"
	"user-service/internal/query"
)

type GetAllCustomersQuery struct {
	PaginationQuery
	SortingQuery
}

func ToGetAllCustomerQueryInput(params GetAllCustomersQuery) queryinput.GetAllCustomerQuery {
	return queryinput.GetAllCustomerQuery{
		Pagination: query.Pagination{
			Page: params.Page,
			Size: params.Size,
		},
		Sorting: query.Sorting{
			SortBy:    params.SortBy,
			SortOrder: params.SortOrder,
		},
	}
}
