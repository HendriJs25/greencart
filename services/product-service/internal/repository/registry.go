package repository

import (
	categoryrepository "product-service/internal/repository/category"
	productrepository "product-service/internal/repository/product"

	"gorm.io/gorm"
)

type Registry struct {
	CategoryRepository categoryrepository.Repository
	ProductRepository  productrepository.Repository
	Transaction        TransactionManager
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		CategoryRepository: categoryrepository.NewRepository(db),
		ProductRepository:  productrepository.NewRepository(db),
		Transaction:        NewTransactionManager(db),
	}
}
