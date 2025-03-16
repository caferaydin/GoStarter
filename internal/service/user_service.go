package service

import (
	"context"
	"errors"
	"go-starter/internal/config"
	"go-starter/internal/model"
	"go-starter/internal/repository"
	"go-starter/internal/util"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Authenticate(ctx context.Context, username, password string) (*model.User, error)
	Register(ctx context.Context, user *model.User) error
	RefreshTokens(ctx context.Context, token string) (string, string, error)
	SaveRefreshToken(ctx context.Context, userID int64, token string, expires time.Time) error
}

type userService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (*model.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *userService) Register(ctx context.Context, user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repo.CreateUser(ctx, user)
}

func (s *userService) RefreshTokens(ctx context.Context, token string) (string, string, error) {
	userID, err := s.repo.ValidateRefreshToken(ctx, token)
	if err != nil {
		return "", "", err
	}

	accessToken, err := util.GenerateJWT(userID, s.cfg.JWTSecret, s.cfg.AccessTokenExpiry)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := util.GenerateRefreshToken(userID, s.cfg.RefreshSecret, s.cfg.RefreshTokenExpiry)
	if err != nil {
		return "", "", err
	}

	expires := time.Now().Add(s.cfg.RefreshTokenExpiry)
	_ = s.repo.SaveRefreshToken(ctx, userID, newRefreshToken, expires)

	return accessToken, newRefreshToken, nil
}

func (s *userService) SaveRefreshToken(ctx context.Context, userID int64, token string, expires time.Time) error {
	return s.repo.SaveRefreshToken(ctx, userID, token, expires)
}
