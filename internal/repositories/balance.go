package repositories

import (
	"context"
	"fmt"
	"go_test_task/internal/models"

	"gorm.io/gorm"
)

type BalanceRepository interface {
	GetBalance(ctx context.Context, userID uint64) (models.Balance, error)
}

type balanceRepository struct {
	DB *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return balanceRepository{DB: db}
}

func (r balanceRepository) GetBalance(ctx context.Context, userID uint64) (balance models.Balance, err error) {
	if err = r.DB.Where("user_id = ?", userID).First(&balance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = fmt.Errorf("balance for user %d not found", userID)
		}
	}
	return balance, err
}
