package content

import "github.com/SergeyBogomolovv/fitflow/internal/domain"

type CreatePostInput struct {
	Content  string
	Audience domain.UserLvl
}
