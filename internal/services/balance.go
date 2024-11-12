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
	UpdateBalance(ctx context.Context, req dto.UpdateBalance) error
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
	if err := s.checkIfUserExists(ctx, userID); err != nil {
		return dto.Balance{}, err
	}

	balance, err := s.balanceRepository.GetBalance(ctx, userID)
	if err != nil {
		return dto.Balance{}, fmt.Errorf("failed to get balance: %v", err)
	}

	return mapBalanceModelToDTO(balance), nil
}

func (s balanceService) UpdateBalance(ctx context.Context, req dto.UpdateBalance) error {
	if err := s.checkIfUserExists(ctx, req.UserID); err != nil {
		return err
	}

	var err error

	entity, err := mapUpdateBalanceDTOToTransactionModel(req)
	if err != nil {
		return err
	}

	if req.State == string(models.Win) {
		err = s.balanceRepository.IncreaseBalance(ctx, entity, req.ParsedAmount)
	}
	if req.State == string(models.Lose) {
		err = s.balanceRepository.DecreaseBalance(ctx, entity, req.ParsedAmount)
	}

	if err != nil {
		return fmt.Errorf("failed to update balance: %v", err)
	}

	return nil
}

func (s balanceService) checkIfUserExists(ctx context.Context, userID uint64) error {
	exists, err := s.userRepository.Exists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}

	if !exists {
		return fmt.Errorf("user does not exist")
	}

	return nil
}

func mapBalanceModelToDTO(model models.Balance) dto.Balance {
	return dto.Balance{
		UserID:        model.UserID,
		CurrentAmount: model.Amount,
	}
}

func mapUpdateBalanceDTOToTransactionModel(req dto.UpdateBalance) (models.Transaction, error) {
	state, err := models.StringToTransactionState(req.State)
	if err != nil {
		return models.Transaction{}, err
	}

	source, err := models.StringToTransactionSource(req.Source)
	if err != nil {
		return models.Transaction{}, err
	}

	return models.Transaction{
		UID:    req.TransactionID,
		UserID: req.UserID,
		Source: source,
		TState: state,
		Amount: req.Amount,
	}, nil
}
