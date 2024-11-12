package repositories

import (
	"context"
	"fmt"
	"go_test_task/internal/models"
	"strconv"

	"gorm.io/gorm"
)

type BalanceRepository interface {
	GetBalance(ctx context.Context, userID uint64) (models.Balance, error)
	IncreaseBalance(ctx context.Context, transaction models.Transaction, amount float64) error
	DecreaseBalance(ctx context.Context, transaction models.Transaction, amount float64) error
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

func (r balanceRepository) IncreaseBalance(ctx context.Context, transaction models.Transaction, amount float64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := r.createTransaction(tx, transaction); err != nil {
			return err
		}

		currentBalance, err := r.getCurrentBalance(tx, transaction)
		if err != nil {
			return err
		}
		currentBalance += amount

		return r.saveBalance(tx, transaction, currentBalance)
	})

}

func (r balanceRepository) DecreaseBalance(ctx context.Context, transaction models.Transaction, amount float64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		currentBalance, err := r.getCurrentBalance(tx, transaction)
		if err != nil {
			return err
		}

		if currentBalance < amount {
			return fmt.Errorf("balance is less then requested amount")
		}

		if err := r.createTransaction(tx, transaction); err != nil {
			return err
		}

		currentBalance -= amount

		return r.saveBalance(tx, transaction, currentBalance)
	})
}

func (r balanceRepository) createTransaction(tx *gorm.DB, transaction models.Transaction) error {
	var userExists int64
	err := tx.Model(&models.Transaction{}).
		Where("uid = ?", transaction.UID).
		Count(&userExists).
		Error

	if err != nil {
		return fmt.Errorf("failed to check if transaction exists: %v", err)
	}

	if userExists != 0 {
		return fmt.Errorf("transaction with id %s already exists", transaction.UID)
	}

	if err := tx.Create(&transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	return nil
}

func (r balanceRepository) getCurrentBalance(tx *gorm.DB, transaction models.Transaction) (float64, error) {
	balance, err := r.GetBalance(context.Background(), transaction.UserID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user balance: %v", err)
	}

	amount, err := strconv.ParseFloat(balance.Amount, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse user balance: %v", err)
	}

	return amount, nil
}

func (r balanceRepository) saveBalance(tx *gorm.DB, transaction models.Transaction, currentBalance float64) error {
	err := r.DB.Model(&models.Balance{}).
		Where("user_id = ?", transaction.UserID).
		Update("amount", strconv.FormatFloat(currentBalance, 'f', 2, 64)).
		Error
	if err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	return nil
}
