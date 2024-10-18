package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/shevchenko-a-v/auth-service/internal/app"
	"github.com/shevchenko-a-v/auth-service/internal/config"
	"github.com/shevchenko-a-v/auth-service/internal/logger"
	"go.uber.org/zap"
)

func main() {
	config := config.MustLoad()
	logger.MustInitLogger(config)
	defer logger.Logger.Sync()
	logger.Logger.Info("starting application", zap.Any("config", config))

	application := app.New(config.GRPC.Port, config.StoragePath, config.TokenTTL)
	go application.GRPCApp.Run()
	// if err := application.GRPCApp.Run(); err != nil {
	// 	logger.Logger.Error("grpc application failed", zap.Error(err))
	// }
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop
	application.GRPCApp.Stop()
	logger.Logger.Info("application stopped", zap.String("signal", sig.String()))
}
