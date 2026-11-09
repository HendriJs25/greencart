package services

import (
	"user-service/internal/config"
	"user-service/internal/messaging/kafka"
	"user-service/internal/repository"
	sessionrepository "user-service/internal/repository/session"
	authservice "user-service/internal/services/auth"
	customerservice "user-service/internal/services/customer"
	jwtservice "user-service/internal/services/jwt"
	userservice "user-service/internal/services/user"
	"user-service/internal/storage"
)

type Registry struct {
	JWTService      jwtservice.Service
	AuthService     authservice.Service
	UserService     userservice.Service
	CustomerService customerservice.Service
}

func NewRegistry(repositories *repository.Registry, sessionRepository sessionrepository.Repository, kafkaProducer kafka.Producer, s3Storage storage.Storage, cfg *config.Config) *Registry {
	jwtService := jwtservice.NewService(cfg.JWT)
	return &Registry{
		JWTService:      jwtService,
		AuthService:     authservice.NewService(jwtService, repositories.UserRepository, sessionRepository, repositories.VerificationTokenRepository, repositories.Transaction, kafkaProducer, cfg.URL),
		UserService:     userservice.NewService(repositories.UserRepository, s3Storage),
		CustomerService: customerservice.NewService(repositories.CustomerRepository, repositories.UserRepository, repositories.Transaction, kafkaProducer),
	}
}
