package user

import (
	"user-service/internal/domain/model"
	"uuid"
)

type GetMeResult struct {
	RoleName []string
	UUID     uuid.UUID
	Name     string
	Email    string
}
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

type UpdateUserDataInput struct {
	Name    string
	Phone   int64
	Lat     *float64
	Lng     *float64
	Address *string
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

func ToGetMeResult(user *model.User) *GetMeResult {
	roleName := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleName = append(roleName, role.Name)
	}
	return &GetMeResult{
		RoleName: roleName,
		UUID:     user.UUID,
		Name:     user.Name,
		Email:    user.Email,
	}
}
