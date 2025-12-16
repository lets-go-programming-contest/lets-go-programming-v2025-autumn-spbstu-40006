package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	wifimocks "github.com/MargotBush/task-6/internal/wifi/mocks"
)

func TestWiFiService_GetAddresses_OK(t *testing.T) {
	m := wifimocks.NewWiFiHandle(t)

	if1 := &wifi.Interface{Name: "wlan0", HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}}
	if2 := &wifi.Interface{Name: "wlan1", HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}}

	m.EXPECT().Interfaces().Return([]*wifi.Interface{if1, if2}, nil)

	service := New(m)
	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{if1.HardwareAddr, if2.HardwareAddr}, addrs)

	m.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_Error(t *testing.T) {
	m := wifimocks.NewWiFiHandle(t)

	m.EXPECT().Interfaces().Return(nil, errors.New("no permission"))

	service := New(m)
	_, err := service.GetAddresses()
	require.Error(t, err) // also covers wrapping "getting interfaces: %w"

	m.AssertExpectations(t)
}

func TestWiFiService_GetNames_OK(t *testing.T) {
	m := wifimocks.NewWiFiHandle(t)

	if1 := &wifi.Interface{Name: "wlan0"}
	if2 := &wifi.Interface{Name: "wlan1"}

	m.EXPECT().Interfaces().Return([]*wifi.Interface{if1, if2}, nil)

	service := New(m)
	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wlan0", "wlan1"}, names)

	m.AssertExpectations(t)
}

func TestWiFiService_GetNames_Error(t *testing.T) {
	m := wifimocks.NewWiFiHandle(t)

	m.EXPECT().Interfaces().Return(nil, errors.New("wifi failure"))

	service := New(m)
	_, err := service.GetNames()
	require.Error(t, err)

	m.AssertExpectations(t)
}
