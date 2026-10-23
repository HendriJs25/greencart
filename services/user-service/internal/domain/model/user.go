package model

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type User struct {
	ID           int64
	UUID         uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Phone        *string
	Photo        *string
	Address      *string
	lat          *string
	Lng          *string
	IsVerified   bool `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    gorm.DeletedAt

	Roles []Role `gorm:"many2many:user_roles"`
}

func (u *User) TableName() string {
	return "users"
}
