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

func TestUserService_SaveUser(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()
	userId := int64(1)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UserExists", ctx, userId).Return(false, nil).Once()
		mockRepo.On("SaveUser", ctx, userId, domain.UserLvlDefault).Return(&domain.User{ID: userId}, nil).Once()
		err := svc.SaveUser(ctx, userId)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockRepo.On("UserExists", ctx, userId).Return(true, nil).Once()
		err := svc.SaveUser(ctx, userId)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("UserExists", ctx, userId).Return(false, nil).Once()
		mockRepo.On("SaveUser", ctx, userId, domain.UserLvlDefault).Return((*domain.User)(nil), assert.AnError).Once()
		err := svc.SaveUser(ctx, userId)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateSubscribed(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()
	userId := int64(1)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UpdateSubscribed", ctx, userId, true).Return(nil).Once()
		err := svc.UpdateSubscribed(ctx, userId, true)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("UpdateSubscribed", ctx, userId, false).Return(domain.ErrUserNotFound).Once()
		err := svc.UpdateSubscribed(ctx, userId, false)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_SubscribersIdsByLvl(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()
	users := []domain.User{{ID: 1}, {ID: 2}}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("SubscribersByLvl", ctx, domain.UserLvlDefault).Return(users, nil)
		ids, err := svc.SubscribersIdsByLvl(ctx, domain.UserLvlDefault)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []int64{1, 2}, ids)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUserLvl(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()
	userId := int64(1)
	newLvl := domain.UserLvlAdvanced

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UpdateUserLvl", ctx, userId, newLvl).Return(nil).Once()
		err := svc.UpdateUserLvl(ctx, userId, newLvl)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("UpdateUserLvl", ctx, userId, newLvl).Return(domain.ErrUserNotFound).Once()
		err := svc.UpdateUserLvl(ctx, userId, newLvl)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_SubscribersIds(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()
	users := []domain.User{{ID: 1}, {ID: 2}, {ID: 3}}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Subscribers", ctx).Return(users, nil).Once()
		ids, err := svc.SubscribersIds(ctx)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []int64{1, 2, 3}, ids)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("Subscribers", ctx).Return(([]domain.User)(nil), assert.AnError).Once()
		ids, err := svc.SubscribersIds(ctx)
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

func (m *mockUserRepo) UpdateSubscribed(ctx context.Context, id int64, subscribed bool) error {
	args := m.Called(ctx, id, subscribed)
	return args.Error(0)
}

func (m *mockUserRepo) UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error {
	args := m.Called(ctx, id, lvl)
	return args.Error(0)
}

func (m *mockUserRepo) SubscribersByLvl(ctx context.Context, lvl domain.UserLvl) ([]domain.User, error) {
	args := m.Called(ctx, lvl)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *mockUserRepo) Subscribers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *mockUserRepo) UserExists(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
