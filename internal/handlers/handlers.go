package handlers

import (
	"context"
	"gophermart/internal/domain"
)

type Repository interface {
	SaveUser(ctx context.Context, user *domain.Credentials) error
	GetUser(ctx context.Context, login string) (*domain.UserIDPassword, error)
	GetOrderWithUserID(ctx context.Context, number string) (*domain.OrderWithUserID, error)
	CreateOrder(ctx context.Context, data *domain.OrderWithUserID) error
	GetAllOrders(ctx context.Context, userID int) ([]domain.Order, error)
}

type Handler struct {
	utils *domain.Utils
	repo  Repository
}

func New(
	utils *domain.Utils,
	repo Repository,
) *Handler {
	return &Handler{
		utils: utils,
		repo:  repo,
	}
}
