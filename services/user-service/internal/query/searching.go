package query

import (
	"strings"

	"gorm.io/gorm"
)

// Search performs case-insensitive partial matching using PostgreSQL ILIKE.
// Columns must be trusted database column names, never raw user input.
// Multiple columns are combined using OR.
func Search(keyword string, columns ...string) func(*gorm.DB) *gorm.DB {
	keyword = strings.TrimSpace(keyword)

	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" || len(columns) == 0 {
			return db
		}

		// Escape ILIKE wildcard characters
		escaped := strings.NewReplacer(
			`\`, `\\`,
			`%`, `\%`,
			`_`, `\_`,
		).Replace(keyword)

		pattern := "%" + escaped + "%"

		conditions := make([]string, 0, len(columns))
		args := make([]any, 0, len(columns))

		for _, column := range columns {
			conditions = append(conditions, column+" ILIKE ?")
			args = append(args, pattern)
		}

		return db.Where(
			"("+strings.Join(conditions, " OR ")+")",
			args...,
		)
	}
}
