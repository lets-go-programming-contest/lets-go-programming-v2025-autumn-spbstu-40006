package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

// WiFiHandle is a mock implementation of WiFiHandle.
type WiFiHandle struct {
	mock.Mock
}

// Interfaces provides a mock function with no fields.
func (_m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	var r0 []*wifi.Interface

	if rf, ok := ret.Get(0).(func() []*wifi.Interface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*wifi.Interface)
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

// Close provides a mock function with no fields.
func (_m *WiFiHandle) Close() error {
	ret := _m.Called()

	var r0 error

	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWiFiHandle creates a new instance of WiFiHandle.
// It registers a testing interface on the mock and a cleanup function to assert expectations.
// The first argument is typically a *testing.T value.
func NewWiFiHandle(t interface {
	mock.TestingT
	Cleanup(func())
}) *WiFiHandle {
	m := &WiFiHandle{}

	m.Mock.Test(t)

	t.Cleanup(func() {
		m.AssertExpectations(t)
	})

	return m
}

// Ensure interface compliance.
var _ fmt.Stringer = (*WiFiHandle)(nil)
