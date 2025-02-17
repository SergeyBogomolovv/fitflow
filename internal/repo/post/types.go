package post

import (
	"context"
	"time"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/lib/pq"
)

type SavePostInput struct {
	Content  string
	Audience domain.UserLvl
	Images   []string
}

type Post struct {
	ID        int64          `db:"post_id"`
	Content   string         `db:"content"`
	Audience  domain.UserLvl `db:"audience"`
	Images    pq.StringArray `db:"images"`
	CreatedAt time.Time      `db:"created_at"`
	Posted    bool           `db:"posted"`
}

func (p Post) ToDomain() domain.Post {
	return domain.Post{
		ID:       p.ID,
		Content:  p.Content,
		Audience: p.Audience,
		Images:   p.Images,
	}
}

func mapPostsToDomain(posts []Post) []domain.Post {
	res := make([]domain.Post, 0, len(posts))
	for _, post := range posts {
		res = append(res, post.ToDomain())
	}
	return res
}

type PostRepo interface {
	LatestByAudience(ctx context.Context, audience domain.UserLvl) (domain.Post, error)
	MarkAsPosted(ctx context.Context, id int64) error
	Save(ctx context.Context, in SavePostInput) (domain.Post, error)
	Remove(ctx context.Context, id int64) (domain.Post, error)
	List(ctx context.Context, audience domain.UserLvl, incoming bool) ([]domain.Post, error)
}
