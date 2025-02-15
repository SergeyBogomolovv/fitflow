package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
)

type UserRepo interface {
	SaveUser(ctx context.Context, id int64, lvl domain.UserLvl) error
	UserExists(ctx context.Context, id int64) (bool, error)
	UpdateSubscribed(ctx context.Context, id int64, subscribed bool) error
	UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error
	Subscribers(ctx context.Context, lvl domain.UserLvl, all bool) ([]domain.User, error)
}

type service struct {
	logger   *slog.Logger
	userRepo UserRepo
}

func New(logger *slog.Logger, userRepo UserRepo) *service {
	return &service{logger, userRepo}
}

func (s *service) EnsureUserExists(ctx context.Context, id int64) error {
	const op = "user.EnsureUserExists"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))

	exists, err := s.userRepo.UserExists(ctx, id)
	if err != nil {
		logger.Error("failed to check user exists", "error", err)
		return err
	}
	if exists {
		return nil
	}

	if err := s.userRepo.SaveUser(ctx, id, domain.UserLvlDefault); err != nil {
		logger.Error("failed to save user", "error", err)
		return err
	}

	return nil
}

func (s *service) UpdateSubscribed(ctx context.Context, id int64, subscribed bool) error {
	const op = "user.UpdateSubscribed"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))

	if err := s.userRepo.UpdateSubscribed(ctx, id, subscribed); err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			logger.Error("failed to update user subscribed", "error", err)
		}
		return err
	}

	return nil
}

func (s *service) UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error {
	const op = "user.UpdateUserLvl"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))

	if err := s.userRepo.UpdateUserLvl(ctx, id, lvl); err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			logger.Error("failed to update user level", "error", err)
		}
		return err
	}

	return nil
}

// if lvl is default it returns all subscribers ids
func (s *service) SubscribersIds(ctx context.Context, lvl domain.UserLvl) ([]int64, error) {
	const op = "user.SubscribersIds"
	logger := s.logger.With(slog.String("op", op))

	var all bool
	if lvl == domain.UserLvlDefault {
		all = true
	}

	users, err := s.userRepo.Subscribers(ctx, lvl, all)
	if err != nil {
		logger.Error("failed to get subscribers", "error", err)
		return nil, err
	}
	res := make([]int64, len(users))
	for i, user := range users {
		res[i] = user.ID
	}
	return res, nil
}
