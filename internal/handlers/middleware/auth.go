package middleware

import (
	"gophermart/internal/utils/session"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HasUserID(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := session.GetUserID(c.Request.Header.Get("Authorization"))
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
