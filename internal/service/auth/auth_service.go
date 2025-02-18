package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/SergeyBogomolovv/fitflow/pkg/auth"
)

type AdminRepo interface {
	AdminByLogin(ctx context.Context, login string) (domain.Admin, error)
}

type service struct {
	logger    *slog.Logger
	adminRepo AdminRepo
	jwtSecret []byte
	jwtTTL    time.Duration
}

func New(logger *slog.Logger, repo AdminRepo, jwtSecret []byte, jwtTTL time.Duration) *service {
	return &service{logger, repo, jwtSecret, jwtTTL}
}

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	const op = "auth.Login"
	logger := s.logger.With(slog.String("op", op), slog.String("login", login))

	admin, err := s.adminRepo.AdminByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrAdminNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		logger.Error("failed to get admin", "error", err)
		return "", err
	}

	if !auth.ComparePassword(admin.Password, password) {
		return "", domain.ErrInvalidCredentials
	}

	token, err := auth.SignJWT(admin.Login, s.jwtSecret, s.jwtTTL)
	if err != nil {
		logger.Error("failed to sign JWT", "error", err)
		return "", err
	}

	return token, nil
}

func (s *service) AuthFunc(tokenString string) (context.Context, error) {
	aud, err := auth.VerifyJWT(tokenString, s.jwtSecret)
	if err != nil {
		return nil, err
	}
	return context.WithValue(context.Background(), auth.AdminLoginKey{}, aud), nil
}
