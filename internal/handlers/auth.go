package handlers

import (
	"context"
	"errors"
	"gophermart/internal/domain"
	"gophermart/internal/repository"
	"gophermart/internal/utils/bcrypt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) Register(c *gin.Context) {
	var user domain.Credentials
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, "invalid input")
		return
	}
	hashPass, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		h.utils.L.Warn("error hashing password", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}
	user.Password = hashPass
	id, err := h.repo.SaveUser(context.TODO(), &user)
	if err != nil {
		var duplicateError *repository.DuplicateError
		if errors.As(err, &duplicateError) {
			c.Status(http.StatusConflict)
			return
		}
		h.utils.L.Warn("error saving user", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}
	token, err := h.utils.S.CreateToken(id)
	if err != nil {
		h.utils.L.Warn("error creating token", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("Authorization", `Bearer `+token)
	c.Status(http.StatusOK)
}

func (h *Handler) Login(c *gin.Context) {
	var user domain.Credentials
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	storedUser, err := h.repo.GetUser(context.TODO(), user.Login)
	if err != nil {
		h.utils.L.Warn("error getting user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !bcrypt.ComparePasswords(user.Password, storedUser.Password) {
		h.utils.L.Warn("passwords do not match", zap.String("login", user.Login))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := h.utils.S.CreateToken(storedUser.ID)
	if err != nil {
		h.utils.L.Warn("error creating token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("Authorization", `Bearer `+token)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
