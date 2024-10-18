package auth

import (
	"context"

	ssov1 "github.com/shevchenko-a-v/protofiles/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyAppId  int32 = 0
	emptyUserId int64 = 0
)

type AuthInterface interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth AuthInterface
}

func Register(gRPC *grpc.Server, auth AuthInterface) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, request *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := validateLogin(request); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, request.GetEmail(), request.GetPassword(), int(request.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, request *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(request); err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, request *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(request); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, request.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateLogin(request *ssov1.LoginRequest) error {
	if request.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "empty email")
	}
	if request.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "empty password")
	}
	if request.GetAppId() == emptyAppId {
		return status.Error(codes.InvalidArgument, "empty app_id")
	}
	return nil
}

func validateRegister(request *ssov1.RegisterRequest) error {
	if request.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "empty email")
	}
	if request.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "empty password")
	}
	return nil
}

func validateIsAdmin(request *ssov1.IsAdminRequest) error {
	if request.GetUserId() == emptyUserId {
		return status.Error(codes.InvalidArgument, "empty user_id")
	}
	return nil
}
