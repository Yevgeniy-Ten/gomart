package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophermart/internal/domain"
	"testing"
)

func TestNewUtils(t *testing.T) {
	cfg := &domain.Config{
		Address:     ":8080",
		DatabaseURL: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		AccrualHost: "http://localhost:8081",
	}

	utils, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, utils)

	assert.NotNil(t, utils.L)
	assert.Equal(t, cfg, utils.C)
}
