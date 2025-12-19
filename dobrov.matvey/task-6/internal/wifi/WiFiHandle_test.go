package wifi_test

import (
	"fmt"

	myWiFi "github.com/HorekProgrammer/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type WiFiHandleMock struct {
	mock.Mock
}

func (m *WiFiHandleMock) Interfaces() ([]*wifi.Interface, error) {
	ret := m.Called()

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

func (m *WiFiHandleMock) Close() error {
	ret := m.Called()

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

func NewWiFiHandleMock(t interface {
	mock.TestingT
	Cleanup(f func())
},
) *WiFiHandleMock {
	m := &WiFiHandleMock{}
	m.Mock.Test(t)

	t.Cleanup(func() {
		m.AssertExpectations(t)
	})

	return m
}

var _ myWiFi.WiFiHandle = (*WiFiHandleMock)(nil)
