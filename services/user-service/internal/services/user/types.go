package user

import (
	"user-service/internal/domain/model"
	"uuid"
)

type GetProfileResult struct {
	RoleName []string
	ID       int64
	UUID     uuid.UUID
	Name     string
	Email    string
	Phone    *string
	Lat      *string
	Lng      *string
	Address  *string
	Photo    *string
}

func ToGetProfileResult(user *model.User) *GetProfileResult {
	roleName := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleName = append(roleName, role.Name)
	}
	return &GetProfileResult{
		RoleName: roleName,
		ID:       user.ID,
		UUID:     user.UUID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Lat:      user.Lat,
		Lng:      user.Lng,
		Address:  user.Address,
		Photo:    user.Photo,
	}
}
