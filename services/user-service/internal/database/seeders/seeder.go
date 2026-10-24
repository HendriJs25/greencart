package seeders

import (
	"context"
	"user-service/internal/config"
	"user-service/internal/domain/model"

	"gorm.io/gorm"
)

func Run(ctx context.Context, db *gorm.DB, cfg *config.Seed) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		adminRole, err := seedRoles(tx)
		if err != nil {
			return err
		}

		err = seedAdmin(tx, cfg, []model.Role{*adminRole})
		if err != nil {
			return err
		}
		return nil
	})
}
