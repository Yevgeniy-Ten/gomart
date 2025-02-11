package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) Register(c *gin.Context) {
	c.JSON(200, "Hello")
}

func (h *Handler) Login(c *gin.Context) {
	c.JSON(200, "Hello")
}
