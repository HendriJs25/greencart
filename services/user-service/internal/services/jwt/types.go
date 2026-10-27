package jwt

import (
	"time"
	"uuid"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type AccessToken struct {
	Value     string
	ExpiresAt time.Time
}

type claims struct {
	UserUUID uuid.UUID
	jwtlib.RegisteredClaims
}

type AccessTokenClaims struct {
	UserUUID  uuid.UUID
	SessionID string
	Issuer    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}
