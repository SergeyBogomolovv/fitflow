package domain

import (
	"errors"
)

type Post struct {
	ID       int64
	Content  string
	Audience UserLvl
	Images   []string
}

var (
	ErrPostNotFound = errors.New("post not found")
	ErrNoPosts      = errors.New("no posts")
)
