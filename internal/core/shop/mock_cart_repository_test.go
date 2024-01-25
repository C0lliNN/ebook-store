// Code generated by mockery v2.40.1. DO NOT EDIT.

package shop

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockCartRepository is an autogenerated mock type for the CartRepository type
type MockCartRepository struct {
	mock.Mock
}

// DeleteByUserID provides a mock function with given fields: ctx, userID
func (_m *MockCartRepository) DeleteByUserID(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByUserID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByUserID provides a mock function with given fields: ctx, userID
func (_m *MockCartRepository) FindByUserID(ctx context.Context, userID string) (*Cart, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for FindByUserID")
	}

	var r0 *Cart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*Cart, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *Cart); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Cart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, cart
func (_m *MockCartRepository) Save(ctx context.Context, cart *Cart) error {
	ret := _m.Called(ctx, cart)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *Cart) error); ok {
		r0 = rf(ctx, cart)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockCartRepository creates a new instance of MockCartRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCartRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCartRepository {
	mock := &MockCartRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}