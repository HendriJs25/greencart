package model

import (
	"time"

	"gorm.io/gorm"
)

type VerificationToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	TokenType string
	ExpireAt  time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt

	User User `gorm:"foreignkey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}
