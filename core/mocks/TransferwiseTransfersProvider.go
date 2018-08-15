// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import core "github.com/endiangroup/transferwiser/core"
import mock "github.com/stretchr/testify/mock"

// TransferwiseTransfersProvider is an autogenerated mock type for the TransferwiseTransfersProvider type
type TransferwiseTransfersProvider struct {
	mock.Mock
}

// IsAuthenticated provides a mock function with given fields:
func (_m *TransferwiseTransfersProvider) IsAuthenticated() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Transfers provides a mock function with given fields:
func (_m *TransferwiseTransfersProvider) Transfers() ([]*core.Transfer, error) {
	ret := _m.Called()

	var r0 []*core.Transfer
	if rf, ok := ret.Get(0).(func() []*core.Transfer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Transfer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UseAuthentication provides a mock function with given fields: _a0
func (_m *TransferwiseTransfersProvider) UseAuthentication(_a0 *core.RefreshTokenData) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*core.RefreshTokenData) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
