package admin_test

import (
	"context"
	"testing"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	adminSvc "github.com/SergeyBogomolovv/fitflow/internal/service/admin"
	"github.com/SergeyBogomolovv/fitflow/internal/service/admin/mocks"
	"github.com/SergeyBogomolovv/fitflow/pkg/auth"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAdminService_CreateAdmin(t *testing.T) {
	type args struct {
		ctx      context.Context
		login    string
		password string
	}

	type MockBehavior func(repo *mocks.AdminRepo, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         error
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				login:    "login",
				password: "password",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				repo.EXPECT().AdminExists(args.ctx, args.login).Return(false, nil).Once()
				repo.EXPECT().SaveAdmin(args.ctx, args.login, mock.Anything).Return(nil).Once()
			},
			want: nil,
		},
		{
			name: "admin already exists",
			args: args{
				ctx:      context.Background(),
				login:    "login",
				password: "password",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				repo.EXPECT().AdminExists(args.ctx, args.login).Return(true, nil).Once()
			},
			want: domain.ErrAdminAlreadyExists,
		},
		{
			name: "failed to save admin",
			args: args{
				ctx:      context.Background(),
				login:    "login",
				password: "password",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				repo.EXPECT().AdminExists(args.ctx, args.login).Return(false, nil).Once()
				repo.EXPECT().SaveAdmin(args.ctx, args.login, mock.Anything).Return(assert.AnError).Once()
			},
			want: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewAdminRepo(t)
			adminSvc := adminSvc.New(testutils.NewTestLogger(), repo)
			tc.mockBehavior(repo, tc.args)

			err := adminSvc.CreateAdmin(tc.args.ctx, tc.args.login, tc.args.password)
			assert.ErrorIs(t, tc.want, err)
		})
	}
}

func TestAdminService_UpdatePassword(t *testing.T) {
	type args struct {
		ctx     context.Context
		login   string
		oldPass string
		newPass string
	}

	type MockBehavior func(repo *mocks.AdminRepo, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         error
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				login:   "login",
				oldPass: "oldPass",
				newPass: "newPass",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				hashedPassword, err := auth.HashPassword(args.oldPass)
				require.NoError(t, err)
				repo.EXPECT().AdminByLogin(args.ctx, args.login).Return(
					domain.Admin{
						Login:    args.login,
						Password: hashedPassword,
					}, nil).Once()
				repo.EXPECT().UpdatePassword(args.ctx, args.login, mock.Anything).Return(nil).Once()
			},
			want: nil,
		},
		{
			name: "invalid login",
			args: args{
				ctx:     context.Background(),
				login:   "login",
				oldPass: "oldPass",
				newPass: "newPass",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				repo.EXPECT().AdminByLogin(args.ctx, args.login).Return(domain.Admin{}, domain.ErrAdminNotFound).Once()
			},
			want: domain.ErrInvalidCredentials,
		},
		{
			name: "invalid password",
			args: args{
				ctx:     context.Background(),
				login:   "login",
				oldPass: "oldPass",
				newPass: "newPass",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				repo.EXPECT().AdminByLogin(args.ctx, args.login).Return(
					domain.Admin{
						Login:    args.login,
						Password: []byte("wrongPassword"),
					}, nil).Once()
			},
			want: domain.ErrInvalidCredentials,
		},
		{
			name: "failed to update password",
			args: args{
				ctx:     context.Background(),
				login:   "login",
				oldPass: "oldPass",
				newPass: "newPass",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				hashedPassword, err := auth.HashPassword(args.oldPass)
				require.NoError(t, err)
				repo.EXPECT().AdminByLogin(args.ctx, args.login).Return(
					domain.Admin{
						Login:    args.login,
						Password: hashedPassword,
					}, nil).Once()
				repo.EXPECT().UpdatePassword(args.ctx, args.login, mock.Anything).Return(assert.AnError).Once()
			},
			want: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewAdminRepo(t)
			adminSvc := adminSvc.New(testutils.NewTestLogger(), repo)
			tc.mockBehavior(repo, tc.args)

			err := adminSvc.UpdatePassword(tc.args.ctx, tc.args.login, tc.args.oldPass, tc.args.newPass)
			assert.ErrorIs(t, tc.want, err)
		})
	}
}

func TestAdminService_RemoveAdmin(t *testing.T) {
	type args struct {
		ctx   context.Context
		login string
	}

	type MockBehavior func(repo *mocks.AdminRepo, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         error
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				login: "login",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				repo.EXPECT().DeleteAdmin(args.ctx, args.login).Return(nil).Once()
			},
			want: nil,
		},
		{
			name: "admin not found",
			args: args{
				ctx:   context.Background(),
				login: "login",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				repo.EXPECT().DeleteAdmin(args.ctx, args.login).Return(domain.ErrAdminNotFound).Once()
			},
			want: domain.ErrAdminNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewAdminRepo(t)
			adminSvc := adminSvc.New(testutils.NewTestLogger(), repo)
			tc.mockBehavior(repo, tc.args)

			err := adminSvc.RemoveAdmin(tc.args.ctx, tc.args.login)
			assert.ErrorIs(t, tc.want, err)
		})
	}
}
