package handler

import healthhandler "user-service/internal/handler/health"

type Registry struct {
	HealthHandler *healthhandler.Handler
}

func NewRegistry() *Registry {
	return &Registry{
		HealthHandler: healthhandler.NewHandler(),
	}
}
