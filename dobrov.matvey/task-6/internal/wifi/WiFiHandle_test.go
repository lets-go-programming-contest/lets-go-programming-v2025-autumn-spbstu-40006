package wifi_test

import (
	"fmt"

	wifi "github.com/mdlayher/wifi"
	mock "github.com/stretchr/testify/mock"
)

// WiFiHandle is a mock implementation of WiFiHandle.
type WiFiHandle struct {
	mock.Mock
}

// Interfaces provides a mock function with no fields.
func (_m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
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
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*wifi.Interface)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		err := ret.Error(1)
		if err != nil {
			r1 = fmt.Errorf("ret.Error(1): %w", err)
		}
	}

	return r0, r1
}

// NewWiFiHandle creates a new instance of WiFiHandle.
// It registers a testing interface on the mock and a cleanup function to assert expectations.
// The first argument is typically a *testing.T value.
func NewWiFiHandle(t interface {
	mock.TestingT
	Cleanup(fn func())
}) *WiFiHandle {
	mock := &WiFiHandle{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
