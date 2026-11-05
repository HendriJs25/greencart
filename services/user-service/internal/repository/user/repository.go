package user

import (
	"context"
	"errors"
	"fmt"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/domain/model"
	"uuid"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindByID(ctx context.Context, userID int64) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUUID(ctx context.Context, userUUID uuid.UUID) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	IsEmailExist(ctx context.Context, email string) (bool, error)
	MarkAsVerified(ctx context.Context, userID int64) error
	UpdateHashedPasswordByID(ctx context.Context, userID int64, hashedPassword string) error
	UpdatePhotoKeyByUUID(ctx context.Context, userUUID uuid.UUID, photoKey string) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByID(ctx context.Context, userID int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ? AND is_verified = ?", userID, true).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("email = ? AND is_verified = ?", email, true).Preload("Roles").Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

func (r *repository) FindByUUID(ctx context.Context, userUUID uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Roles").Where("uuid = ? AND is_verified = ?", userUUID, true).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by uuid: %w", err)
	}
	return &user, nil
}

func (r *repository) Create(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errConstant.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *repository) IsEmailExist(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}
	return count > 0, nil
}

func (r *repository) MarkAsVerified(ctx context.Context, userID int64) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("is_verified", true)
	if result.Error != nil {
		return fmt.Errorf("failed to mark as verified: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errConstant.ErrUserNotFound
	}
	return nil
}

func (r *repository) UpdateHashedPasswordByID(ctx context.Context, userID int64, hashedPassword string) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("password_hash", hashedPassword)
	if result.Error != nil {
		return fmt.Errorf("failed to update hashed password: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrUserNotFound
	}
	return nil
}

func (r *repository) UpdatePhotoKeyByUUID(ctx context.Context, userUUID uuid.UUID, photoKey string) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("uuid = ?", userUUID).Update("photo", photoKey)
	if result.Error != nil {
		return fmt.Errorf("failed to update photo key: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrUserNotFound
	}
	return nil
}
