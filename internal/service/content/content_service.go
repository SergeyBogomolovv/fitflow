package content

import (
	"context"
	"io"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	postRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/post"
)

type PostRepo interface {
	SavePost(ctx context.Context, in postRepo.SavePostInput) error
}

type S3Client interface {
	Upload(ctx context.Context, key string, body io.Reader) (string, error)
	Delete(ctx context.Context, key string) error
}

type AiGenerator interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}

type postService struct {
	logger   *slog.Logger
	postRepo PostRepo
	ai       AiGenerator
	s3       S3Client
}

func New(logger *slog.Logger, repo PostRepo, ai AiGenerator, s3 S3Client) *postService {
	return &postService{logger, repo, ai, s3}
}

func (s *postService) GenerateContent(ctx context.Context, theme string) (string, error) {
	return s.ai.GenerateContent(ctx, theme)
}

func (s *postService) CreatePost(ctx context.Context, in CreatePostInput) (domain.Post, error) {
	return domain.Post{}, nil
}
