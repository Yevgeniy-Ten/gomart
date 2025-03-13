package config

import (
	"gophermart/internal/domain"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_DefaultValues(t *testing.T) {
	cfg, err := New()
	require.NoError(t, err)
	assert.Equal(t, ":8080", cfg.Address)
	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", cfg.DatabaseURL)
	assert.Equal(t, "http://localhost:8081", cfg.AccrualHost)
}
func TestParseEnv(t *testing.T) {
	os.Setenv("RUN_ADDRESS", ":9090")
	os.Setenv("DATABASE_URI", "postgres://test:test@localhost:5432/test")
	os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://test-host:8082")

	defer func() {
		os.Unsetenv("RUN_ADDRESS")
		os.Unsetenv("DATABASE_URI")
		os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
	}()

	cfg := &domain.Config{
		Address:     ":8080",
		DatabaseURL: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		AccrualHost: "http://localhost:8081",
	}

	err := parseEnv(cfg)
	require.NoError(t, err)

	assert.Equal(t, ":9090", cfg.Address)
	assert.Equal(t, "postgres://test:test@localhost:5432/test", cfg.DatabaseURL)
	assert.Equal(t, "http://test-host:8082", cfg.AccrualHost)
}
