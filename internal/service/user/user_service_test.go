package user_test

import (
	"context"
	"testing"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	userSvc "github.com/SergeyBogomolovv/fitflow/internal/service/user"
	"github.com/SergeyBogomolovv/fitflow/internal/service/user/mocks"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestUserService_EnsureUserExists(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}

	type MockBehavior func(repo *mocks.UserRepo, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         error
	}{
		{
			name: "not exists, need to save",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().UserExists(args.ctx, args.id).Return(false, nil).Once()
				repo.EXPECT().SaveUser(args.ctx, args.id, domain.UserLvlDefault).Return(nil).Once()
			},
			want: nil,
		},
		{
			name: "already exists",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().UserExists(args.ctx, args.id).Return(true, nil).Once()
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewUserRepo(t)
			tc.mockBehavior(repo, tc.args)
			svc := userSvc.New(testutils.NewTestLogger(), repo)
			err := svc.EnsureUserExists(tc.args.ctx, tc.args.id)

			if tc.want == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func TestUserService_UpdateSubscribed(t *testing.T) {
	type args struct {
		ctx        context.Context
		id         int64
		subscribed bool
	}

	type MockBehavior func(repo *mocks.UserRepo, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         error
	}{
		{
			name: "success",
			args: args{
				ctx:        context.Background(),
				id:         1,
				subscribed: true,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().UpdateSubscribed(args.ctx, args.id, args.subscribed).Return(nil).Once()
			},
			want: nil,
		},
		{
			name: "user not found",
			args: args{
				ctx:        context.Background(),
				id:         1,
				subscribed: true,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().UpdateSubscribed(args.ctx, args.id, args.subscribed).Return(domain.ErrUserNotFound).Once()
			},
			want: domain.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewUserRepo(t)
			tc.mockBehavior(repo, tc.args)
			svc := userSvc.New(testutils.NewTestLogger(), repo)
			err := svc.UpdateSubscribed(tc.args.ctx, tc.args.id, tc.args.subscribed)

			if tc.want == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func TestUserService_UpdateUserLvl(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
		lvl domain.UserLvl
	}

	type MockBehavior func(repo *mocks.UserRepo, args args)

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
				lvl: domain.UserLvlBeginner,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().UpdateUserLvl(args.ctx, args.id, args.lvl).Return(nil).Once()
			},
			want: nil,
		},
		{
			name: "user not found",
			args: args{
				ctx: context.Background(),
				id:  1,
				lvl: domain.UserLvlBeginner,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().UpdateUserLvl(args.ctx, args.id, args.lvl).Return(domain.ErrUserNotFound).Once()
			},
			want: domain.ErrUserNotFound,
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				id:  1,
				lvl: domain.UserLvlBeginner,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().UpdateUserLvl(args.ctx, args.id, args.lvl).Return(assert.AnError).Once()
			},
			want: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewUserRepo(t)
			tc.mockBehavior(repo, tc.args)
			svc := userSvc.New(testutils.NewTestLogger(), repo)
			err := svc.UpdateUserLvl(tc.args.ctx, tc.args.id, tc.args.lvl)

			if tc.want == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func TestUserService_SubscribersIds(t *testing.T) {
	type args struct {
		ctx context.Context
		lvl domain.UserLvl
	}

	type MockBehavior func(repo *mocks.UserRepo, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         []int64
		wantErr      error
	}{
		{
			name: "default lvl",
			args: args{
				ctx: context.Background(),
				lvl: domain.UserLvlDefault,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().Subscribers(args.ctx, args.lvl, true).Return([]domain.User{{ID: 1}}, nil).Once()
			},
			want: []int64{1},
		},
		{
			name: "beginner lvl",
			args: args{
				ctx: context.Background(),
				lvl: domain.UserLvlBeginner,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().Subscribers(args.ctx, args.lvl, false).Return([]domain.User{{ID: 1}}, nil).Once()
			},
			want: []int64{1},
		},
		{
			name: "no subscribers",
			args: args{
				ctx: context.Background(),
				lvl: domain.UserLvlBeginner,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().Subscribers(args.ctx, args.lvl, false).Return([]domain.User{}, nil).Once()
			},
			want: []int64{},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				lvl: domain.UserLvlBeginner,
			},
			mockBehavior: func(repo *mocks.UserRepo, args args) {
				repo.EXPECT().Subscribers(args.ctx, args.lvl, false).Return([]domain.User{}, assert.AnError).Once()
			},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewUserRepo(t)
			tc.mockBehavior(repo, tc.args)
			svc := userSvc.New(testutils.NewTestLogger(), repo)
			got, err := svc.SubscribersIds(tc.args.ctx, tc.args.lvl)

			if tc.wantErr == nil {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tc.want, got)
				return
			}
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
