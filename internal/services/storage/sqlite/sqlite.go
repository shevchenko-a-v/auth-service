package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/shevchenko-a-v/auth-service/internal/domain/models"
	"github.com/shevchenko-a-v/auth-service/internal/services/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (userID int64, err error) {
	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("SaveUser failed: %w", err)
	}

	result, err := stmt.Exec(email, passHash)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) && sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("SaveUser failed: %w", storage.ErrUserExists)
		}
		return 0, fmt.Errorf("SaveUser failed: %w", err)
	}
	uid, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("SaveUser failed: %w", err)
	}
	return uid, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("User failed: %w", err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var retUser models.User
	if err := row.Scan(&retUser.ID, &retUser.Email, &retUser.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("User failed: %w", storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("User failed: %w", err)
	}
	return retUser, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("IsAdmin failed: %w", err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool
	if err := row.Scan(&isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("IsAdmin failed: %w", storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("IsAdmin failed: %w", err)
	}
	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("App failed: %w", err)
	}

	row := stmt.QueryRowContext(ctx, appID)

	var app models.App
	if err := row.Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("App failed: %w", storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("App failed: %w", err)
	}
	return app, nil
}
