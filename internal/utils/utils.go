package utils

import (
	"gophermart/internal/domain"
	"gophermart/internal/utils/session"
)

func New(c *domain.Config) (*domain.Utils, error) {
	l, err := NewLogger()
	if err != nil {
		return nil, err
	}
	return &domain.Utils{L: l, C: c, S: session.NewSession()}, nil
}
