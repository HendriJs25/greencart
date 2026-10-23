package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
}

func (r *Role) TableName() string {
	return "roles"
}
