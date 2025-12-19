package wifi

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type WiFiHandle struct {
	mock.Mock
}

func (_m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	var r0 []*wifi.Interface

	if rf, ok := ret.Get(0).(func() []*wifi.Interface); ok {
		r0 = rf()
	} else if v, ok := ret.Get(0).([]*wifi.Interface); ok {
		r0 = v
	}

	var r1 error

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		if err := ret.Error(1); err != nil {
			r1 = fmt.Errorf("mock error: %w", err)
		}
	}

	return r0, r1
}

func (_m *WiFiHandle) Close() error {
	ret := _m.Called()

	var r0 error

	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		if err := ret.Error(0); err != nil {
			r0 = fmt.Errorf("mock error: %w", err)
		}
	}

	return r0
}

func NewWiFiHandle(t interface {
	mock.TestingT
	Cleanup(f func())
}) *WiFiHandle {
	m := &WiFiHandle{}
	m.Mock.Test(t)

	t.Cleanup(func() {
		m.AssertExpectations(t)
	})

	return m
}
