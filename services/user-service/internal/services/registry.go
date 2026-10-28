package services

import (
	"user-service/internal/config"
	"user-service/internal/repository"
	sessionrepository "user-service/internal/repository/session"
	authservice "user-service/internal/services/auth"
	jwtservice "user-service/internal/services/jwt"
)

type Registry struct {
	JWTService  jwtservice.Service
	AuthService authservice.Service
}

func NewRegistry(repositories *repository.Registry, sessionRepository sessionrepository.Repository, cfg *config.JWT) *Registry {
	jwtService := jwtservice.NewService(cfg)
	return &Registry{
		JWTService:  jwtService,
		AuthService: authservice.NewService(jwtService, repositories.UserRepository, sessionRepository, repositories.Transaction),
	}
}
