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
	Lat                  *float64
	Lng                  *float64
	Role                 []string
}

type SignInInput struct {
	Email    string
	Password string
}

type SignUpInput struct {
	Name     string
	Email    string
	Password string
}

type ResetPasswordInput struct {
	Token       string
	NewPassword string
}

func ToSignInResult(user *model.User, accessToken *jwt.AccessToken, csrfToken string) *SignInResult {
	roleName := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleName = append(roleName, role.Name)
	}
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
		Role:                 roleName,
	}
}
