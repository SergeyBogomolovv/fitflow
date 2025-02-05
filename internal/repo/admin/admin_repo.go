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

var ErrAdminNotFound = errors.New("admin not found")

func (r *adminRepo) AdminByLogin(ctx context.Context, login string) (*domain.Admin, error) {
	var admin admin
	query := `SELECT login, password FROM admins WHERE login = $1`
	if err := r.db.GetContext(ctx, &admin, query, login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAdminNotFound
		}
		return nil, err
	}
	return admin.ToDomain(), nil
}

func (r *adminRepo) SaveAdmin(ctx context.Context, login string, password []byte) error {
	query := `INSERT INTO admins (login, password) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, login, password)
	return err
}

func (r *adminRepo) UpdatePassword(ctx context.Context, login string, password []byte) error {
	query := `UPDATE admins SET password = $1 WHERE login = $2`
	_, err := r.db.ExecContext(ctx, query, password, login)
	return err
}

func (r *adminRepo) DeleteAdmin(ctx context.Context, login string) error {
	query := `DELETE FROM admins WHERE login = $1`
	_, err := r.db.ExecContext(ctx, query, login)
	return err
}
