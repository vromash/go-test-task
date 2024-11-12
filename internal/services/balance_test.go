package services

import (
	"context"
	"fmt"
	"go_test_task/internal/dto"
	mock_repositories "go_test_task/internal/mocks/repositories"
	"go_test_task/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type balanceServiceTestSuite struct {
	suite.Suite
	ctrl   *gomock.Controller
	logger *zap.Logger

	balanceRepository *mock_repositories.MockBalanceRepository
	userRepository    *mock_repositories.MockUserRepository

	service balanceService
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(balanceServiceTestSuite))
}

func (s *balanceServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.logger = zap.NewNop()

	s.balanceRepository = mock_repositories.NewMockBalanceRepository(s.ctrl)
	s.userRepository = mock_repositories.NewMockUserRepository(s.ctrl)

	s.service = NewBalanceService(s.logger, s.balanceRepository, s.userRepository).(balanceService)
}

func (s *balanceServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *balanceServiceTestSuite) TestGetBalance() {
	ctx := context.Background()
	expectedUserID := uint64(1)

	s.Run("balance retrieved", func() {
		expected := dto.Balance{
			UserID:        1,
			CurrentAmount: "10.50",
		}

		s.userRepository.EXPECT().Exists(ctx, expectedUserID).Return(true, nil)
		s.balanceRepository.EXPECT().GetBalance(ctx, expectedUserID).Return(models.Balance{
			UserID: 1,
			Amount: "10.50",
		}, nil)

		actual, err := s.service.GetBalance(ctx, 1)
		s.NoError(err)
		s.Equal(expected, actual)
	})

	s.Run("error occurs while getting balance", func() {
		expected := dto.Balance{}
		expectedErr := fmt.Errorf("failed to get balance: db failed")

		s.userRepository.EXPECT().Exists(ctx, expectedUserID).Return(true, nil)
		s.balanceRepository.EXPECT().GetBalance(ctx, expectedUserID).Return(models.Balance{}, fmt.Errorf("db failed"))

		actual, err := s.service.GetBalance(ctx, 1)
		s.Equal(expectedErr, err)
		s.Equal(expected, actual)
	})
}

func (s *balanceServiceTestSuite) TestUpdateBalance() {
	ctx := context.Background()
	expectedUserID := uint64(1)
	defaultReq := dto.UpdateBalance{
		UserID:        expectedUserID,
		Source:        "game",
		Amount:        "10.50",
		ParsedAmount:  10.50,
		TransactionID: "111-222-333-444",
	}

	s.Run("balance increased", func() {
		req := defaultReq
		req.State = "win"

		s.userRepository.EXPECT().Exists(ctx, expectedUserID).Return(true, nil)
		s.balanceRepository.EXPECT().IncreaseBalance(ctx, models.Transaction{
			UID:    "111-222-333-444",
			UserID: expectedUserID,
			Source: "game",
			TState: "win",
			Amount: "10.50",
		}, 10.50).Return(nil)

		err := s.service.UpdateBalance(ctx, req)
		s.NoError(err)
	})

	s.Run("balance decreased", func() {
		req := defaultReq
		req.State = "lose"

		s.userRepository.EXPECT().Exists(ctx, expectedUserID).Return(true, nil)
		s.balanceRepository.EXPECT().DecreaseBalance(ctx, models.Transaction{
			UID:    "111-222-333-444",
			UserID: expectedUserID,
			Source: "game",
			TState: "lose",
			Amount: "10.50",
		}, 10.50).Return(nil)

		err := s.service.UpdateBalance(ctx, req)
		s.NoError(err)
	})

	s.Run("error occurs while increasing balance", func() {
		req := defaultReq
		req.State = "win"
		expectedErr := fmt.Errorf("failed to update balance: db failed")

		s.userRepository.EXPECT().Exists(ctx, expectedUserID).Return(true, nil)
		s.balanceRepository.EXPECT().IncreaseBalance(ctx, models.Transaction{
			UID:    "111-222-333-444",
			UserID: expectedUserID,
			Source: "game",
			TState: "win",
			Amount: "10.50",
		}, 10.50).Return(fmt.Errorf("db failed"))

		err := s.service.UpdateBalance(ctx, req)
		s.Equal(expectedErr, err)
	})

	s.Run("error occurs while decreasing balance", func() {
		req := defaultReq
		req.State = "lose"
		expectedErr := fmt.Errorf("failed to update balance: db failed")

		s.userRepository.EXPECT().Exists(ctx, expectedUserID).Return(true, nil)
		s.balanceRepository.EXPECT().DecreaseBalance(ctx, models.Transaction{
			UID:    "111-222-333-444",
			UserID: expectedUserID,
			Source: "game",
			TState: "lose",
			Amount: "10.50",
		}, 10.50).Return(fmt.Errorf("db failed"))

		err := s.service.UpdateBalance(ctx, req)
		s.Equal(expectedErr, err)
	})
}

func (s *balanceServiceTestSuite) Test_checkIfUserExists() {
	ctx := context.Background()
	expectedUserID := uint64(1)

	s.Run("balance retrieved", func() {
		s.userRepository.EXPECT().Exists(ctx, expectedUserID).Return(true, nil)

		err := s.service.checkIfUserExists(ctx, 1)
		s.NoError(err)
	})

	s.Run("error occurs when user does not exist", func() {
		expectedErr := fmt.Errorf("user does not exist")

		s.userRepository.EXPECT().Exists(ctx, expectedUserID).Return(false, nil)

		err := s.service.checkIfUserExists(ctx, 1)
		s.Equal(expectedErr, err)
	})

	s.Run("error occurs while checking if user exists", func() {
		expectedErr := fmt.Errorf("failed to find user: db failed")

		s.userRepository.EXPECT().Exists(ctx, expectedUserID).Return(false, fmt.Errorf("db failed"))

		err := s.service.checkIfUserExists(ctx, 1)
		s.Equal(expectedErr, err)
	})
}

func (s *balanceServiceTestSuite) Test_mapUpdateBalanceDTOToTransactionModel() {
	expectedUserID := uint64(1)
	defaultReq := dto.UpdateBalance{
		UserID:        expectedUserID,
		Amount:        "10.50",
		ParsedAmount:  10.50,
		TransactionID: "111-222-333-444",
	}

	s.Run("request mapped", func() {
		req := defaultReq
		req.Source = "game"
		req.State = "win"
		expected := models.Transaction{
			UID:    "111-222-333-444",
			UserID: expectedUserID,
			Source: "game",
			TState: "win",
			Amount: "10.50",
		}

		actual, err := mapUpdateBalanceDTOToTransactionModel(req)
		s.NoError(err)
		s.Equal(expected, actual)
	})

	s.Run("error occurs while parsing transaction state", func() {
		req := defaultReq
		req.Source = "game"
		expected := models.Transaction{}
		expectedErr := fmt.Errorf("invalid transaction state value")

		actual, err := mapUpdateBalanceDTOToTransactionModel(req)
		s.Equal(expectedErr, err)
		s.Equal(expected, actual)
	})

	s.Run("error occurs while parsing transaction source", func() {
		req := defaultReq
		req.State = "win"
		expected := models.Transaction{}
		expectedErr := fmt.Errorf("invalid transaction source value")

		actual, err := mapUpdateBalanceDTOToTransactionModel(req)
		s.Equal(expectedErr, err)
		s.Equal(expected, actual)
	})
}
