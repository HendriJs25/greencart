package jwt

import (
	"fmt"
	"time"
	"user-service/internal/config"
	errConstant "user-service/internal/constants/error"
	"uuid"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type service struct {
	secretKey      []byte
	issuer         string
	accessTokenTTL time.Duration
}

type Service interface {
	GenerateAccessToken(userUUID uuid.UUID, jti uuid.UUID) (*AccessToken, error)
	ValidateAccessToken(accessToken string) (*AccessTokenClaims, error)
}

func NewService(cfg *config.JWT) Service {
	return &service{
		secretKey:      []byte(cfg.SecretKey),
		issuer:         cfg.Issuer,
		accessTokenTTL: cfg.AccessTokenTTl,
	}
}

func (s *service) GenerateAccessToken(userUUID uuid.UUID, jti uuid.UUID) (*AccessToken, error) {
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

func (s *service) ValidateAccessToken(accessToken string) (*AccessTokenClaims, error) {
	tokenClaims := &claims{}

	token, err := jwtlib.ParseWithClaims(
		accessToken,
		tokenClaims,
		func(token *jwtlib.Token) (any, error) {
			return s.secretKey, nil
		},
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}),
		jwtlib.WithIssuer(s.issuer),
		jwtlib.WithExpirationRequired(),
		jwtlib.WithIssuedAt(),
	)

	if err != nil {
		return nil, fmt.Errorf("%w: validate access token: %v", errConstant.ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("%w: access token is invalid", errConstant.ErrInvalidToken)
	}

	return &AccessTokenClaims{
		UserUUID:  tokenClaims.UserUUID,
		SessionID: tokenClaims.ID,
		Issuer:    tokenClaims.Issuer,
		IssuedAt:  tokenClaims.IssuedAt.Time,
		ExpiresAt: tokenClaims.ExpiresAt.Time,
	}, nil
}
