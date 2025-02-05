package domain

import "time"

type Post struct {
	ID          int64     `db:"post_id"`
	Content     string    `db:"content"`
	Audience    UserLvl   `db:"audience"`
	Images      []string  `db:"images"`
	ScheduledAt time.Time `db:"scheduled_at"`
}
