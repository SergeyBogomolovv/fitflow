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

type PostRepo interface {
	LatestPostByAudience(ctx context.Context, audience domain.UserLvl) (domain.Post, error)
	MarkAsPosted(ctx context.Context, id int64) error
	SavePost(ctx context.Context, in SavePostInput) error
}
