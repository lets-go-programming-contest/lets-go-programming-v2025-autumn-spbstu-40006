package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	service "github.com/Mishaa105/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock error: %w", err)
	}

	if val := args.Get(0); val != nil {
		if interfaces, ok := val.([]*wifi.Interface); ok {
			return interfaces, nil
		}
	}

	return nil, nil
}

var errInterfacesFailed = errors.New("failed to get interfaces")

func TestWiFiService_New(t *testing.T) {
	t.Parallel()

	mockHandle := &MockWiFiHandle{}
	svc := service.New(mockHandle)

	assert.NotNil(t, svc)
	assert.Same(t, mockHandle, svc.WiFi)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success with multiple interfaces", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		mac1, _ := net.ParseMAC("aa:bb:cc:00:00:01")
		mac2, _ := net.ParseMAC("aa:bb:cc:00:00:02")

		interfaces := []*wifi.Interface{
			{HardwareAddr: mac1},
			{HardwareAddr: mac2},
		}

		mockHandle.On("Interfaces").Return(interfaces, nil).Once()

		addresses, err := svc.GetAddresses()

		require.NoError(t, err)
		assert.Equal(t, []net.HardwareAddr{mac1, mac2}, addresses)
		mockHandle.AssertExpectations(t)
	})

	t.Run("success with empty result", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()

		addresses, err := svc.GetAddresses()

		require.NoError(t, err)
		assert.Empty(t, addresses)
		mockHandle.AssertExpectations(t)
	})

	t.Run("error from Interfaces", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errInterfacesFailed).Once()

		addresses, err := svc.GetAddresses()

		require.Error(t, err)
		assert.Nil(t, addresses)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success with multiple names", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		interfaces := []*wifi.Interface{
			{Name: "wlp3s0"},
			{Name: "wlan0"},
		}

		mockHandle.On("Interfaces").Return(interfaces, nil).Once()

		names, err := svc.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"wlp3s0", "wlan0"}, names)
		mockHandle.AssertExpectations(t)
	})

	t.Run("success with empty result", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()

		names, err := svc.GetNames()

		require.NoError(t, err)
		assert.Empty(t, names)
		mockHandle.AssertExpectations(t)
	})

	t.Run("error from Interfaces", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errInterfacesFailed).Once()

		names, err := svc.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockHandle.AssertExpectations(t)
	})
}
