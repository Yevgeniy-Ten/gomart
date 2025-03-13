package domain

import (
	"go.uber.org/zap"
)

type Config struct {
	Address     string `env:"RUN_ADDRESS"`
	DatabaseURL string `env:"DATABASE_URI"`
	AccrualHost string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JobInterval int
}
type SessionInterface interface {
	CreateToken(userID int) (string, error)
	GetUserID(authHeader string) (int, error)
}
type Utils struct {
	L *zap.Logger
	C *Config
	S SessionInterface
}
