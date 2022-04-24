// Code generated by mockery v2.12.0. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/c0llinn/ebook-store/internal/auth"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// FindByEmail provides a mock function with given fields: ctx, email
func (_m *Repository) FindByEmail(ctx context.Context, email string) (auth.User, error) {
	ret := _m.Called(ctx, email)

	var r0 auth.User
	if rf, ok := ret.Get(0).(func(context.Context, string) auth.User); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(auth.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, user
func (_m *Repository) Save(ctx context.Context, user *auth.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *auth.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, user
func (_m *Repository) Update(ctx context.Context, user *auth.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *auth.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t testing.TB) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
