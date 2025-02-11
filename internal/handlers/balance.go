package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) Balance(c *gin.Context) {
	c.JSON(200, "Hello")
}

func (h *Handler) BalanceWithdraw(c *gin.Context) {
	c.JSON(200, "Hello")
}
func (h *Handler) Withdrawals(c *gin.Context) {
	c.JSON(200, "Hello")
}
