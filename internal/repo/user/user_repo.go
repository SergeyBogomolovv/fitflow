package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) SaveUser(ctx context.Context, id int64, lvl domain.UserLvl) (*domain.User, error) {
	query := `INSERT INTO users (user_id, lvl) VALUES ($1, $2) RETURNING user_id, lvl`
	user := new(user)
	if err := r.db.GetContext(ctx, &user, query, id, lvl); err != nil {
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *userRepo) UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error {
	query := `UPDATE users SET lvl = $1 WHERE user_id = $2 `
	res, err := r.db.ExecContext(ctx, query, lvl, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return domain.ErrUserNotFound
	}
	return err
}

func (r *userRepo) UserByID(ctx context.Context, id int64) (*domain.User, error) {
	var user user
	query := `SELECT user_id, lvl FROM users WHERE user_id = $1`
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE user_id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepo) UsersByLvl(ctx context.Context, lvl domain.UserLvl) ([]domain.User, error) {
	var entities []user
	query := `SELECT user_id, lvl FROM users WHERE lvl = $1`
	if err := r.db.SelectContext(ctx, &entities, query, lvl); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.User{}, nil
		}
		return nil, err
	}
	return mapUsers(entities), nil
}

func (r *userRepo) AllUsers(ctx context.Context) ([]domain.User, error) {
	var entities []user
	query := `SELECT user_id, lvl FROM users`
	if err := r.db.SelectContext(ctx, &entities, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.User{}, nil
		}
		return nil, err
	}
	return mapUsers(entities), nil
}
