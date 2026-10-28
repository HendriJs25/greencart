package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type transactionManager struct {
	db *gorm.DB
}

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(repositories *Registry) error) error
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{
		db: db,
	}
}

func (t *transactionManager) WithinTransaction(ctx context.Context, fn func(repositories *Registry) error) error {
	err := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		transactionRepository := NewRegistry(tx)

		return fn(transactionRepository)
	})

	if err != nil {
		return fmt.Errorf("failed to execute database transaction: %w", err)
	}

	return nil
}
