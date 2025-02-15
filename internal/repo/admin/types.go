package admin

import (
	"context"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
)

type Admin struct {
	Login    string `db:"login"`
	Password []byte `db:"password"`
}

func (a Admin) ToDomain() domain.Admin {
	return domain.Admin{
		Login:    a.Login,
		Password: a.Password,
	}
}

type AdminRepo interface {
	AdminByLogin(ctx context.Context, login string) (domain.Admin, error)
	AdminExists(ctx context.Context, login string) (bool, error)
	SaveAdmin(ctx context.Context, login string, password []byte) error
	UpdatePassword(ctx context.Context, login string, password []byte) error
	DeleteAdmin(ctx context.Context, login string) error
}
