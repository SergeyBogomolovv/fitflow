package user

import (
	"database/sql"
	"errors"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) UsersByLvl(lvl domain.UserLvl) ([]domain.User, error) {
	var entities []user
	query := `SELECT user_id, lvl FROM users WHERE lvl = $1`
	if err := r.db.Select(&entities, query, lvl); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.User{}, nil
		}
		return nil, err
	}
	return mapUsers(entities), nil
}

func (r *userRepo) AllUsers() ([]domain.User, error) {
	var entities []user
	query := `SELECT user_id, lvl FROM users`
	if err := r.db.Select(&entities, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.User{}, nil
		}
		return nil, err
	}
	return mapUsers(entities), nil
}
