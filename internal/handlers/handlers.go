package handlers

import (
	"gophermart/internal/domain"
)

type Repository interface {
	SaveUser(user *domain.Credentials) error
	GetUser(login string) (*domain.Credentials, error)
	GetOrderWithUserID(number string) (*domain.OrderWithUserID, error)
	CreateOrder(*domain.OrderWithUserID) error
}

type Handler struct {
	utils *domain.Utils
	repo  Repository
}

func New(
	utils *domain.Utils,
) *Handler {
	return &Handler{
		utils: utils,
	}
}
