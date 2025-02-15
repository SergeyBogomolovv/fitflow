package user

import (
	"context"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
)

type User struct {
	ID         int64          `db:"user_id"`
	Lvl        domain.UserLvl `db:"lvl"`
	Subscribed bool           `db:"subscribed"`
}

func (u User) ToDomain() domain.User {
	return domain.User{ID: u.ID, Lvl: u.Lvl}
}

func mapUsersToDomain(users []User) []domain.User {
	res := make([]domain.User, len(users))
	for i, user := range users {
		res[i] = user.ToDomain()
	}
	return res
}

type UserRepo interface {
	SaveUser(ctx context.Context, id int64, lvl domain.UserLvl) error
	UserExists(ctx context.Context, id int64) (bool, error)
	UpdateSubscribed(ctx context.Context, id int64, subscribed bool) error
	UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error
	Subscribers(ctx context.Context, lvl domain.UserLvl, all bool) ([]domain.User, error)
}
