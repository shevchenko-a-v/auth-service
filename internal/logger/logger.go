package logger

import (
	"github.com/shevchenko-a-v/auth-service/internal/config"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func MustInitLogger(cfg *config.Config) {
	if cfg.IsProd() {
		Logger = zap.Must(zap.NewProduction())
	} else {
		Logger = zap.Must(zap.NewDevelopment())
	}
}
