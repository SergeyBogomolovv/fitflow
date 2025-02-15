package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	authSvc "github.com/SergeyBogomolovv/fitflow/internal/service/auth"
	"github.com/SergeyBogomolovv/fitflow/internal/service/auth/mocks"
	"github.com/SergeyBogomolovv/fitflow/pkg/auth"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Login(t *testing.T) {
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
		want         string
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				login:    "login",
				password: "password",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				hashedPassword, err := auth.HashPassword(args.password)
				require.NoError(t, err)
				repo.EXPECT().AdminByLogin(args.ctx, args.login).Return(domain.Admin{Password: hashedPassword}, nil).Once()
			},
			want:    string(mock.AnythingOfType("string")),
			wantErr: nil,
		},
		{
			name: "admin not found",
			args: args{
				ctx:      context.Background(),
				login:    "login",
				password: "password",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				repo.EXPECT().AdminByLogin(args.ctx, args.login).Return(domain.Admin{}, domain.ErrAdminNotFound).Once()
			},
			want:    "",
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			args: args{
				ctx:      context.Background(),
				login:    "login",
				password: "password",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				repo.EXPECT().AdminByLogin(args.ctx, args.login).Return(domain.Admin{Password: []byte("hash")}, nil).Once()
			},
			want:    "",
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "failed to get admin",
			args: args{
				ctx:      context.Background(),
				login:    "login",
				password: "password",
			},
			mockBehavior: func(repo *mocks.AdminRepo, args args) {
				repo.EXPECT().AdminByLogin(args.ctx, args.login).Return(domain.Admin{}, assert.AnError).Once()
			},
			want:    "",
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewAdminRepo(t)
			tc.mockBehavior(repo, tc.args)

			svc := authSvc.New(testutils.NewTestLogger(), repo, []byte("secret"), time.Hour)
			got, err := svc.Login(tc.args.ctx, tc.args.login, tc.args.password)

			if tc.wantErr != nil {
				assert.ErrorIs(t, tc.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}
