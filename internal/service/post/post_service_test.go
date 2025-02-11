package post_test

import (
	"context"
	"testing"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	postRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/post"
	"github.com/SergeyBogomolovv/fitflow/internal/service/post"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostService_PickLatest(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockPostRepo)
	svc := post.New(testutils.NewTestLogger(), mockRepo)

	t.Run("success", func(t *testing.T) {
		expected := &domain.Post{Content: "content"}
		mockRepo.On("LatestPostByAudience", ctx, domain.UserLvlDefault).Return(expected, nil).Once()
		result, err := svc.PickLatest(ctx, domain.UserLvlDefault)
		assert.NoError(t, err)
		assert.Equal(t, expected.Content, result.Content)
		mockRepo.AssertExpectations(t)
	})

	t.Run("no posts", func(t *testing.T) {
		mockRepo.On("LatestPostByAudience", ctx, domain.UserLvlDefault).Return((*domain.Post)(nil), domain.ErrNoPosts).Once()
		result, err := svc.PickLatest(ctx, domain.UserLvlDefault)
		assert.ErrorIs(t, err, domain.ErrNoPosts)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("LatestPostByAudience", ctx, domain.UserLvlDefault).Return((*domain.Post)(nil), assert.AnError).Once()
		result, err := svc.PickLatest(ctx, domain.UserLvlDefault)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestPostService_MarkAsPosted(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockPostRepo)
	svc := post.New(testutils.NewTestLogger(), mockRepo)

	t.Run("succes", func(t *testing.T) {
		mockRepo.On("MarkAsPosted", ctx, int64(1)).Return(nil).Once()
		err := svc.MarkAsPosted(ctx, int64(1))
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("MarkAsPosted", ctx, int64(1)).Return(assert.AnError).Once()
		err := svc.MarkAsPosted(ctx, int64(1))
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("post not found", func(t *testing.T) {
		mockRepo.On("MarkAsPosted", ctx, int64(1)).Return(domain.ErrPostNotFound).Once()
		err := svc.MarkAsPosted(ctx, int64(1))
		assert.ErrorIs(t, err, domain.ErrPostNotFound)
		mockRepo.AssertExpectations(t)
	})
}

type mockPostRepo struct {
	mock.Mock
}

func (r *mockPostRepo) SavePost(ctx context.Context, post postRepo.CreatePostInput) error {
	args := r.Called(ctx, post)
	return args.Error(0)
}

func (r *mockPostRepo) MarkAsPosted(ctx context.Context, id int64) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *mockPostRepo) LatestPostByAudience(ctx context.Context, audience domain.UserLvl) (*domain.Post, error) {
	args := r.Called(ctx, audience)
	return args.Get(0).(*domain.Post), args.Error(1)
}
