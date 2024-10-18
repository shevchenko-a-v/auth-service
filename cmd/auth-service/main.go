package main

import (
	"github.com/shevchenko-a-v/auth-service/internal/config"
	"github.com/shevchenko-a-v/auth-service/internal/logger"
	"go.uber.org/zap"
)

func main() {
	config := config.MustLoad()
	logger.MustInitLogger(config)
	defer logger.Logger.Sync()
	logger.Logger.Info("starting application", zap.Any("config", config))
}
