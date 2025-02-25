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
	userAPI.POST("/orders", middleware.HasUserID(utils.L), handlers.CreateOrder)
	userAPI.GET("/orders", middleware.HasUserID(utils.L), handlers.Orders)
	userAPI.GET("/balance", middleware.HasUserID(utils.L), handlers.Balance)
	userAPI.POST("/balance/withdraw", middleware.HasUserID(utils.L), handlers.BalanceWithdraw)
	userAPI.GET("/withdrawals", middleware.HasUserID(utils.L), handlers.Withdrawals)
	return r
}
