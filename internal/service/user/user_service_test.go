package user_test

import (
	"context"
	"testing"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	userSvc "github.com/SergeyBogomolovv/fitflow/internal/service/user"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_AddUser(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("SaveUser", ctx, int64(1), domain.UserLvlDefault).Return(&domain.User{ID: 1}, nil)
		err := svc.AddUser(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("SaveUser", ctx, int64(2), domain.UserLvlDefault).Return((*domain.User)(nil), assert.AnError)
		err := svc.AddUser(ctx, 2)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_RemoveUser(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteUser", ctx, int64(1)).Return(nil)
		err := svc.RemoveUser(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("DeleteUser", ctx, int64(2)).Return(domain.ErrUserNotFound)
		err := svc.RemoveUser(ctx, 2)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("DeleteUser", ctx, int64(3)).Return(assert.AnError)
		err := svc.RemoveUser(ctx, 3)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUserLvl(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()
	level := domain.UserLvlAdvanced

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UpdateUserLvl", ctx, int64(1), level).Return(nil)
		err := svc.UpdateUserLvl(ctx, int64(1), level)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("UpdateUserLvl", ctx, int64(2), level).Return(domain.ErrUserNotFound)
		err := svc.UpdateUserLvl(ctx, int64(2), level)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("UpdateUserLvl", ctx, int64(3), level).Return(assert.AnError)
		err := svc.UpdateUserLvl(ctx, int64(3), level)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UserIdsByLvl(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()
	users := []domain.User{{ID: 1}, {ID: 2}}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UsersByLvl", ctx, domain.UserLvlDefault).Return(users, nil)
		ids, err := svc.UserIdsByLvl(ctx, domain.UserLvlDefault)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []int64{1, 2}, ids)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo.On("UsersByLvl", ctx, domain.UserLvlBeginner).Return([]domain.User{}, nil)
		ids, err := svc.UserIdsByLvl(ctx, domain.UserLvlBeginner)
		assert.NoError(t, err)
		assert.Empty(t, ids)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("UsersByLvl", ctx, domain.UserLvlIntermediate).Return(([]domain.User)(nil), assert.AnError)
		ids, err := svc.UserIdsByLvl(ctx, domain.UserLvlIntermediate)
		assert.Error(t, err)
		assert.Nil(t, ids)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_AllUserIds(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()
	users := []domain.User{{ID: 1}, {ID: 2}}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("AllUsers", ctx).Return(users, nil).Once()
		ids, err := svc.AllUserIds(ctx)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []int64{1, 2}, ids)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo.On("AllUsers", ctx).Return([]domain.User{}, nil).Once()
		ids, err := svc.AllUserIds(ctx)
		assert.NoError(t, err)
		assert.Empty(t, ids)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("AllUsers", ctx).Return(([]domain.User)(nil), assert.AnError).Once()
		ids, err := svc.AllUserIds(ctx)
		assert.Error(t, err)
		assert.Nil(t, ids)
		mockRepo.AssertExpectations(t)
	})
}

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) SaveUser(ctx context.Context, id int64, lvl domain.UserLvl) (*domain.User, error) {
	args := m.Called(ctx, id, lvl)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error {
	args := m.Called(ctx, id, lvl)
	return args.Error(0)
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserRepo) UsersByLvl(ctx context.Context, lvl domain.UserLvl) ([]domain.User, error) {
	args := m.Called(ctx, lvl)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *mockUserRepo) AllUsers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}
