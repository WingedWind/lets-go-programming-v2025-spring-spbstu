// Code generated by mockery v2.53.2. DO NOT EDIT.

package wifi_test

import (
	wifi "github.com/mdlayher/wifi"
	mock "github.com/stretchr/testify/mock"
)

// WiFi is an autogenerated mock type for the WiFi type
type WiFi struct {
	mock.Mock
}

// Interfaces provides a mock function with no fields
func (_m *WiFi) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Interfaces")
	}

	var r0 []*wifi.Interface
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*wifi.Interface, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*wifi.Interface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*wifi.Interface)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewWiFi creates a new instance of WiFi. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWiFi(t interface {
	mock.TestingT
	Cleanup(func())
}) *WiFi {
	mock := &WiFi{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
