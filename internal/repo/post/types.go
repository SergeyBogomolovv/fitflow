package post

import (
	"time"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/lib/pq"
)

type CreatePostInput struct {
	Content  string
	Audience domain.UserLvl
	Images   []string
}

type post struct {
	ID        int64          `db:"post_id"`
	Content   string         `db:"content"`
	Audience  domain.UserLvl `db:"audience"`
	Images    pq.StringArray `db:"images"`
	CreatedAt time.Time      `db:"created_at"`
	Posted    bool           `db:"posted"`
}

func (p post) ToDomain() *domain.Post {
	return &domain.Post{
		ID:       p.ID,
		Content:  p.Content,
		Audience: p.Audience,
		Images:   p.Images,
	}
}

func mapPosts(posts []post) []domain.Post {
	res := make([]domain.Post, 0, len(posts))
	for _, p := range posts {
		res = append(res, *p.ToDomain())
	}
	return res
}
