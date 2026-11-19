package query

import "gorm.io/gorm"

const (
	DefaultPage = 1
	DefaultSize = 10
	MaxSize     = 100
)

type Pagination struct {
	Page int
	Size int
}

type Meta struct {
	Page       int
	Size       int
	Total      int64
	TotalPages int64
	HasNext    bool
	HasPrev    bool
}

func normalize(page, size int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}

	if size < 1 {
		size = DefaultSize
	}

	if size > MaxSize {
		size = MaxSize
	}

	return page, size
}

func Paginate(page, size int) func(*gorm.DB) *gorm.DB {
	normalizedPage, normalizedSize := normalize(page, size)

	return func(db *gorm.DB) *gorm.DB {
		offset := (normalizedPage - 1) * normalizedSize
		return db.Offset(offset).Limit(normalizedSize)
	}
}

func Metadata(page int, size int, total int64) *Meta {
	normalizedPage, normalizedSize := normalize(page, size)

	totalPages := total / int64(normalizedSize)
	if total%int64(normalizedSize) != 0 {
		totalPages++
	}

	return &Meta{
		Page:       normalizedPage,
		Size:       normalizedSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    int64(normalizedPage) < totalPages,
		HasPrev:    normalizedPage > 1 && totalPages > 1,
	}
}
