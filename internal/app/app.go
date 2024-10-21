package app

import (
	"time"

	grpcapp "github.com/shevchenko-a-v/auth-service/internal/app/grpc"
	"github.com/shevchenko-a-v/auth-service/internal/services/auth"
	"github.com/shevchenko-a-v/auth-service/internal/services/storage/sqlite"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(authService, grpcPort)
	return &App{GRPCApp: grpcApp}
}
