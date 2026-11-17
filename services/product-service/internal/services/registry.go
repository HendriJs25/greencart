package services

import (
	"product-service/internal/config"
	"product-service/internal/repository"
	categoryservice "product-service/internal/services/category"
	jwtservice "product-service/internal/services/jwt"
)

type Registry struct {
	JWTService      jwtservice.Service
	CategoryService categoryservice.Service
}

func NewRegistry(repositories *repository.Registry, cfg *config.Config) *Registry {
	jwtService := jwtservice.NewService(cfg.JWT)
	return &Registry{
		JWTService:      jwtService,
		CategoryService: categoryservice.NewService(repositories.CategoryRepository),
	}
}
