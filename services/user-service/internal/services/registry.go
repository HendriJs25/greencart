package services

import (
	"user-service/internal/repository"
	authservice "user-service/internal/services/auth"
)

type Registry struct {
	AuthService authservice.Service
}

func NewRegistry(repositories *repository.Registry) *Registry {
	return &Registry{
		AuthService: authservice.NewService(repositories.UserRepository),
	}
}
