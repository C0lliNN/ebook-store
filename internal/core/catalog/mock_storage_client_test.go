// Code generated by mockery v2.40.1. DO NOT EDIT.

package catalog

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockStorageClient is an autogenerated mock type for the StorageClient type
type MockStorageClient struct {
	mock.Mock
}

// GenerateGetPreSignedUrl provides a mock function with given fields: ctx, key
func (_m *MockStorageClient) GenerateGetPreSignedUrl(ctx context.Context, key string) (string, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for GenerateGetPreSignedUrl")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GeneratePutPreSignedUrl provides a mock function with given fields: ctx, key
func (_m *MockStorageClient) GeneratePutPreSignedUrl(ctx context.Context, key string) (string, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for GeneratePutPreSignedUrl")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockStorageClient creates a new instance of MockStorageClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockStorageClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockStorageClient {
	mock := &MockStorageClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}