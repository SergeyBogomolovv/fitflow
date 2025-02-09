package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	authSvc "github.com/SergeyBogomolovv/fitflow/internal/service/auth"
	"github.com/SergeyBogomolovv/fitflow/pkg/auth"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(adminRepoMock)
	svc := authSvc.New(testutils.NewTestLogger(), mockRepo, []byte("secret"), time.Hour)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		hashed, err := auth.HashPassword("password")
		require.NoError(t, err)

		mockRepo.On("AdminByLogin", ctx, "success").Return(&domain.Admin{Login: "success", Password: hashed}, nil)

		token, err := svc.Login(ctx, "success", "password")
		res, err := auth.VerifyJWT(token, []byte("secret"))
		require.NoError(t, err)

		assert.Equal(t, res, "success")
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		hashed, err := auth.HashPassword("password")
		require.NoError(t, err)

		mockRepo.On("AdminByLogin", ctx, "wrongpass").Return(&domain.Admin{Login: "wrongpass", Password: hashed}, nil)

		_, err = svc.Login(ctx, "wrongpass", "wrong")

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong login", func(t *testing.T) {
		mockRepo.On("AdminByLogin", ctx, "notexists").Return((*domain.Admin)(nil), domain.ErrAdminNotFound)

		_, err := svc.Login(ctx, "notexists", "password")

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("AdminByLogin", ctx, "error").Return((*domain.Admin)(nil), assert.AnError)

		_, err := svc.Login(ctx, "error", "password")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

type adminRepoMock struct {
	mock.Mock
}

func (m *adminRepoMock) AdminByLogin(ctx context.Context, login string) (*domain.Admin, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(*domain.Admin), args.Error(1)
}
