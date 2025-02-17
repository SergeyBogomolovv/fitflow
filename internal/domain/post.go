package domain

import (
	"errors"
	"mime/multipart"
)

type Post struct {
	ID       int64    `json:"id"`
	Content  string   `json:"content"`
	Audience UserLvl  `json:"audience"`
	Images   []string `json:"images"`
}

var (
	ErrPostNotFound = errors.New("post not found")
	ErrNoPosts      = errors.New("no posts")
)

type CreatePostDTO struct {
	Content  string                  `validate:"required"`
	Audience UserLvl                 `validate:"required,oneof=beginner intermediate advanced default"`
	Images   []*multipart.FileHeader `validate:"required,min=1,dive,required"`
}
