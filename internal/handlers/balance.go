package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gophermart/internal/domain"
	"gophermart/internal/repository"
	"gophermart/internal/utils/session"
	"net/http"
	"time"
)

func (h *Handler) Balance(c *gin.Context) {
	time.Sleep(3 * time.Second)
	requestUserID, _ := session.GetUserID(c.Request.Header.Get("Authorization"))
	balance, err := h.repo.GetUserBalance(context.TODO(), requestUserID)
	if err != nil {
		h.utils.L.Warn("error getting balance", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (h *Handler) BalanceWithdraw(c *gin.Context) {
	requestUserID, err := session.GetUserID(c.Request.Header.Get("Authorization"))
	if err != nil {
		h.utils.L.Warn("error getting user id", zap.Error(err))
		c.Status(http.StatusUnauthorized)
		return
	}
	var orderToWithdraw domain.OrderToWithdraw
	if err := c.BindJSON(&orderToWithdraw); err != nil {
		h.utils.L.Warn("error binding json", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}
	if orderToWithdraw.Sum <= 0 {
		h.utils.L.Warn("error sum", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}
	err = h.repo.BalanceWithdraw(context.TODO(), requestUserID, &orderToWithdraw)
	if err != nil {
		var shouldPositiveErr *repository.ShouldBePositiveError
		if errors.As(err, &shouldPositiveErr) {
			c.Status(http.StatusPaymentRequired)
			return
		}
		h.utils.L.Warn("error withdraw balance", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
func (h *Handler) Withdrawals(c *gin.Context) {
	requestUserID, err := session.GetUserID(c.Request.Header.Get("Authorization"))
	if err != nil {
		h.utils.L.Warn("error getting user id", zap.Error(err))
		c.Status(http.StatusUnauthorized)
		return
	}
	withdraws, err := h.repo.GetWithdraws(context.TODO(), requestUserID)
	if err != nil {
		h.utils.L.Warn("error getting withdraws", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(withdraws) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, withdraws)
}
