// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	models "github.com/developeerz/restorio-auth/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CheckVerificationCode provides a mock function with given fields: userCode
func (_m *Repository) CheckVerificationCode(userCode *models.UserCode) (int64, error) {
	ret := _m.Called(userCode)

	if len(ret) == 0 {
		panic("no return value specified for CheckVerificationCode")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.UserCode) (int64, error)); ok {
		return rf(userCode)
	}
	if rf, ok := ret.Get(0).(func(*models.UserCode) int64); ok {
		r0 = rf(userCode)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(*models.UserCode) error); ok {
		r1 = rf(userCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: _a0
func (_m *Repository) CreateUser(_a0 *models.User) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUserAuth provides a mock function with given fields: userAuth
func (_m *Repository) CreateUserAuth(userAuth *models.UserAuth) error {
	ret := _m.Called(userAuth)

	if len(ret) == 0 {
		panic("no return value specified for CreateUserAuth")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.UserAuth) error); ok {
		r0 = rf(userAuth)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateVerificationCode provides a mock function with given fields: userCode
func (_m *Repository) CreateVerificationCode(userCode *models.UserCode) error {
	ret := _m.Called(userCode)

	if len(ret) == 0 {
		panic("no return value specified for CreateVerificationCode")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.UserCode) error); ok {
		r0 = rf(userCode)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteVerificationCode provides a mock function with given fields: userCode
func (_m *Repository) DeleteVerificationCode(userCode *models.UserCode) error {
	ret := _m.Called(userCode)

	if len(ret) == 0 {
		panic("no return value specified for DeleteVerificationCode")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.UserCode) error); ok {
		r0 = rf(userCode)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByTelegram provides a mock function with given fields: telegram
func (_m *Repository) FindByTelegram(telegram string) (*models.User, error) {
	ret := _m.Called(telegram)

	if len(ret) == 0 {
		panic("no return value specified for FindByTelegram")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.User, error)); ok {
		return rf(telegram)
	}
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(telegram)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(telegram)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByTelegramWithAuths provides a mock function with given fields: telegram
func (_m *Repository) FindByTelegramWithAuths(telegram string) (*models.UserWithAuths, error) {
	ret := _m.Called(telegram)

	if len(ret) == 0 {
		panic("no return value specified for FindByTelegramWithAuths")
	}

	var r0 *models.UserWithAuths
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.UserWithAuths, error)); ok {
		return rf(telegram)
	}
	if rf, ok := ret.Get(0).(func(string) *models.UserWithAuths); ok {
		r0 = rf(telegram)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserWithAuths)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(telegram)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveUser provides a mock function with given fields: _a0
func (_m *Repository) SaveUser(_a0 *models.User) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SaveUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetUserAuth provides a mock function with given fields: userAuth
func (_m *Repository) SetUserAuth(userAuth *models.UserAuth) error {
	ret := _m.Called(userAuth)

	if len(ret) == 0 {
		panic("no return value specified for SetUserAuth")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.UserAuth) error); ok {
		r0 = rf(userAuth)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
