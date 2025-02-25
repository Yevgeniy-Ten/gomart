package config

import (
	"errors"
	"flag"
	"gophermart/internal/domain"

	"github.com/caarlos0/env/v11"
)

func New() (*domain.Config, error) {
	config := &domain.Config{
		Address:     ":8080",
		DatabaseURL: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		JobInterval: 5,
		AccrualHost: "http://localhost:8081",
	}

	parseFlags(config)
	if err := parseEnv(config); err != nil {
		return nil, err
	}
	if config.DatabaseURL == "" || config.AccrualHost == "" || config.Address == "" {
		return nil, errors.New("no settings provided")
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
	if envConfig.DatabaseURL != "" {
		config.DatabaseURL = envConfig.DatabaseURL
	}
	if envConfig.AccrualHost != "" {
		config.AccrualHost = envConfig.AccrualHost
	}
	return nil
}

func parseFlags(config *domain.Config) {
	flag.StringVar(&config.Address, "a", config.Address, "address for server")
	flag.Parse()
}
