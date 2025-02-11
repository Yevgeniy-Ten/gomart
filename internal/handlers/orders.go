package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) Orders(c *gin.Context) {
	c.JSON(200, "Hello")
}

func (h *Handler) CreateOrder(c *gin.Context) {
	c.JSON(200, "Hello")
}
