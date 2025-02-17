package post

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
)

type PostRepo interface {
	LatestByAudience(ctx context.Context, audience domain.UserLvl) (domain.Post, error)
	MarkAsPosted(ctx context.Context, id int64) error
}

type postService struct {
	logger   *slog.Logger
	postRepo PostRepo
}

func New(logger *slog.Logger, repo PostRepo) *postService {
	return &postService{logger, repo}
}

func (s *postService) PickLatest(ctx context.Context, audience domain.UserLvl) (domain.Post, error) {
	const op = "post.PickLatest"
	logger := s.logger.With(slog.String("op", op), slog.String("audience", string(audience)))

	post, err := s.postRepo.LatestByAudience(ctx, audience)
	if err != nil {
		if errors.Is(err, domain.ErrNoPosts) {
			return domain.Post{}, domain.ErrNoPosts
		}
		logger.Error("failed to get latest post", "error", err)
		return domain.Post{}, err
	}
	return post, nil
}

func (s *postService) MarkAsPosted(ctx context.Context, id int64) error {
	const op = "post.MarkAsPosted"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))

	if err := s.postRepo.MarkAsPosted(ctx, id); err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return domain.ErrPostNotFound
		}
		logger.Error("failed to mark post as posted", "error", err)
		return err
	}
	return nil
}
