package auth

import (
	"user-service/internal/domain/model"
	"uuid"
)

type SignInResult struct {
	CSRFToken string
	UUID      uuid.UUID
	Name      string
	Email     string
	Phone     *string
	Lat       *string
	Lng       *string
}

type SignInInput struct {
	Email    string
	Password string
}

func ToSignInResult(user *model.User, CSRFToken string) *SignInResult {
	return &SignInResult{
		CSRFToken: CSRFToken,
		UUID:      user.UUID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Lat:       user.Lat,
		Lng:       user.Lng,
	}
}
