package routes

import (
	"github.com/gin-gonic/gin"
	h "gophermart/internal/handlers"
	"gophermart/internal/utils"
)

func Init(
	utils *utils.Utils,
) *gin.Engine {
	r := gin.Default()
	handlers := h.New(utils)
	userApi := r.Group("/api/user")
	{
		userApi.POST("/register", handlers.Register)
		userApi.POST("/login", handlers.Login)
		userApi.POST("/orders", handlers.CreateOrder)
		userApi.GET("/orders", handlers.Orders)
		userApi.GET("/balance", handlers.Balance)
		userApi.POST("/balance/withdraw", handlers.BalanceWithdraw)
		userApi.GET("/withdrawals", handlers.Withdrawals)
	}

	return r
}
