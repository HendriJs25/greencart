package verificationtoken

import (
	"context"
	"errors"
	"fmt"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/domain/model"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	Create(ctx context.Context, verificationToken *model.VerificationToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*model.VerificationToken, error)
	DeleteByID(ctx context.Context, id int64) error
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

func (r *repository) FindByTokenHash(ctx context.Context, tokenHash string) (*model.VerificationToken, error) {
	var verificationToken model.VerificationToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).Take(&verificationToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find verification token: %w", err)
	}
	return &verificationToken, nil
}

func (r *repository) DeleteByID(ctx context.Context, id int64) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.VerificationToken{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete verification token: %w", err)
	}
	return nil
}
