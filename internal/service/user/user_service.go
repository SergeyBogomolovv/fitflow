package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
)

type UserRepo interface {
	SaveUser(ctx context.Context, id int64, lvl domain.UserLvl) (*domain.User, error)
	UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error
	DeleteUser(ctx context.Context, id int64) error
	UsersByLvl(ctx context.Context, lvl domain.UserLvl) ([]domain.User, error)
	AllUsers(ctx context.Context) ([]domain.User, error)
}

type service struct {
	logger *slog.Logger
	repo   UserRepo
}

func New(logger *slog.Logger, repo UserRepo) UserService {
	return &service{logger, repo}
}

func (s *service) AddUser(ctx context.Context, id int64) error {
	const op = "user.AddUser"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))
	logger.Debug("adding user")

	user, err := s.repo.SaveUser(ctx, id, domain.UserLvlDefault)
	if err != nil {
		logger.Error("failed to save user", "error", err)
		return err
	}

	logger.Info("user added", "user", user)
	return nil
}

func (s *service) RemoveUser(ctx context.Context, id int64) error {
	const op = "user.RemoveUser"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))
	logger.Debug("deleting user")

	if err := s.repo.DeleteUser(ctx, id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Debug("user not exists")
			return domain.ErrUserNotFound
		}
		logger.Error("failed to delete user", "error", err)
		return err
	}

	logger.Info("user deleted")
	return nil
}

func (s *service) UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error {
	const op = "user.UpdateUserLvl"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))
	logger.Debug("updating user level")

	if err := s.repo.UpdateUserLvl(ctx, id, lvl); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Debug("user not exists")
			return domain.ErrUserNotFound
		}
		logger.Error("failed to update user level", "error", err)
		return err
	}

	logger.Info("user level updated", "level", lvl)
	return nil
}

func (s *service) UserIdsByLvl(ctx context.Context, lvl domain.UserLvl) ([]int64, error) {
	const op = "user.UsersByLvl"
	logger := s.logger.With(slog.String("op", op), slog.String("lvl", string(lvl)))
	logger.Debug("getting users by level")

	users, err := s.repo.UsersByLvl(ctx, lvl)
	if err != nil {
		logger.Error("failed to get users by level", "error", err)
		return nil, err
	}
	res := make([]int64, len(users))
	for i, user := range users {
		res[i] = user.ID
	}
	return res, nil
}

func (s *service) AllUserIds(ctx context.Context) ([]int64, error) {
	const op = "user.AllUserIds"
	logger := s.logger.With(slog.String("op", op))
	logger.Debug("getting all users")

	users, err := s.repo.AllUsers(ctx)
	if err != nil {
		logger.Error("failed to get all users", "error", err)
		return nil, err
	}
	res := make([]int64, len(users))
	for i, user := range users {
		res[i] = user.ID
	}
	return res, nil
}
