package repositories

import (
	"context"
	"go_test_task/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, entity models.Transaction) error
}

type transactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return transactionRepository{DB: db}
}

func (r transactionRepository) CreateTransaction(ctx context.Context, entity models.Transaction) error {
	return nil
}
