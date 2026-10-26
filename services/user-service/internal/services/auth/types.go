package auth

import (
	"time"
	"user-service/internal/domain/model"
	"user-service/internal/services/jwt"
	"uuid"
)

type SignInResult struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
	CSRFToken            string
	UUID                 uuid.UUID
	Name                 string
	Email                string
	Phone                *string
	Lat                  *string
	Lng                  *string
}

type SignInInput struct {
	Email    string
	Password string
}

func ToSignInResult(user *model.User, accessToken *jwt.AccessToken, csrfToken string) *SignInResult {
	return &SignInResult{
		AccessToken:          accessToken.Value,
		AccessTokenExpiresAt: accessToken.ExpiresAt,
		CSRFToken:            csrfToken,
		UUID:                 user.UUID,
		Name:                 user.Name,
		Email:                user.Email,
		Phone:                user.Phone,
		Lat:                  user.Lat,
		Lng:                  user.Lng,
	}
}
