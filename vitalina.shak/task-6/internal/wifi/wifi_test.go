package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	internalwifi "github.com/vitsh1/task-6/internal/wifi"
)

func TestWiFi_New(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	service := internalwifi.New(mockWiFi)

	require.Equal(t, mockWiFi, service.WiFi)
}

func TestWiFi_GetAddresses_OK(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	addr1, err := net.ParseMAC("00:11:22:33:44:55")
	require.NoError(t, err)
	addr2, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	require.NoError(t, err)

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{
		{Name: "wlan0", HardwareAddr: addr1},
		{Name: "wlan1", HardwareAddr: addr2},
	}, nil).Once()

	service := internalwifi.New(mockWiFi)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{addr1, addr2}, addrs)
}

func TestWiFi_GetAddresses_Empty(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()

	service := internalwifi.New(mockWiFi)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Empty(t, addrs)
}

func TestWiFi_GetAddresses_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	expected := errors.New("error")
	mockWiFi.On("Interfaces").Return(([]*wifi.Interface)(nil), expected).Once()

	service := internalwifi.New(mockWiFi)

	addrs, err := service.GetAddresses()
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Contains(t, err.Error(), "getting interfaces")
	require.Nil(t, addrs)
}

func TestWiFiService_GetNames_OK(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	addr1, err := net.ParseMAC("00:11:22:33:44:55")
	require.NoError(t, err)
	addr2, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	require.NoError(t, err)

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{
		{Name: "wlan0", HardwareAddr: addr1},
		{Name: "wlan1", HardwareAddr: addr2},
	}, nil).Once()

	service := internalwifi.New(mockWiFi)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wlan0", "wlan1"}, names)
}

func TestWiFiService_GetNames_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	expected := errors.New("error")
	mockWiFi.On("Interfaces").Return(([]*wifi.Interface)(nil), expected).Once()

	service := internalwifi.New(mockWiFi)

	names, err := service.GetNames()
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Contains(t, err.Error(), "getting interfaces")
	require.Nil(t, names)
}

func TestWiFi_GetNames_EmptyList(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()

	service := internalwifi.New(mockWiFi)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Empty(t, names)
}
