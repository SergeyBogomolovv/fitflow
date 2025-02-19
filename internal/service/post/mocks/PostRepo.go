// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/SergeyBogomolovv/fitflow/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// PostRepo is an autogenerated mock type for the PostRepo type
type PostRepo struct {
	mock.Mock
}

type PostRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *PostRepo) EXPECT() *PostRepo_Expecter {
	return &PostRepo_Expecter{mock: &_m.Mock}
}

// LatestByAudience provides a mock function with given fields: ctx, audience
func (_m *PostRepo) LatestByAudience(ctx context.Context, audience domain.UserLvl) (domain.Post, error) {
	ret := _m.Called(ctx, audience)

	if len(ret) == 0 {
		panic("no return value specified for LatestByAudience")
	}

	var r0 domain.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserLvl) (domain.Post, error)); ok {
		return rf(ctx, audience)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserLvl) domain.Post); ok {
		r0 = rf(ctx, audience)
	} else {
		r0 = ret.Get(0).(domain.Post)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.UserLvl) error); ok {
		r1 = rf(ctx, audience)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PostRepo_LatestByAudience_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LatestByAudience'
type PostRepo_LatestByAudience_Call struct {
	*mock.Call
}

// LatestByAudience is a helper method to define mock.On call
//   - ctx context.Context
//   - audience domain.UserLvl
func (_e *PostRepo_Expecter) LatestByAudience(ctx interface{}, audience interface{}) *PostRepo_LatestByAudience_Call {
	return &PostRepo_LatestByAudience_Call{Call: _e.mock.On("LatestByAudience", ctx, audience)}
}

func (_c *PostRepo_LatestByAudience_Call) Run(run func(ctx context.Context, audience domain.UserLvl)) *PostRepo_LatestByAudience_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.UserLvl))
	})
	return _c
}

func (_c *PostRepo_LatestByAudience_Call) Return(_a0 domain.Post, _a1 error) *PostRepo_LatestByAudience_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PostRepo_LatestByAudience_Call) RunAndReturn(run func(context.Context, domain.UserLvl) (domain.Post, error)) *PostRepo_LatestByAudience_Call {
	_c.Call.Return(run)
	return _c
}

// MarkAsPosted provides a mock function with given fields: ctx, id
func (_m *PostRepo) MarkAsPosted(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for MarkAsPosted")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PostRepo_MarkAsPosted_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MarkAsPosted'
type PostRepo_MarkAsPosted_Call struct {
	*mock.Call
}

// MarkAsPosted is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *PostRepo_Expecter) MarkAsPosted(ctx interface{}, id interface{}) *PostRepo_MarkAsPosted_Call {
	return &PostRepo_MarkAsPosted_Call{Call: _e.mock.On("MarkAsPosted", ctx, id)}
}

func (_c *PostRepo_MarkAsPosted_Call) Run(run func(ctx context.Context, id int64)) *PostRepo_MarkAsPosted_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *PostRepo_MarkAsPosted_Call) Return(_a0 error) *PostRepo_MarkAsPosted_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PostRepo_MarkAsPosted_Call) RunAndReturn(run func(context.Context, int64) error) *PostRepo_MarkAsPosted_Call {
	_c.Call.Return(run)
	return _c
}

// NewPostRepo creates a new instance of PostRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostRepo {
	mock := &PostRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
