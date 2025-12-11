package wifi_test

import (
	"errors"
	"net"
	"testing"

	svc "github.com/identicalaffiliation/task-6/internal/wifi"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errWiFi = errors.New("wifi error")

func TestNew(t *testing.T) {
	t.Parallel()

	mockWiFiHandle := &MockWiFiHandle{}

	service := svc.New(mockWiFiHandle)
	assert.NotNil(t, service)
	assert.Equal(t, mockWiFiHandle, service.WiFi)
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("completed", func(t *testing.T) {
		t.Parallel()

		mockWiFiHandle := &MockWiFiHandle{}
		service := svc.New(mockWiFiHandle)
		addr1, err := net.ParseMAC("11:22:33:44:55:66")
		require.NoError(t, err)
		addr2, err := net.ParseMAC("66:55:44:33:22:11")
		require.NoError(t, err)

		interfaces := []*wifi.Interface{
			{HardwareAddr: addr1},
			{HardwareAddr: addr2},
		}

		mockWiFiHandle.On("Interfaces").Return(interfaces, nil).Once()

		addresses, err := service.GetAddresses()
		require.NoError(t, err)
		assert.Equal(t, []net.HardwareAddr{addr1, addr2}, addresses)
		mockWiFiHandle.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)

		service := svc.New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errWiFi).Once()

		addrs, err := service.GetAddresses()
		require.Error(t, err)
		assert.Nil(t, addrs)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockHandle.AssertExpectations(t)
	})
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	t.Run("completed", func(t *testing.T) {
		t.Parallel()

		mockWiFiHandle := &MockWiFiHandle{}
		service := svc.New(mockWiFiHandle)

		interfaces := []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "eth0"},
		}

		mockWiFiHandle.On("Interfaces").Return(interfaces, nil).Once()

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "eth0"}, names)
		mockWiFiHandle.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		t.Parallel()

		mockWiFiHandle := &MockWiFiHandle{}
		service := svc.New(mockWiFiHandle)

		mockWiFiHandle.On("Interfaces").Return(nil, errWiFi).Once()

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockWiFiHandle.AssertExpectations(t)
	})
}
