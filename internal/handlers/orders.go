package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gophermart/internal/domain"
	"io"
	"net/http"
)

func (h *Handler) Orders(c *gin.Context) {
	c.JSON(200, "Hello")
}

func (h *Handler) CreateOrder(c *gin.Context) {
	//middleware logic 401
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(400, "Read error")
		return
	}
	orderNum := string(body)

	if orderNum == "" {
		c.Status(400)
		return
	}
	existOrder, err := h.repo.GetOrderWithUserID(orderNum)
	if err != nil {
		h.utils.L.Warn("error getting order", zap.Error(err))
		c.Status(500)
		return
	}
	//check user id in request

	//after create order

	err = h.repo.CreateOrder(&domain.OrderWithUserID{
		Number: orderNum,
		UserID: 1,
	})

	if err != nil {
		h.utils.L.Error("error creating order", zap.Error(err))
		c.Status(500)
		return
	}

	c.JSON(http.StatusAccepted, "Hello")
}
