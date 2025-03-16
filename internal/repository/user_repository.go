package repository

import (
	"context"
	"errors"
	"go-starter/internal/model"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

type userRepository struct {
	users map[string]*model.User
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users: map[string]*model.User{
			"admin": {ID: 1, Username: "admin", Password: "1234"},
		},
	}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, ok := r.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
