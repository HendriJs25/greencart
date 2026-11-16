package services

import (
	"product-service/internal/config"
	jwtservice "product-service/internal/services/jwt"
)

type Registry struct {
	JWTService jwtservice.Service
}

func NewRegistry(cfg *config.Config) *Registry {
	jwtService := jwtservice.NewService(cfg.JWT)
	return &Registry{
		JWTService: jwtService,
	}
}
