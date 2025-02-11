package post

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type postRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) LatestPostByAudience(ctx context.Context, audience domain.UserLvl) (*domain.Post, error) {
	post := new(post)
	query := `
	SELECT post_id, content, audience, images, created_at, posted
	FROM posts 
	WHERE posted = false AND audience = $1
	ORDER BY created_at DESC
	LIMIT 1`
	if err := r.db.GetContext(ctx, post, query, audience); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoPosts
		}
		return nil, err
	}
	return post.ToDomain(), nil
}

func (r *postRepo) MarkAsPosted(ctx context.Context, id int64) error {
	query := `UPDATE posts SET posted = true WHERE post_id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return domain.ErrPostNotFound
	}
	return nil
}

func (r *postRepo) SavePost(ctx context.Context, post CreatePostInput) error {
	query := `INSERT INTO posts (content, audience, images) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, post.Content, post.Audience, pq.Array(post.Images))
	return err
}
