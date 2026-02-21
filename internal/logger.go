package internal

import (
	"go.uber.org/zap"
)

func NewLogger(environment string) *zap.Logger {
	var (
		logger *zap.Logger
		err    error
	)

	if environment == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	return logger
}
