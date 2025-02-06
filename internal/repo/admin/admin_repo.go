package admin

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/jmoiron/sqlx"
)

type adminRepo struct {
	db *sqlx.DB
}

func NewAdminRepo(db *sqlx.DB) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) AdminByLogin(ctx context.Context, login string) (*domain.Admin, error) {
	var admin admin
	query := `SELECT login, password FROM admins WHERE login = $1`
	if err := r.db.GetContext(ctx, &admin, query, login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAdminNotFound
		}
		return nil, err
	}
	return admin.ToDomain(), nil
}

func (r *adminRepo) AdminExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	query := `SELECT TRUE FROM admins WHERE login = $1`
	if err := r.db.GetContext(ctx, &exists, query, login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (r *adminRepo) SaveAdmin(ctx context.Context, login string, password []byte) error {
	query := `INSERT INTO admins (login, password) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, login, password)
	return err
}

func (r *adminRepo) UpdatePassword(ctx context.Context, login string, password []byte) error {
	query := `UPDATE admins SET password = $1 WHERE login = $2`
	res, err := r.db.ExecContext(ctx, query, password, login)
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
	query := `DELETE FROM admins WHERE login = $1`
	res, err := r.db.ExecContext(ctx, query, login)
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
