package seeders

import (
	"errors"
	"fmt"
	"user-service/internal/constants"
	"user-service/internal/domain/model"

	"gorm.io/gorm"
)

func seedRoles(tx *gorm.DB) (*model.Role, error) {
	roleSeeds := []*model.Role{
		{
			Name: constants.Admin,
		},
		{
			Name: constants.Customer,
		},
	}

	var roleAdmin *model.Role

	for _, roleSeed := range roleSeeds {
		role, err := findOrCreateRole(tx, roleSeed)
		if err != nil {
			return nil, err
		}
		if role.Name == constants.Admin {
			roleAdmin = role
		}
	}
	return roleAdmin, nil

}

func findOrCreateRole(tx *gorm.DB, roleSeed *model.Role) (*model.Role, error) {
	var role model.Role

	err := tx.Unscoped().Where("name = ?", roleSeed.Name).First(&role).Error

	switch {
	case err == nil:
		if role.DeletedAt.Valid {
			err := tx.Unscoped().Model(&role).Where("name = ?", roleSeed.Name).Update("deleted_at", nil).Error
			if err != nil {
				return nil, fmt.Errorf("failed to restore roles: %w", err)
			}
			role.DeletedAt = gorm.DeletedAt{}
		}
		return &role, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		role = model.Role{
			Name: roleSeed.Name,
		}
		err := tx.Create(&role).Error
		if err != nil {
			return nil, fmt.Errorf("failed to seed role: %w", err)
		}
		return &role, nil
	default:
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
}
