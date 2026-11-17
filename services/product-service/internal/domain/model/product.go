package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           int64
	ParentID     *int64
	CategorySlug string
	Name         string
	Image        string
	Description  *string
	RegularPrice int64
	SalePrice    int64
	Unit         string
	Weight       int64
	Variant      int
	Status       string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    gorm.DeletedAt
}

func (p *Product) TableName() string {
	return "products"
}
