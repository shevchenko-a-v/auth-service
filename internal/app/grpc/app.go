package grpcapp

import (
	"fmt"
	"net"

	authgrpc "github.com/shevchenko-a-v/auth-service/internal/grpc/auth"
	"github.com/shevchenko-a-v/auth-service/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
}

func New(auth authgrpc.AuthInterface, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, auth)
	return &App{gRPCServer: gRPCServer, port: port}
}

func (a *App) Run() error {
	logger.Logger.Info("starting grpc server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("couldn't start grpc server: %w", err)
	}

	logger.Logger.Info("grpc server is running", zap.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("grpc server error: %w", err)
	}
	return nil
}

func (a *App) Stop() {
	logger.Logger.Info("stopping grpc server", zap.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
