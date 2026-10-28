package repository

import (
	rolerepository "user-service/internal/repository/role"
	userrepository "user-service/internal/repository/user"
	verificationtokenrepository "user-service/internal/repository/verificationtoken"

	"gorm.io/gorm"
)

type Registry struct {
	UserRepository              userrepository.Repository
	VerificationTokenRepository verificationtokenrepository.Repository
	RoleRepository              rolerepository.Repository
	Transaction                 TransactionManager
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		UserRepository:              userrepository.NewRepository(db),
		VerificationTokenRepository: verificationtokenrepository.NewRepository(db),
		RoleRepository:              rolerepository.NewRepository(db),
		Transaction:                 NewTransactionManager(db),
	}
}
