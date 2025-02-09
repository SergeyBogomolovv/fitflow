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

	t.Run("success", func(t *testing.T) {
		mockRepo.On("SaveUser", ctx, int64(1), domain.UserLvlDefault).Return(&domain.User{ID: 1}, nil)
		err := svc.SaveUser(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.On("SaveUser", ctx, int64(2), domain.UserLvlDefault).Return((*domain.User)(nil), assert.AnError)
		err := svc.SaveUser(ctx, 2)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateSubscribed(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := userSvc.New(testutils.NewTestLogger(), mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UpdateSubscribed", ctx, int64(1), true).Return(nil)
		err := svc.UpdateSubscribed(ctx, 1, true)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("UpdateSubscribed", ctx, int64(2), false).Return(domain.ErrUserNotFound)
		err := svc.UpdateSubscribed(ctx, 2, false)
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
