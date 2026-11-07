package request

type SortingQuery struct {
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
}
