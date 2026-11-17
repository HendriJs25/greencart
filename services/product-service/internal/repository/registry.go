package repository

import (
	categoryrepository "product-service/internal/repository/category"

	"gorm.io/gorm"
)

type Registry struct {
	CategoryRepository categoryrepository.Repository
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		CategoryRepository: categoryrepository.NewRepository(db),
	}
}
