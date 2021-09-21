// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// BcryptWrapper is an autogenerated mock type for the BcryptWrapper type
type BcryptWrapper struct {
	mock.Mock
}

// CompareHashAndPassword provides a mock function with given fields: hashedPassword, password
func (_m *BcryptWrapper) CompareHashAndPassword(hashedPassword string, password string) error {
	ret := _m.Called(hashedPassword, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(hashedPassword, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HashPassword provides a mock function with given fields: password
func (_m *BcryptWrapper) HashPassword(password string) (string, error) {
	ret := _m.Called(password)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}