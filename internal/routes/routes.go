package routes

import (
	"gophermart/internal/domain"
	h "gophermart/internal/handlers"
	"gophermart/internal/handlers/middleware"
	"gophermart/internal/repository"

	"github.com/gin-gonic/gin"
)

func Init(
	utils *domain.Utils,
	repo *repository.Repo,
) *gin.Engine {
	r := gin.Default()
	handlers := h.New(utils, repo)
	userAPI := r.Group("/api/user")
	userAPI.POST("/register", handlers.Register)
	userAPI.POST("/login", handlers.Login)
	userAPI.POST("/orders", middleware.Authorize(utils.L, utils.S), handlers.CreateOrder)
	userAPI.GET("/orders", middleware.Authorize(utils.L, utils.S), handlers.ListOrders)
	userAPI.GET("/balance", middleware.Authorize(utils.L, utils.S), handlers.Balance)
	userAPI.POST("/balance/withdraw", middleware.Authorize(utils.L, utils.S), handlers.BalanceWithdraw)
	userAPI.GET("/withdrawals", middleware.Authorize(utils.L, utils.S), handlers.Withdrawals)
	return r
}
