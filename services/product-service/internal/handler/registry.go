package handler

import (
	categoryhandler "product-service/internal/handler/category"
	healthhandler "product-service/internal/handler/health"
	"product-service/internal/services"
)

type Registry struct {
	HealthHandler   *healthhandler.Handler
	CategoryHandler *categoryhandler.Handler
}

func NewRegistry(services *services.Registry) *Registry {
	return &Registry{
		HealthHandler:   healthhandler.NewHandler(),
		CategoryHandler: categoryhandler.NewHandler(services.CategoryService),
	}
}
