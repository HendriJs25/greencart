package handler

import (
	categoryhandler "product-service/internal/handler/category"
	healthhandler "product-service/internal/handler/health"
	producthandler "product-service/internal/handler/product"
	"product-service/internal/services"

	customValidator "github.com/go-playground/validator/v10"
)

type Registry struct {
	HealthHandler   *healthhandler.Handler
	CategoryHandler *categoryhandler.Handler
	ProductHandler  *producthandler.Handler
}

func NewRegistry(services *services.Registry, validate *customValidator.Validate) *Registry {
	return &Registry{
		HealthHandler:   healthhandler.NewHandler(),
		CategoryHandler: categoryhandler.NewHandler(services.CategoryService, validate),
		ProductHandler:  producthandler.NewHandler(services.ProductService, validate),
	}
}
