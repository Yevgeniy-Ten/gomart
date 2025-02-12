package domain

import "go.uber.org/zap"

type Config struct {
	Address     string `env:"SERVER_ADDRESS"`
	DatabaseURL string `env:"DATABASE_URI"`
}
type Utils struct {
	L *zap.Logger
	C *Config
}
