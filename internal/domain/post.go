package domain

import (
	"errors"
	"mime/multipart"
)

type Post struct {
	ID       int64    `json:"id" example:"123"`
	Content  string   `json:"content" example:"Польза протеина в диете"`
	Audience UserLvl  `json:"audience" example:"beginner"`
	Images   []string `json:"images" example:"image1.jpg,image2.jpg"`
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
