package middleware

import (
	"gophermart/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Authorize(l *zap.Logger, s domain.SessionInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := s.GetUserID(c.Request.Header.Get("Authorization"))
		if err != nil {
			if l != nil {
				l.Debug("error getting user id", zap.Error(err))
			}
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}
