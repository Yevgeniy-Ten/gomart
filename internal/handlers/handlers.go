package handlers

import "gophermart/internal/utils"

type Handler struct {
	utils *utils.Utils
}

func New(
	utils *utils.Utils,
) *Handler {
	return &Handler{
		utils: utils,
	}
}
