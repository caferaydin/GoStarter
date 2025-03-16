package service

import (
	"context"
	"errors"
	"go-starter/internal/model"
	"go-starter/internal/repository"
)

type UserService interface {
	Authenticate(ctx context.Context, username, password string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (*model.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil || user.Password != password {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
