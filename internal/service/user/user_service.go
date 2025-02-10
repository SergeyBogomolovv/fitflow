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
	SubscribersByLvl(ctx context.Context, lvl domain.UserLvl) ([]domain.User, error)
	Subscribers(ctx context.Context) ([]domain.User, error)
}

type service struct {
	logger *slog.Logger
	repo   UserRepo
}

type UserService interface {
	SaveUser(ctx context.Context, id int64) error
	UpdateSubscribed(ctx context.Context, id int64, subscribed bool) error
	UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error
	SubscribersIdsByLvl(ctx context.Context, lvl domain.UserLvl) ([]int64, error)
	SubscribersIds(ctx context.Context) ([]int64, error)
}

func New(logger *slog.Logger, repo UserRepo) UserService {
	return &service{logger, repo}
}

// SaveUser creates new user if not exists
func (s *service) SaveUser(ctx context.Context, id int64) error {
	const op = "user.SaveUser"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))

	exists, err := s.repo.UserExists(ctx, id)
	if err != nil {
		logger.Error("failed to check user exists", "error", err)
		return err
	}
	if exists {
		return nil
	}

	logger.Debug("saving user")

	if err := s.repo.SaveUser(ctx, id, domain.UserLvlDefault); err != nil {
		logger.Error("failed to save user", "error", err)
		return err
	}

	logger.Info("user saved")
	return nil
}

func (s *service) UpdateSubscribed(ctx context.Context, id int64, subscribed bool) error {
	const op = "user.UpdateSubscribed"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))
	logger.Debug("updating user subscribed")

	if err := s.repo.UpdateSubscribed(ctx, id, subscribed); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Debug("user not exists")
			return domain.ErrUserNotFound
		}
		logger.Error("failed to update user subscribed", "error", err)
		return err
	}

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

func (s *service) SubscribersIdsByLvl(ctx context.Context, lvl domain.UserLvl) ([]int64, error) {
	const op = "user.UsersByLvl"
	logger := s.logger.With(slog.String("op", op), slog.String("lvl", string(lvl)))
	logger.Debug("getting subscribers by level")

	users, err := s.repo.SubscribersByLvl(ctx, lvl)
	if err != nil {
		logger.Error("failed to get subscribers by level", "error", err)
		return nil, err
	}
	res := make([]int64, len(users))
	for i, user := range users {
		res[i] = user.ID
	}
	return res, nil
}

func (s *service) SubscribersIds(ctx context.Context) ([]int64, error) {
	const op = "user.SubscribersIds"
	logger := s.logger.With(slog.String("op", op))
	logger.Debug("getting all subscribers")

	users, err := s.repo.Subscribers(ctx)
	if err != nil {
		logger.Error("failed to get all subscribers", "error", err)
		return nil, err
	}
	res := make([]int64, len(users))
	for i, user := range users {
		res[i] = user.ID
	}
	return res, nil
}
