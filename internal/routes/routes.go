package routes

import (
	"github.com/gin-gonic/gin"
	"gophermart/internal/domain"
	h "gophermart/internal/handlers"
	"gophermart/internal/handlers/middleware"
	"gophermart/internal/repository"
)

func Init(
	utils *domain.Utils,
	repo *repository.Repo,
) *gin.Engine {
	r := gin.Default()
	handlers := h.New(utils, repo)
	userApi := r.Group("/api/user")
	{
		userApi.POST("/register", handlers.Register)
		userApi.POST("/login", handlers.Login)
		userApi.POST("/orders", middleware.HasUserID(utils.L), handlers.CreateOrder)
		userApi.GET("/orders", middleware.HasUserID(utils.L), handlers.Orders)
		userApi.GET("/balance", middleware.HasUserID(utils.L), handlers.Balance)
		userApi.POST("/balance/withdraw", middleware.HasUserID(utils.L), handlers.BalanceWithdraw)
		userApi.GET("/withdrawals", middleware.HasUserID(utils.L), handlers.Withdrawals)
	}

	return r
}
