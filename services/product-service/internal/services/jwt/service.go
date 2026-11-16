package jwt

import (
	"fmt"
	"product-service/internal/config"
	errConstant "product-service/internal/constants/error"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type service struct {
	secretKey []byte
	issuer    string
}

type Service interface {
	ValidateAccessToken(accessToken string) (*AccessTokenClaims, error)
}

func NewService(cfg *config.JWT) Service {
	return &service{
		secretKey: []byte(cfg.SecretKey),
		issuer:    cfg.Issuer,
	}
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
