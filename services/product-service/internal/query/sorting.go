package query

import (
	"strings"

	"gorm.io/gorm"
)

// Sort returns a GORM scope for sorting query results.
// allowedColumns must include (default sort column)
// and "id" (secondary sort column for stable pagination).
// All column values must be trusted database column names.

const (
	DefaultSortBy    = "created_at"
	DefaultSortOrder = "desc"
)

type Sorting struct {
	SortBy    string
	SortOrder string
}

func Sort(column string, direction string, allowedColumns map[string]string) func(*gorm.DB) *gorm.DB {
	column = strings.ToLower(strings.TrimSpace(column))
	direction = strings.ToLower(strings.TrimSpace(direction))
	dbColumn, ok := allowedColumns[column]
	if !ok {
		dbColumn = allowedColumns[DefaultSortBy]
	}

	if direction != "asc" && direction != "desc" {
		direction = DefaultSortOrder
	}

	idColumn := allowedColumns["id"]

	return func(db *gorm.DB) *gorm.DB {
		s := db.Order(dbColumn + " " + direction)

		if idColumn != "" && dbColumn != idColumn {
			s = db.Order(idColumn + " " + direction)
		}
		return s
	}
}
