// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mock

import (
	model "github.com/c0llinn/ebook-store/internal/auth/model"
	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the Repository type
type UserRepository struct {
	mock.Mock
}

// Save provides a mock function with given fields: user
func (_m *UserRepository) Save(user *model.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
