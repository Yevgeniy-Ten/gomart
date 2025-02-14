package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gophermart/internal/domain"
	"gophermart/internal/repository"
	"gophermart/internal/utils/lunhchecker"
	"gophermart/internal/utils/session"
	"io"
	"net/http"
)

func (h *Handler) Orders(c *gin.Context) {
	requestUserID, err := session.GetUserID(c.Request.Header.Get("Authorization"))
	if err != nil {
		h.utils.L.Warn("error getting user id", zap.Error(err))
		c.Status(http.StatusUnauthorized)
		return
	}

	allOrders, err := h.repo.GetAllOrders(context.TODO(), requestUserID)
	if err != nil {
		h.utils.L.Warn("error getting allOrders", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, allOrders)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	requestUserID, err := session.GetUserID(c.Request.Header.Get("Authorization"))
	if err != nil {
		h.utils.L.Warn("error getting user id", zap.Error(err))
		c.Status(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Read error")
		return
	}
	orderNum := string(body)
	if orderNum == "" || !lunhchecker.LuhnCheck(orderNum) {
		c.Status(400)
		return
	}
	existOrder, err := h.repo.GetOrderWithUserID(context.TODO(), orderNum)
	if err != nil {
		var notFoundError *repository.NotFoundError
		if !errors.As(err, &notFoundError) {
			h.utils.L.Warn("error getting order", zap.Error(err))
			c.Status(500)
			return
		}
		err = h.repo.CreateOrder(context.TODO(), &domain.OrderWithUserID{
			Number: orderNum,
			UserID: requestUserID,
		})

		if err != nil {
			h.utils.L.Error("error creating order", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusAccepted)
		return
	}
	if requestUserID != existOrder.UserID {
		c.Status(http.StatusConflict)
	}
	c.Status(http.StatusOK)
}
