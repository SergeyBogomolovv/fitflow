package dto

import (
	"time"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
)

type CreatePost struct {
	Content     string         `json:"content" validate:"required"`
	Audience    domain.UserLvl `json:"audience" validate:"required,oneof=default beginner intermediate advanced"`
	Images      []string       `json:"images"`
	ScheduledAt time.Time      `json:"scheduled_at" validate:"required,future"`
}
