package post

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type postRepo struct {
	qb sq.StatementBuilderType
	db *sqlx.DB
}

func New(db *sqlx.DB) PostRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &postRepo{db: db, qb: qb}
}

func (r *postRepo) LatestPostByAudience(ctx context.Context, audience domain.UserLvl) (domain.Post, error) {
	query, args := r.qb.
		Select("post_id", "content", "audience", "images", "created_at", "posted").
		From("posts").
		Where(sq.Eq{"audience": audience, "posted": false}).
		OrderBy("created_at DESC").
		Limit(1).
		MustSql()

	var post Post
	if err := r.db.GetContext(ctx, &post, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Post{}, domain.ErrNoPosts
		}
		return domain.Post{}, fmt.Errorf("failed to get latest post: %w", err)
	}
	return post.ToDomain(), nil
}

func (r *postRepo) MarkAsPosted(ctx context.Context, id int64) error {
	query, args := r.qb.Update("posts").Set("posted", true).Where(sq.Eq{"post_id": id}).MustSql()
	return r.execOrNotFound(ctx, query, args)
}

func (r *postRepo) SavePost(ctx context.Context, in SavePostInput) (domain.Post, error) {
	query, args := r.qb.
		Insert("posts").
		Columns("content", "audience", "images").
		Values(in.Content, in.Audience, pq.Array(in.Images)).
		Suffix("RETURNING post_id, content, audience, images").
		MustSql()

	post := Post{}
	if err := r.db.GetContext(ctx, &post, query, args...); err != nil {
		return domain.Post{}, fmt.Errorf("failed to save post: %w", err)
	}
	return post.ToDomain(), nil
}

func (r *postRepo) RemovePost(ctx context.Context, id int64) (domain.Post, error) {
	query, args := r.qb.
		Delete("posts").
		Where(sq.Eq{"post_id": id}).
		Suffix("RETURNING post_id, content, audience, images").
		MustSql()

	var post Post
	if err := r.db.GetContext(ctx, &post, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Post{}, domain.ErrPostNotFound
		}
		return domain.Post{}, err
	}

	return post.ToDomain(), nil
}

func (r *postRepo) execOrNotFound(ctx context.Context, query string, args []any) error {
	res, err := r.db.ExecContext(ctx, query, args...)
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
