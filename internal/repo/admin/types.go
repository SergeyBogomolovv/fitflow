package admin

import "github.com/SergeyBogomolovv/fitflow/internal/domain"

type admin struct {
	Login    string `db:"login"`
	Password []byte `db:"password"`
}

func (a admin) ToDomain() *domain.Admin {
	return &domain.Admin{
		Login:    a.Login,
		Password: a.Password,
	}
}
