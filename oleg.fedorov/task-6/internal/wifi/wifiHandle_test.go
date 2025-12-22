package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiInterface struct {
	mock.Mock
}

func (m *MockWiFiInterface) GetInterfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		err := args.Error(1)
		if err != nil {
			return []*wifi.Interface{}, fmt.Errorf("mock error: %w", err)
		}

		return []*wifi.Interface{}, nil
	}

	interfaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, fmt.Errorf("type assertion failed")
	}

	err := args.Error(1)
	if err != nil {
		return interfaces, fmt.Errorf("mock error: %w", err)
	}

	return interfaces, nil
}

type MockWiFiHandler struct {
	mock.Mock
}

func (m *MockWiFiHandler) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		err := args.Error(1)
		if err != nil {
			return []*wifi.Interface{}, fmt.Errorf("mock error: %w", err)
		}

		return []*wifi.Interface{}, nil
	}

	interfaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, fmt.Errorf("type assertion failed")
	}

	err := args.Error(1)
	if err != nil {
		return interfaces, fmt.Errorf("mock error: %w", err)
	}

	return interfaces, nil
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
