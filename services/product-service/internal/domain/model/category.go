package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID           int64
	ParentID     *int64
	Name         string
	Icon         *string
	Status       bool
	Slug         string
	Description  *string
	TotalProduct int64 `gorm:"->;-:migration"`
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *gorm.DeletedAt
	Products     []Product `gorm:"foreignkey:CategorySlug;references:Slug;"`

	Parent *Category `gorm:"foreignkey:ParentID;references:ID"`
}

func (c *Category) TableName() string {
	return "categories"
}
