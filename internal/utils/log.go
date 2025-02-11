package utils

import "go.uber.org/zap"

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	myLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	defer myLogger.Sync()
	return myLogger, nil
}
