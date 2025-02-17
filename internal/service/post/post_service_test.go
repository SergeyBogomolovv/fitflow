package post_test

import (
	"context"
	"testing"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	postSvc "github.com/SergeyBogomolovv/fitflow/internal/service/post"
	"github.com/SergeyBogomolovv/fitflow/internal/service/post/mocks"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestPostService_PickLatest(t *testing.T) {
	type args struct {
		ctx      context.Context
		audience domain.UserLvl
	}

	type MockBehavior func(repo *mocks.PostRepo, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         domain.Post
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				audience: domain.UserLvlBeginner,
			},
			mockBehavior: func(repo *mocks.PostRepo, args args) {
				repo.EXPECT().LatestByAudience(args.ctx, args.audience).Return(domain.Post{ID: 1}, nil).Once()
			},
			want: domain.Post{ID: 1},
		},
		{
			name: "no posts",
			args: args{
				ctx:      context.Background(),
				audience: domain.UserLvlBeginner,
			},
			mockBehavior: func(repo *mocks.PostRepo, args args) {
				repo.EXPECT().LatestByAudience(args.ctx, args.audience).Return(domain.Post{}, domain.ErrPostNotFound).Once()
			},
			wantErr: domain.ErrPostNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPostRepo(t)
			tc.mockBehavior(repo, tc.args)

			svc := postSvc.New(testutils.NewTestLogger(), repo)
			got, err := svc.PickLatest(tc.args.ctx, tc.args.audience)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPostService_MarkAsPosted(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}

	type MockBehavior func(repo *mocks.PostRepo, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(repo *mocks.PostRepo, args args) {
				repo.EXPECT().MarkAsPosted(args.ctx, args.id).Return(nil).Once()
			},
			want: nil,
		},
		{
			name: "post not found",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(repo *mocks.PostRepo, args args) {
				repo.EXPECT().MarkAsPosted(args.ctx, args.id).Return(domain.ErrPostNotFound).Once()
			},
			want: domain.ErrPostNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPostRepo(t)
			tc.mockBehavior(repo, tc.args)
			svc := postSvc.New(testutils.NewTestLogger(), repo)

			got := svc.MarkAsPosted(tc.args.ctx, tc.args.id)
			assert.Equal(t, tc.want, got)
		})
	}
}
