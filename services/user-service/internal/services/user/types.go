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
	Lat      *float64
	Lng      *float64
	Address  *string
	PhotoURL *string
}

func ToGetProfileResult(user *model.User, photoURL *string) *GetProfileResult {
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
		PhotoURL: photoURL,
	}
}
