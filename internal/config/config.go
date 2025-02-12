package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"gophermart/internal/domain"
)

func New() (*domain.Config, error) {
	config := &domain.Config{
		Address:     ":8081",
		DatabaseURL: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	}

	parseFlags(config)
	if err := parseEnv(config); err != nil {
		return nil, err
	}

	return config, nil
}

func parseEnv(config *domain.Config) error {
	var envConfig domain.Config
	if err := env.Parse(&envConfig); err != nil {
		return err
	}

	// Обновляем конфигурацию только если переменные окружения заданы
	if envConfig.Address != "" {
		config.Address = envConfig.Address
	}
	return nil
}

func parseFlags(config *domain.Config) {
	flag.StringVar(&config.Address, "a", config.Address, "address for server")
	flag.Parse()
}
