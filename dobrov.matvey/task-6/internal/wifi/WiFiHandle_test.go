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

	if len(ret) == 0 {
		panic("no return value specified for Interfaces")
	}

	// First return value.
	if rf, ok := ret.Get(0).(func() ([]*wifi.Interface, error)); ok {
		return rf()
	}

	var r0 []*wifi.Interface
	v0 := ret.Get(0)

	if rf, ok := v0.(func() []*wifi.Interface); ok {
		r0 = rf()
	} else if v0 != nil {
		var ok bool
		r0, ok = v0.([]*wifi.Interface)
		if !ok {
			panic("invalid type for Interfaces return value")
		}
	}

	// Second return value.
	if rf, ok := ret.Get(1).(func() error); ok {
		return r0, rf()
	}

	err := ret.Error(1)
	if err != nil {
		err = fmt.Errorf("ret.Error(1): %w", err)
	}

	return r0, err
}

// NewWiFiHandle creates a new instance of WiFiHandle.
// It registers a testing interface on the mock and a cleanup function to assert expectations.
// The first argument is typically a *testing.T value.
func NewWiFiHandle(t interface {
	mock.TestingT
	Cleanup(fn func())
},
) *WiFiHandle {
	m := &WiFiHandle{}
	m.Mock.Test(t)

	t.Cleanup(func() { m.AssertExpectations(t) })

	return m
}
