package utils

import "go.uber.org/zap"

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	myLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck // ignore error because it's not important
	defer myLogger.Sync()
	return myLogger, nil
}
