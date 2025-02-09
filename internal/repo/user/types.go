package user

import (
	"github.com/SergeyBogomolovv/fitflow/internal/domain"
)

type user struct {
	ID         int64          `db:"user_id"`
	Lvl        domain.UserLvl `db:"lvl"`
	Subscribed bool           `db:"subscribed"`
}

func (u user) ToDomain() *domain.User {
	return &domain.User{ID: u.ID, Lvl: u.Lvl}
}

func mapUsers(users []user) []domain.User {
	res := make([]domain.User, len(users))
	for i, user := range users {
		res[i] = domain.User{ID: user.ID, Lvl: user.Lvl}
	}
	return res
}
