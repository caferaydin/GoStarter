package repository

import (
	"context"
	"errors"
	"go-starter/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	SaveRefreshToken(ctx context.Context, userID int64, token string, expires time.Time) error
	ValidateRefreshToken(ctx context.Context, token string) (int64, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT id, username, password FROM users WHERE username=$1", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	return err
}

func (r *userRepository) SaveRefreshToken(ctx context.Context, userID int64, token string, expires time.Time) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, expires)
	return err
}

func (r *userRepository) ValidateRefreshToken(ctx context.Context, token string) (int64, error) {
	var userID int64
	var expires time.Time
	err := r.db.QueryRowxContext(ctx,
		"SELECT user_id, expires_at FROM refresh_tokens WHERE token=$1", token).
		Scan(&userID, &expires)
	if err != nil || time.Now().After(expires) {
		return 0, errors.New("invalid or expired refresh token")
	}
	return userID, nil
}
