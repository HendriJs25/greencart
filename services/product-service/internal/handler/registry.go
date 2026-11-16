package handler

import healthhandler "product-service/internal/handler/health"

type Registry struct {
	HealthHandler *healthhandler.Handler
}

func NewRegistry() *Registry {
	return &Registry{
		HealthHandler: healthhandler.NewHandler(),
	}
}
