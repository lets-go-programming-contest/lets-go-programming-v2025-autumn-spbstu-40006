package wifi_test

import (
	"errors"
	"net"
	"testing"

	svc "github.com/filon6/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

var errExpected = errors.New("expected error")

func TestNewWiFiServiceStoresHandle(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	service := svc.New(mockWiFi)

	require.Equal(t, mockWiFi, service.WiFi, "wifi handle should be stored inside service")
	mockWiFi.AssertExpectations(t)
}

func TestGetAddresses_OK(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	hwAddr1, err := net.ParseMAC("00:11:22:33:44:55")
	require.NoError(t, err)
	hwAddr2, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	require.NoError(t, err)

	interfaces := []*wifi.Interface{
		{Name: "wifi-home", HardwareAddr: hwAddr1},
		{Name: "wifi-guest", HardwareAddr: hwAddr2},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Once()

	service := svc.New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Len(t, addrs, 2)
	require.Equal(t, hwAddr1, addrs[0])
	require.Equal(t, hwAddr2, addrs[1])

	mockWiFi.AssertExpectations(t)
}

func TestGetAddresses_EmptyList(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	interfaces := []*wifi.Interface{}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Once()

	service := svc.New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Empty(t, addrs, "no interfaces should give empty slice")

	mockWiFi.AssertExpectations(t)
}

func TestGetAddresses_ErrorFromHandle(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, errExpected).Once()

	service := svc.New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.Error(t, err)
	require.Nil(t, addrs)
	require.Contains(t, err.Error(), "getting interfaces")

	mockWiFi.AssertExpectations(t)
}

func TestGetNames_OK(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	hwAddr, err := net.ParseMAC("00:11:22:33:44:55")
	require.NoError(t, err)

	interfaces := []*wifi.Interface{
		{Name: "wifi-home", HardwareAddr: hwAddr},
		{Name: "wifi-guest", HardwareAddr: nil},
		{Name: "ethernet-main", HardwareAddr: nil},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Once()

	service := svc.New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	require.Len(t, names, 3)
	require.Equal(t, []string{"wifi-home", "wifi-guest", "ethernet-main"}, names)

	mockWiFi.AssertExpectations(t)
}

func TestGetNames_EmptyList(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	interfaces := []*wifi.Interface{}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Once()

	service := svc.New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	require.Empty(t, names)

	mockWiFi.AssertExpectations(t)
}

func TestGetNames_ErrorFromHandle(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, errExpected).Once()

	service := svc.New(mockWiFi)
	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "getting interfaces")

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_AddressesAndNamesFromSameInput(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	hwAddr1, err := net.ParseMAC("00:11:22:33:44:55")
	require.NoError(t, err)
	hwAddr2, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	require.NoError(t, err)
	hwAddr3, err := net.ParseMAC("11:22:33:44:55:66")
	require.NoError(t, err)

	interfaces := []*wifi.Interface{
		{Name: "wifi-home", HardwareAddr: hwAddr1},
		{Name: "wifi-guest", HardwareAddr: hwAddr2},
		{Name: "ethernet-main", HardwareAddr: hwAddr3},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Twice()

	service := svc.New(mockWiFi)

	addresses, err := service.GetAddresses()
	require.NoError(t, err, "GetAddresses should not fail")
	require.Len(t, addresses, 3, "expected 3 addresses")

	names, err := service.GetNames()
	require.NoError(t, err, "GetNames should not fail")
	require.Len(t, names, 3, "expected 3 names")

	require.Equal(t, "wifi-home", names[0], "first interface name mismatch")
	require.Equal(t, "wifi-guest", names[1], "second interface name mismatch")
	require.Equal(t, "ethernet-main", names[2], "third interface name mismatch")

	mockWiFi.AssertExpectations(t)
}

func TestGetAddresses_NilHardwareAddr(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	hwAddr, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	require.NoError(t, err)

	interfaces := []*wifi.Interface{
		{
			Name:         "wifi-home",
			HardwareAddr: nil,
		},
		{Name: "wifi-guest", HardwareAddr: hwAddr},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Once()

	service := svc.New(mockWiFi)
	addresses, err := service.GetAddresses()
	require.NoError(t, err)
	require.Len(t, addresses, 2, "expected 2 addresses")
	require.Nil(t, addresses[0], "first address should be nil")
	require.Equal(t, hwAddr.String(), addresses[1].String(), "second address mismatch")

	mockWiFi.AssertExpectations(t)
}
