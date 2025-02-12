package handlers

import (
	"gophermart/internal/domain"
)

type Handler struct {
	utils *domain.Utils
}

func New(
	utils *domain.Utils,
) *Handler {
	return &Handler{
		utils: utils,
	}
}
