package request

import "user-service/internal/query"

type Pagination struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func ToPaginationInput(pagination Pagination) query.Pagination {
	return query.Pagination{
		Page: pagination.Page,
		Size: pagination.Size,
	}
}
