package repository

import (
	userrepository "user-service/internal/repository/user"

	"gorm.io/gorm"
)

type Registry struct {
	UserRepository userrepository.Repository
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		UserRepository: userrepository.NewRepository(db),
	}
}
