package content

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	postRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/post"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type PostRepo interface {
	Save(ctx context.Context, in postRepo.SavePostInput) (domain.Post, error)
	Remove(ctx context.Context, id int64) (domain.Post, error)
	List(ctx context.Context, audience domain.UserLvl, incoming bool) ([]domain.Post, error)
}

type S3Client interface {
	Upload(ctx context.Context, key string, body io.Reader) (string, error)
	Delete(ctx context.Context, url string) error
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

const ImagesFolder = "images"

func New(logger *slog.Logger, repo PostRepo, ai AiGenerator, s3 S3Client) *postService {
	return &postService{logger, repo, ai, s3}
}

func (s *postService) GenerateContent(ctx context.Context, theme string) (string, error) {
	return s.ai.GenerateContent(ctx, theme)
}

func (s *postService) CreatePost(ctx context.Context, in domain.CreatePostDTO) (domain.Post, error) {
	const op = "content.CreatePost"
	logger := s.logger.With(slog.String("op", op))
	input := postRepo.SavePostInput{
		Content:  in.Content,
		Images:   make([]string, 0, len(in.Images)),
		Audience: in.Audience,
	}

	eg, uploadCtx := errgroup.WithContext(ctx)
	for _, imageHeader := range in.Images {
		eg.Go(func() error {
			image, err := imageHeader.Open()
			if err != nil {
				logger.Error("failed to open image", "error", err)
				return err
			}
			defer image.Close()
			key, err := s.s3.Upload(uploadCtx, fmt.Sprintf("%s/%s.jpg", ImagesFolder, uuid.NewString()), image)
			if err != nil {
				logger.Error("failed to upload image", "error", err)
				return err
			}
			input.Images = append(input.Images, key)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return domain.Post{}, err
	}

	post, err := s.postRepo.Save(ctx, input)
	if err != nil {
		logger.Error("failed to save post", "error", err)
		return domain.Post{}, err
	}
	return post, nil
}

func (s *postService) RemovePost(ctx context.Context, id int64) error {
	const op = "content.RemovePost"
	logger := s.logger.With(slog.String("op", op), slog.Int64("id", id))

	post, err := s.postRepo.Remove(ctx, id)
	if err != nil {
		if !errors.Is(err, domain.ErrPostNotFound) {
			logger.Error("failed to remove post", "err", err)
		}
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, url := range post.Images {
		eg.Go(func() error {
			err := s.s3.Delete(ctx, url)
			if err != nil {
				logger.Error("failed to remove image", "err", err)
			}
			return err
		})
	}

	return eg.Wait()
}

func (s *postService) Posts(ctx context.Context, audience domain.UserLvl, incoming bool) ([]domain.Post, error) {
	return s.postRepo.List(ctx, audience, incoming)
}
