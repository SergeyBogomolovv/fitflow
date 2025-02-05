package post

import (
	"time"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
)

type post struct {
	ID          int64          `db:"post_id"`
	Content     string         `db:"content"`
	Audience    domain.UserLvl `db:"audience"`
	Images      []string       `db:"images"`
	CreatedAt   time.Time      `db:"created_at"`
	ScheduledAt time.Time      `db:"scheduled_at"`
	Posted      bool           `db:"posted"`
}

func (p post) ToDomain() domain.Post {
	return domain.Post{
		ID:          p.ID,
		Content:     p.Content,
		Audience:    p.Audience,
		Images:      p.Images,
		ScheduledAt: p.ScheduledAt,
	}
}

func mapPosts(entities []post) []domain.Post {
	posts := make([]domain.Post, len(entities))
	for i, entity := range entities {
		posts[i] = entity.ToDomain()
	}
	return posts
}
