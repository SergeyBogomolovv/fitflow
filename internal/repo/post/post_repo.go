package post

import (
	"database/sql"
	"errors"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) PostsByAudience(audience domain.UserLvl) ([]domain.Post, error) {
	var entities []post
	query := `SELECT post_id, content, audience, images, created_at, scheduled_at, posted FROM posts WHERE audience = $1`
	if err := r.db.Select(&entities, query, audience); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Post{}, nil
		}
		return nil, err
	}
	return mapPosts(entities), nil
}
