package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/jmoiron/sqlx"
)

type adminRepo struct {
	qb sq.StatementBuilderType
	db *sqlx.DB
}

func New(db *sqlx.DB) AdminRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &adminRepo{db: db, qb: qb}
}

func (r *adminRepo) AdminByLogin(ctx context.Context, login string) (domain.Admin, error) {
	var admin Admin
	query, args := r.qb.Select("login", "password").From("admins").Where(sq.Eq{"login": login}).MustSql()
	if err := r.db.GetContext(ctx, &admin, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Admin{}, domain.ErrAdminNotFound
		}
		return domain.Admin{}, fmt.Errorf("failed to get admin: %w", err)
	}
	return admin.ToDomain(), nil
}

func (r *adminRepo) AdminExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	query, args := r.qb.Select("TRUE").From("admins").Where(sq.Eq{"login": login}).MustSql()
	if err := r.db.GetContext(ctx, &exists, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check admin: %w", err)
	}
	return exists, nil
}

func (r *adminRepo) SaveAdmin(ctx context.Context, login string, password []byte) error {
	query, args := r.qb.Insert("admins").Columns("login", "password").Values(login, password).MustSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *adminRepo) UpdatePassword(ctx context.Context, login string, password []byte) error {
	query, args := r.qb.Update("admins").Set("password", password).Where(sq.Eq{"login": login}).MustSql()
	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return domain.ErrAdminNotFound
	}
	return nil
}

func (r *adminRepo) DeleteAdmin(ctx context.Context, login string) error {
	query, args := r.qb.Delete("admins").Where(sq.Eq{"login": login}).MustSql()
	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return domain.ErrAdminNotFound
	}
	return nil
}
