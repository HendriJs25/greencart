package query

import "gorm.io/gorm"

func ProductPriceRange(minPrice, maxPrice *int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if minPrice != nil {
			db = db.Where("prices.min_price >= ?", *minPrice)
		}

		if maxPrice != nil {
			db = db.Where("prices.max_price <= ?", *maxPrice)
		}

		return db
	}
}
