package app

import (
	"time"

	grpcapp "github.com/shevchenko-a-v/auth-service/internal/app/grpc"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	grpcApp := grpcapp.New(grpcPort)
	return &App{GRPCApp: grpcApp}
}
