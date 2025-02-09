package user

import (
	"context"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
)

type UserService interface {
	AllUserIds(ctx context.Context) ([]int64, error)
	UserIdsByLvl(ctx context.Context, lvl domain.UserLvl) ([]int64, error)
	AddUser(ctx context.Context, id int64) error
	RemoveUser(ctx context.Context, id int64) error
	UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error
}
