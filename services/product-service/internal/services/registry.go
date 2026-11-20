package services

import (
	"product-service/internal/config"
	"product-service/internal/repository"
	categoryservice "product-service/internal/services/category"
	jwtservice "product-service/internal/services/jwt"
	productservice "product-service/internal/services/product"
	"product-service/internal/storage"
)

type Registry struct {
	JWTService      jwtservice.Service
	CategoryService categoryservice.Service
	ProductService  productservice.Service
}

func NewRegistry(repositories *repository.Registry, s3Storage storage.Storage, cfg *config.Config) *Registry {
	jwtService := jwtservice.NewService(cfg.JWT)
	return &Registry{
		JWTService:      jwtService,
		CategoryService: categoryservice.NewService(repositories.CategoryRepository, s3Storage),
		ProductService:  productservice.NewService(repositories.ProductRepository, repositories.CategoryRepository),
	}
}
