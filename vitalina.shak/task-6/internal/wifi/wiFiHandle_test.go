package wifi_test

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type WiFiHandle struct {
	mock.Mock
}

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

func (m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var ifaces []*wifi.Interface
	if first := args.Get(0); first != nil {
		casted, ok := first.([]*wifi.Interface)
		if ok {
			ifaces = casted
		}
	}

	return ifaces, args.Error(1)
}
