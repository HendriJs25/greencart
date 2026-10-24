package seeders

import (
	"errors"
	"fmt"
	"user-service/internal/common/helper"
	"user-service/internal/config"
	"user-service/internal/domain/model"
	"uuid"

	"gorm.io/gorm"
)

func seedAdmin(tx *gorm.DB, cfg *config.Seed, roles []model.Role) error {
	var admin model.User

	err := tx.Unscoped().Where("email = ?", cfg.Email).First(&admin).Error
	switch {
	case err == nil:
		updates := make(map[string]any)
		if !admin.IsVerified {
			updates["is_verified"] = true
		}

		if admin.DeletedAt.Valid {
			updates["deleted_at"] = nil
		}

		if len(updates) > 0 {
			if err := tx.Unscoped().Model(&admin).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update admin: %w", err)
			}
		}

		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		hashedPassword, err := helper.HashPassword(cfg.Password)
		if err != nil {
			return err
		}
		admin = model.User{
			UUID:         uuid.NewV7(),
			Name:         "Super Admin",
			Email:        cfg.Email,
			PasswordHash: hashedPassword,
			IsVerified:   true,
			Roles:        roles,
		}
		if err := tx.Create(&admin).Error; err != nil {
			return fmt.Errorf("failed to seed admin: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("failed to find admin: %w", err)
	}

}
