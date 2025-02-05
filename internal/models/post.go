package models

import "time"

type Post struct {
	ID          int64     `db:"post_id"`
	Content     string    `db:"content"`
	Audience    UserLvl   `db:"audience"`
	Images      []string  `db:"images"`
	CreatedAt   time.Time `db:"created_at"`
	ScheduledAt time.Time `db:"scheduled_at"`
	Posted      bool      `db:"posted"`
}

type CreatePostDTO struct {
	Content     string    `json:"content" validate:"required"`
	Audience    UserLvl   `json:"audience" validate:"required,oneof=default beginner intermediate advanced"`
	Images      []string  `json:"images"`
	ScheduledAt time.Time `json:"scheduled_at" validate:"required,future"`
}
