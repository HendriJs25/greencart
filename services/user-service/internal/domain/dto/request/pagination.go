package request

type PaginationQuery struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
