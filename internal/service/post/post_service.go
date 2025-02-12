package post

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	postRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/post"
)

type PostRepo interface {
	SavePost(ctx context.Context, post postRepo.CreatePostInput) error
	MarkAsPosted(ctx context.Context, id int64) error
	LatestPostByAudience(ctx context.Context, audience domain.UserLvl) (*domain.Post, error)
}

type postService struct {
	logger   *slog.Logger
	postRepo PostRepo
}

func New(logger *slog.Logger, repo PostRepo) *postService {
	return &postService{logger, repo}
}

func (s *postService) PickLatest(ctx context.Context, audience domain.UserLvl) (*domain.Post, error) {
	const op = "post.PickLatest"
	logger := s.logger.With(slog.String("op", op), slog.String("audience", string(audience)))

	post, err := s.postRepo.LatestPostByAudience(ctx, audience)
	if err != nil {
		if errors.Is(err, domain.ErrNoPosts) {
			logger.Debug("no posts")
			return nil, domain.ErrNoPosts
		}
		logger.Error("failed to get latest post", "error", err)
		return nil, err
	}
	return post, nil
}

func (s *postService) MarkAsPosted(ctx context.Context, id int64) error {
	const op = "post.MarkAsPosted"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))
	logger.Debug("marking post as posted")

	if err := s.postRepo.MarkAsPosted(ctx, id); err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			logger.Debug("post not exists")
			return domain.ErrPostNotFound
		}
		logger.Error("failed to mark post as posted", "error", err)
		return err
	}
	return nil
}
