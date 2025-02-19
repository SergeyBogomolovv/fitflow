package content_test

import (
	"context"
	"mime/multipart"
	"testing"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	postRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/post"
	"github.com/SergeyBogomolovv/fitflow/internal/service/content"
	"github.com/SergeyBogomolovv/fitflow/internal/service/content/mocks"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContentService_CreatePost(t *testing.T) {
	type MockBehavior func(repo *mocks.PostRepo, s3 *mocks.S3Client, in domain.CreatePostDTO)

	testCases := []struct {
		name         string
		in           domain.CreatePostDTO
		mockBehavior MockBehavior
		want         domain.Post
		wantErr      bool
	}{
		{
			name: "success",
			in: domain.CreatePostDTO{
				Content:  "test content",
				Audience: domain.UserLvlBeginner,
				Images: []*multipart.FileHeader{
					testutils.CreateTestFile(t, "test.jpg", "test content"),
				},
			},
			mockBehavior: func(repo *mocks.PostRepo, s3 *mocks.S3Client, in domain.CreatePostDTO) {
				s3.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return("test.jpg", nil).Once()
				repo.EXPECT().Save(mock.Anything, postRepo.SavePostInput{
					Content:  in.Content,
					Images:   []string{"test.jpg"},
					Audience: in.Audience,
				}).Return(domain.Post{ID: 1}, nil).Once()
			},
			want:    domain.Post{ID: 1},
			wantErr: false,
		},
		{
			name: "no images",
			in: domain.CreatePostDTO{
				Content:  "test content",
				Audience: domain.UserLvlBeginner,
				Images:   []*multipart.FileHeader{},
			},
			mockBehavior: func(repo *mocks.PostRepo, s3 *mocks.S3Client, in domain.CreatePostDTO) {
				repo.EXPECT().Save(mock.Anything, postRepo.SavePostInput{
					Content:  in.Content,
					Images:   []string{},
					Audience: in.Audience,
				}).Return(domain.Post{ID: 1}, nil).Once()
			},
			want:    domain.Post{ID: 1},
			wantErr: false,
		},
		{
			name: "failed to upload",
			in: domain.CreatePostDTO{
				Content:  "test content",
				Audience: domain.UserLvlBeginner,
				Images: []*multipart.FileHeader{
					testutils.CreateTestFile(t, "test.jpg", "test content"),
				},
			},
			mockBehavior: func(repo *mocks.PostRepo, s3 *mocks.S3Client, in domain.CreatePostDTO) {
				s3.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return("", assert.AnError).Once()
			},
			wantErr: true,
		},
		{
			name: "failed to save",
			in: domain.CreatePostDTO{
				Content:  "test content",
				Audience: domain.UserLvlBeginner,
				Images: []*multipart.FileHeader{
					testutils.CreateTestFile(t, "test.jpg", "test content"),
				},
			},
			mockBehavior: func(repo *mocks.PostRepo, s3 *mocks.S3Client, in domain.CreatePostDTO) {
				s3.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return("test.jpg", nil).Once()
				repo.EXPECT().Save(mock.Anything, postRepo.SavePostInput{
					Content:  in.Content,
					Images:   []string{"test.jpg"},
					Audience: in.Audience,
				}).Return(domain.Post{}, assert.AnError).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPostRepo(t)
			s3 := mocks.NewS3Client(t)
			tc.mockBehavior(repo, s3, tc.in)

			svc := content.New(testutils.NewTestLogger(), repo, nil, s3)
			got, err := svc.CreatePost(context.Background(), tc.in)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestContentService_RemovePost(t *testing.T) {
	type MockBehavior func(repo *mocks.PostRepo, s3 *mocks.S3Client, id int64)

	testCases := []struct {
		name         string
		id           int64
		mockBehavior MockBehavior
		want         error
	}{
		{
			name: "success",
			id:   1,
			mockBehavior: func(repo *mocks.PostRepo, s3 *mocks.S3Client, id int64) {
				repo.EXPECT().Remove(mock.Anything, id).Return(domain.Post{Images: []string{"test.jpg"}}, nil).Once()
				s3.EXPECT().Delete(mock.Anything, "test.jpg").Return(nil).Once()
			},
			want: nil,
		},
		{
			name: "post not found",
			id:   1,
			mockBehavior: func(repo *mocks.PostRepo, s3 *mocks.S3Client, id int64) {
				repo.EXPECT().Remove(mock.Anything, id).Return(domain.Post{}, domain.ErrPostNotFound).Once()
			},
			want: domain.ErrPostNotFound,
		},
		{
			name: "failed to delete image",
			id:   1,
			mockBehavior: func(repo *mocks.PostRepo, s3 *mocks.S3Client, id int64) {
				repo.EXPECT().Remove(mock.Anything, id).Return(domain.Post{Images: []string{"test.jpg"}}, nil).Once()
				s3.EXPECT().Delete(mock.Anything, "test.jpg").Return(assert.AnError).Once()
			},
			want: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPostRepo(t)
			s3 := mocks.NewS3Client(t)
			tc.mockBehavior(repo, s3, tc.id)

			svc := content.New(testutils.NewTestLogger(), repo, nil, s3)
			got := svc.RemovePost(context.Background(), tc.id)
			assert.ErrorIs(t, got, tc.want)
		})
	}
}
