package services

import (
	"context"
	"fmt"
	"go_test_task/internal/dto"
	"go_test_task/internal/models"
	"go_test_task/internal/repositories"
)

type BalanceService interface {
	GetBalance(ctx context.Context, userID uint64) (dto.Balance, error)
}

type balanceService struct {
	balanceRepository repositories.BalanceRepository
	userRepository    repositories.UserRepository
}

func NewBalanceService(
	balanceRepository repositories.BalanceRepository,
	userRepository repositories.UserRepository,
) BalanceService {
	return balanceService{
		balanceRepository: balanceRepository,
		userRepository:    userRepository,
	}
}

func (s balanceService) GetBalance(ctx context.Context, userID uint64) (dto.Balance, error) {
	exists, err := s.userRepository.Exists(ctx, userID)
	if err != nil {
		return dto.Balance{}, err
	}

	if !exists {
		return dto.Balance{}, fmt.Errorf("user not exists")
	}

	balance, err := s.balanceRepository.GetBalance(ctx, userID)
	if err != nil {
		return dto.Balance{}, err
	}

	return mapBalanceModelToDTO(balance), nil
}

func mapBalanceModelToDTO(model models.Balance) dto.Balance {
	return dto.Balance{
		UserID:  model.UserID,
		Balance: model.Amount,
	}
}
