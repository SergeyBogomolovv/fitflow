package admin

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/SergeyBogomolovv/fitflow/pkg/auth"
)

type AdminRepo interface {
	AdminByLogin(ctx context.Context, login string) (*domain.Admin, error)
	AdminExists(ctx context.Context, login string) (bool, error)
	SaveAdmin(ctx context.Context, login string, password []byte) error
	UpdatePassword(ctx context.Context, login string, password []byte) error
	DeleteAdmin(ctx context.Context, login string) error
}

type adminService struct {
	logger *slog.Logger
	repo   AdminRepo
}

func NewAdminService(logger *slog.Logger, repo AdminRepo) *adminService {
	return &adminService{logger, repo}
}

func (s *adminService) CreateAdmin(ctx context.Context, login, password string) error {
	const op = "admin.CreateAdmin"
	logger := s.logger.With(slog.String("op", op), slog.String("login", login))
	logger.Debug("creating admin")

	exists, err := s.repo.AdminExists(ctx, login)
	if err != nil {
		logger.Error("failed to check admin exists", "error", err)
		return err
	}
	if exists {
		logger.Debug("admin already exists")
		return domain.ErrAdminAlreadyExists
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		logger.Error("failed to hash password", "error", err)
		return err
	}
	if err := s.repo.SaveAdmin(ctx, login, hash); err != nil {
		logger.Error("failed to save admin", "error", err)
		return err
	}
	logger.Info("admin created")

	return nil
}

func (s *adminService) UpdatePassword(ctx context.Context, login, oldPass, newPass string) error {
	const op = "admin.UpdatePassword"
	logger := s.logger.With(slog.String("op", op), slog.String("login", login))
	logger.Debug("updating password")

	admin, err := s.repo.AdminByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrAdminNotFound) {
			logger.Debug("admin not exists")
			return domain.ErrInvalidCredentials
		}
		logger.Error("failed to get admin", "error", err)
		return err
	}

	if !auth.ComparePassword(admin.Password, oldPass) {
		logger.Debug("invalid old password")
		return domain.ErrInvalidCredentials
	}

	hash, err := auth.HashPassword(newPass)
	if err != nil {
		logger.Error("failed to hash password", "error", err)
		return err
	}
	if err := s.repo.UpdatePassword(ctx, login, hash); err != nil {
		logger.Error("failed to update password", "error", err)
		return err
	}
	logger.Info("password updated")
	return nil
}

func (s *adminService) RemoveAdmin(ctx context.Context, login string) error {
	const op = "admin.RemoveAdmin"
	logger := s.logger.With(slog.String("op", op), slog.String("login", login))
	logger.Debug("removing admin")

	if err := s.repo.DeleteAdmin(ctx, login); err != nil {
		if errors.Is(err, domain.ErrAdminNotFound) {
			logger.Debug("admin not exists")
			return domain.ErrAdminNotFound
		}
		logger.Error("failed to remove admin", "error", err)
		return err
	}
	logger.Info("admin removed")
	return nil
}
