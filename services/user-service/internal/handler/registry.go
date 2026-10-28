package handler

import (
	authhandler "user-service/internal/handler/auth"
	healthhandler "user-service/internal/handler/health"
	userhandler "user-service/internal/handler/user"
	"user-service/internal/services"

	customValidator "github.com/go-playground/validator/v10"
)

type Registry struct {
	HealthHandler *healthhandler.Handler
	AuthHandler   *authhandler.Handler
	UserHandler   *userhandler.Handler
}

func NewRegistry(services *services.Registry, validate *customValidator.Validate) *Registry {
	return &Registry{
		HealthHandler: healthhandler.NewHandler(),
		AuthHandler:   authhandler.NewHandler(services.AuthService, validate),
		UserHandler:   userhandler.NewHandler(),
	}
}
