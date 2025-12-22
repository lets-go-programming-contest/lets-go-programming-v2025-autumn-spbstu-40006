package wifi_test

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiInterface struct {
	mock.Mock
}

func (m *MockWiFiInterface) GetInterfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return []*wifi.Interface{}, args.Error(1)
	}

	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

type MockWiFiHandler struct {
	mock.Mock
}

func (m *MockWiFiHandler) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return []*wifi.Interface{}, args.Error(1)
	}

	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

func (m *MockWiFiHandler) SetupSuccess(interfaces []*wifi.Interface) {
	m.On("Interfaces").Return(interfaces, nil)
}

func (m *MockWiFiHandler) SetupFailure(err error) {
	m.On("Interfaces").Return([]*wifi.Interface{}, err)
}

func (m *MockWiFiHandler) SetupEmpty() {
	m.On("Interfaces").Return([]*wifi.Interface{}, nil)
}
