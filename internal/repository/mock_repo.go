package repository

import (
	"context"
	"gophermart/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveUser(ctx context.Context, user *domain.Credentials) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetUser(ctx context.Context, login string) (*domain.UserIDPassword, error) {
	args := m.Called(ctx, login)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*domain.UserIDPassword), args.Error(1)
}

func (m *MockRepository) GetOrderWithUserID(ctx context.Context, number string) (*domain.OrderWithUserID, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(*domain.OrderWithUserID), args.Error(1)
}

func (m *MockRepository) CreateOrder(ctx context.Context, data *domain.OrderWithUserID) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockRepository) GetAllOrders(ctx context.Context, userID int) ([]domain.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Order), args.Error(1)
}

func (m *MockRepository) GetUserBalance(ctx context.Context, userID int) (*domain.Balance, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.Balance), args.Error(1)
}

func (m *MockRepository) BalanceWithdraw(ctx context.Context, userID int, withdraw *domain.OrderToWithdraw) error {
	args := m.Called(ctx, userID, withdraw)
	return args.Error(0)
}

func (m *MockRepository) GetWithdraws(ctx context.Context, userID int) ([]domain.Withdraw, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Withdraw), args.Error(1)
}
