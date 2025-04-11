// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	dto "github.com/developeerz/restorio-auth/internal/handler/user/dto"
	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Login provides a mock function with given fields: req
func (_m *Service) Login(req *dto.LoginRequest) (int, *dto.JwtAccessResponse, string, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 int
	var r1 *dto.JwtAccessResponse
	var r2 string
	var r3 error
	if rf, ok := ret.Get(0).(func(*dto.LoginRequest) (int, *dto.JwtAccessResponse, string, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*dto.LoginRequest) int); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(*dto.LoginRequest) *dto.JwtAccessResponse); ok {
		r1 = rf(req)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dto.JwtAccessResponse)
		}
	}

	if rf, ok := ret.Get(2).(func(*dto.LoginRequest) string); ok {
		r2 = rf(req)
	} else {
		r2 = ret.Get(2).(string)
	}

	if rf, ok := ret.Get(3).(func(*dto.LoginRequest) error); ok {
		r3 = rf(req)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// SignUp provides a mock function with given fields: req
func (_m *Service) SignUp(req *dto.SignUpRequest) (int, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for SignUp")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(*dto.SignUpRequest) (int, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*dto.SignUpRequest) int); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(*dto.SignUpRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Verify provides a mock function with given fields: req
func (_m *Service) Verify(req *dto.VerificationRequest) (int, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for Verify")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(*dto.VerificationRequest) (int, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*dto.VerificationRequest) int); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(*dto.VerificationRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
