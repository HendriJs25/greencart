package verificationtoken

import (
	"context"
	"fmt"
	"user-service/internal/domain/model"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	Create(ctx context.Context, verificationToken *model.VerificationToken) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, verificationToken *model.VerificationToken) error {
	err := r.db.WithContext(ctx).Create(verificationToken).Error
	if err != nil {
		return fmt.Errorf("failed to create verification token: %w", err)
	}
	return nil
}
