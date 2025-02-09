package admin_test

import (
	"context"
	"testing"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	adminSvc "github.com/SergeyBogomolovv/fitflow/internal/service/admin"
	"github.com/SergeyBogomolovv/fitflow/pkg/auth"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdminService_CreateAdmin(t *testing.T) {
	mockRepo := new(mockAdminRepo)
	svc := adminSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("AdminExists", ctx, "success").Return(false, nil)
		mockRepo.On("SaveAdmin", ctx, "success", mock.Anything).Return(nil)

		err := svc.CreateAdmin(ctx, "success", "password")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("login already exists", func(t *testing.T) {
		mockRepo.On("AdminExists", ctx, "existing").Return(true, nil)

		err := svc.CreateAdmin(ctx, "existing", "password")

		assert.ErrorIs(t, err, domain.ErrAdminAlreadyExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SaveAdmin error", func(t *testing.T) {
		mockRepo.On("AdminExists", ctx, "error").Return(false, nil)
		mockRepo.On("SaveAdmin", ctx, "error", mock.Anything).Return(assert.AnError)

		err := svc.CreateAdmin(ctx, "error", "password")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("AdminExists error", func(t *testing.T) {
		mockRepo.On("AdminExists", ctx, "error").Return(false, assert.AnError)

		err := svc.CreateAdmin(ctx, "error", "password")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAdminService_UpdatePassword(t *testing.T) {
	mockRepo := new(mockAdminRepo)
	svc := adminSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		hashed, err := auth.HashPassword("old")
		assert.NoError(t, err)
		mockRepo.On("AdminByLogin", ctx, "success").Return(&domain.Admin{Login: "success", Password: hashed}, nil)
		mockRepo.On("UpdatePassword", ctx, "success", mock.Anything).Return(nil)

		err = svc.UpdatePassword(ctx, "success", "old", "password")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong old password", func(t *testing.T) {
		mockRepo.On("AdminByLogin", ctx, "login").Return(&domain.Admin{Password: []byte("diff")}, nil)

		err := svc.UpdatePassword(ctx, "login", "old", "password")

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})

	t.Run("admin not found", func(t *testing.T) {
		mockRepo.On("AdminByLogin", ctx, "notexists").Return((*domain.Admin)(nil), domain.ErrAdminNotFound)

		err := svc.UpdatePassword(ctx, "notexists", "old", "new")

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		hashed, err := auth.HashPassword("old")
		assert.NoError(t, err)
		mockRepo.On("AdminByLogin", ctx, "err").Return(&domain.Admin{Password: hashed}, nil)
		mockRepo.On("UpdatePassword", ctx, "err", mock.Anything).Return(assert.AnError)

		err = svc.UpdatePassword(ctx, "err", "old", "password")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAdminService_RemoveAdmin(t *testing.T) {
	mockRepo := new(mockAdminRepo)
	svc := adminSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteAdmin", ctx, "success").Return(nil)
		err := svc.RemoveAdmin(ctx, "success")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("admin not found", func(t *testing.T) {
		mockRepo.On("DeleteAdmin", ctx, "notexists").Return(domain.ErrAdminNotFound)
		err := svc.RemoveAdmin(ctx, "notexists")
		assert.ErrorIs(t, err, domain.ErrAdminNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("DeleteAdmin", ctx, "err").Return(assert.AnError)
		err := svc.RemoveAdmin(ctx, "err")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

type mockAdminRepo struct {
	mock.Mock
}

func (m *mockAdminRepo) SaveAdmin(ctx context.Context, login string, password []byte) error {
	args := m.Called(ctx, login, password)
	return args.Error(0)
}

func (m *mockAdminRepo) AdminByLogin(ctx context.Context, login string) (*domain.Admin, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(*domain.Admin), args.Error(1)
}

func (m *mockAdminRepo) AdminExists(ctx context.Context, login string) (bool, error) {
	args := m.Called(ctx, login)
	return args.Bool(0), args.Error(1)
}

func (m *mockAdminRepo) UpdatePassword(ctx context.Context, login string, password []byte) error {
	args := m.Called(ctx, login, password)
	return args.Error(0)
}

func (m *mockAdminRepo) DeleteAdmin(ctx context.Context, login string) error {
	args := m.Called(ctx, login)
	return args.Error(0)
}
