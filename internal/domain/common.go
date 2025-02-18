package domain

import "go.uber.org/zap"

type Config struct {
	Address     string `env:"RUN_ADDRESS"`
	DatabaseURL string `env:"DATABASE_URI"`
	AccrualHost string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JobInterval int
}
type Utils struct {
	L *zap.Logger
	C *Config
}
