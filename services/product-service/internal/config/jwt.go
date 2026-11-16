package config

import (
	"fmt"
)

const (
	minimumJWTSecretBytes = 32
)

type JWT struct {
	SecretKey string
	Issuer    string
}

func loadJWT() (*JWT, error) {
	secretKey, err := getEnvRequired("JWT_SECRET_KEY")
	if err != nil {
		return nil, err
	}

	issuer, err := getEnvRequired("JWT_ISSUER")
	if err != nil {
		return nil, err
	}

	return &JWT{
		SecretKey: secretKey,
		Issuer:    issuer,
	}, nil
}

func (j *JWT) validate() error {
	if len(j.SecretKey) < minimumJWTSecretBytes {
		return fmt.Errorf("JWT_SECRET_KEY must be at least %d bytes", minimumJWTSecretBytes)
	}

	return nil
}
