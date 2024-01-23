// Code generated by mockery v2.40.1. DO NOT EDIT.

package shop

import mock "github.com/stretchr/testify/mock"

// MockIDGenerator is an autogenerated mock type for the IDGenerator type
type MockIDGenerator struct {
	mock.Mock
}

// NewID provides a mock function with given fields:
func (_m *MockIDGenerator) NewID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewMockIDGenerator creates a new instance of MockIDGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIDGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIDGenerator {
	mock := &MockIDGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
