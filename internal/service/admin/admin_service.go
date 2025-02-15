package admin

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/SergeyBogomolovv/fitflow/pkg/auth"
)

type AdminRepo interface {
	AdminByLogin(ctx context.Context, login string) (domain.Admin, error)
	AdminExists(ctx context.Context, login string) (bool, error)
	SaveAdmin(ctx context.Context, login string, password []byte) error
	UpdatePassword(ctx context.Context, login string, password []byte) error
	DeleteAdmin(ctx context.Context, login string) error
}

type adminService struct {
	logger    *slog.Logger
	adminRepo AdminRepo
}

func New(logger *slog.Logger, adminRepo AdminRepo) *adminService {
	return &adminService{logger, adminRepo}
}

func (s *adminService) CreateAdmin(ctx context.Context, login, password string) error {
	const op = "admin.CreateAdmin"
	logger := s.logger.With(slog.String("op", op), slog.String("login", login))

	exists, err := s.adminRepo.AdminExists(ctx, login)
	if err != nil {
		logger.Error("failed to check admin exists", "error", err)
		return err
	}
	if exists {
		return domain.ErrAdminAlreadyExists
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		logger.Error("failed to hash password", "error", err)
		return err
	}
	if err := s.adminRepo.SaveAdmin(ctx, login, hash); err != nil {
		logger.Error("failed to save admin", "error", err)
		return err
	}

	logger.Info("admin created")
	return nil
}

func (s *adminService) UpdatePassword(ctx context.Context, login, oldPass, newPass string) error {
	const op = "admin.UpdatePassword"
	logger := s.logger.With(slog.String("op", op), slog.String("login", login))

	admin, err := s.adminRepo.AdminByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrAdminNotFound) {
			return domain.ErrInvalidCredentials
		}
		logger.Error("failed to get admin", "error", err)
		return err
	}

	if !auth.ComparePassword(admin.Password, oldPass) {
		return domain.ErrInvalidCredentials
	}

	hash, err := auth.HashPassword(newPass)
	if err != nil {
		logger.Error("failed to hash password", "error", err)
		return err
	}
	if err := s.adminRepo.UpdatePassword(ctx, login, hash); err != nil {
		logger.Error("failed to update password", "error", err)
		return err
	}

	logger.Info("admin password updated")
	return nil
}

func (s *adminService) RemoveAdmin(ctx context.Context, login string) error {
	const op = "admin.RemoveAdmin"
	logger := s.logger.With(slog.String("op", op), slog.String("login", login))

	if err := s.adminRepo.DeleteAdmin(ctx, login); err != nil {
		if errors.Is(err, domain.ErrAdminNotFound) {
			return domain.ErrAdminNotFound
		}
		logger.Error("failed to remove admin", "error", err)
		return err
	}

	logger.Info("admin removed")
	return nil
}
