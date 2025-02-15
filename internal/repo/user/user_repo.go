package user

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	qb sq.StatementBuilderType
	db *sqlx.DB
}

func New(db *sqlx.DB) UserRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &userRepo{db: db, qb: qb}
}

func (r *userRepo) SaveUser(ctx context.Context, id int64, lvl domain.UserLvl) error {
	query, args := r.qb.Insert("users").Columns("user_id", "lvl").Values(id, lvl).MustSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *userRepo) UserExists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	query, args := r.qb.Select("TRUE").From("users").Where(sq.Eq{"user_id": id}).MustSql()
	if err := r.db.GetContext(ctx, &exists, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (r *userRepo) UpdateSubscribed(ctx context.Context, id int64, subscribed bool) error {
	query, args := r.qb.Update("users").Set("subscribed", subscribed).Where(sq.Eq{"user_id": id}).MustSql()
	res, err := r.db.ExecContext(ctx, query, args...)
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

func (r *userRepo) UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error {
	query, args := r.qb.Update("users").Set("lvl", lvl).Where(sq.Eq{"user_id": id}).MustSql()
	res, err := r.db.ExecContext(ctx, query, args...)
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

func (r *userRepo) Subscribers(ctx context.Context, lvl domain.UserLvl, all bool) ([]domain.User, error) {
	var entities []User
	q := r.qb.Select("user_id", "lvl").From("users").Where(sq.Eq{"subscribed": true})
	if !all {
		q = q.Where(sq.Eq{"lvl": lvl})
	}
	query, args := q.MustSql()

	if err := r.db.SelectContext(ctx, &entities, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.User{}, nil
		}
		return nil, err
	}
	return mapUsersToDomain(entities), nil
}
