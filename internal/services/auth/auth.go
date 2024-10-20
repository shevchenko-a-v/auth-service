package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shevchenko-a-v/auth-service/internal/domain/models"
	"github.com/shevchenko-a-v/auth-service/internal/jwt"
	"github.com/shevchenko-a-v/auth-service/internal/logger"
	"github.com/shevchenko-a-v/auth-service/internal/services/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (userID int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidAppID       = errors.New("invalid app ID")
)

func New(userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{userSaver: userSaver, userProvider: userProvider, appProvider: appProvider, tokenTTL: tokenTTL}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	log := logger.Logger.With(
		zap.String("op", "Login"),
	)
	log.Info("trying to login user")
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", zap.Error(err))
			return "", fmt.Errorf("Login failed: %w", ErrInvalidCredentials)
		}
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", zap.Error(err))
			return "", fmt.Errorf("Login failed: %w", ErrInvalidAppID)
		}

		log.Error("failed to get user", zap.Error(err))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Warn("invalid credentials", zap.Error(err))

		return "", fmt.Errorf("Login failed: %w", ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		log.Warn("application not found", zap.Int("id", appID))
		return "", fmt.Errorf("Application not found: %w", err)
	}

	log.Info("login successful")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", zap.Error(err))
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email string, password string) (int64, error) {
	log := logger.Logger.With(
		zap.String("op", "Register"),
		// gzap.String("email", email),
	)
	log.Info("registering user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", zap.Error(err))
			return 0, fmt.Errorf("Register failed: %w", ErrInvalidCredentials)
		}
		log.Error("failed to generate password hash", zap.Error(err))
		return 0, err
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", zap.Error(err))
		return 0, err
	}

	log.Info("user registered")
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	log := logger.Logger.With(
		zap.String("op", "IsAdmin"),
	)

	log.Info("checking if user is an admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", zap.Error(err))
			return false, fmt.Errorf("IsAdmin failed: %w", ErrInvalidUserID)
		}
		log.Error("failed to check if user is an admin", zap.Error(err))
		return false, fmt.Errorf("failed to check if user is an admin: %w", err)
	}

	log.Info("checked if user is an admin", zap.Bool("IsAdmin", isAdmin))

	return isAdmin, nil
}
