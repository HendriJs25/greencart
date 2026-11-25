package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           int64
	ParentID     *int64
	CategorySlug *string
	Name         *string
	Image        *string
	Description  *string
	RegularPrice *int64
	SalePrice    *int64
	Unit         *string
	Weight       *int64
	Stock        *int64
	Variant      *int
	Status       *string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    gorm.DeletedAt

	VariantDetail []Product `gorm:"foreignkey:ParentID;references:ID"`
	Category      Category  `gorm:"foreignkey:CategorySlug;references:Slug"`

	MinPrice *int64 `gorm:"column:min_price;->;-:migration"`
	MaxPrice *int64 `gorm:"column:max_price;->;-:migration"`
}

func (p *Product) TableName() string {
	return "products"
}
