package domain

import (
	"go.uber.org/zap"
	"gophermart/internal/utils/session"
)

type Config struct {
	Address     string `env:"RUN_ADDRESS"`
	DatabaseURL string `env:"DATABASE_URI"`
	AccrualHost string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JobInterval int
}
type Utils struct {
	L *zap.Logger
	C *Config
	S *session.Session
}
