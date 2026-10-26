package jwt

import (
	"fmt"
	"time"
	"user-service/internal/config"
	"uuid"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type service struct {
	secretKey      []byte
	issuer         string
	accessTokenTTL time.Duration
}

type Service interface {
	GenerateAccessToken(userUUID uuid.UUID) (*AccessToken, error)
}

func NewService(cfg *config.JWT) Service {
	return &service{
		secretKey:      []byte(cfg.SecretKey),
		issuer:         cfg.Issuer,
		accessTokenTTL: cfg.AccessTokenTTl,
	}
}

func (s *service) GenerateAccessToken(userUUID uuid.UUID) (*AccessToken, error) {
	jti := uuid.NewV4()

	now := time.Now().UTC()
	expiresAt := now.Add(s.accessTokenTTL)

	tokenClaims := claims{
		UserUUID: userUUID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ID:        jti.String(),
			Issuer:    s.issuer,
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	return &AccessToken{
		Value:     tokenString,
		ExpiresAt: expiresAt,
	}, nil
}
