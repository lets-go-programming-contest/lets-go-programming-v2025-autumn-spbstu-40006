package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type WiFiHandle struct {
	mock.Mock
}

func NewWiFiHandle(t interface {
	mock.TestingT
	Cleanup(f func())
},
) *WiFiHandle {
	m := &WiFiHandle{}
	m.Mock.Test(t)

	t.Cleanup(func() {
		m.AssertExpectations(t)
	})

	return m
}

func (m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var ifaces []*wifi.Interface

	first := args.Get(0)

	if first != nil {
		casted, ok := first.([]*wifi.Interface)
		if ok {
			ifaces = casted
		}
	}

	err := args.Error(1)

	if err != nil {
		return ifaces, fmt.Errorf("mock error: %w", err)
	}

	return ifaces, nil
}
