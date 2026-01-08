package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/dizey5k/task-6/internal/wifi"
	md_wifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errSystem = errors.New("system error")

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mockInterfaces := []*md_wifi.Interface{{Name: "wlan0"}}

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)

		service := wifi.New(mockHandle)
		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0"}, names)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, errSystem)

		service := wifi.New(mockHandle)
		_, err := service.GetNames()

		require.Error(t, err)
	})
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mac, _ := net.ParseMAC("00:00:5e:00:53:01")
		mockInterfaces := []*md_wifi.Interface{{HardwareAddr: mac}}

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)

		service := wifi.New(mockHandle)
		addrs, err := service.GetAddresses()

		require.NoError(t, err)
		assert.Len(t, addrs, 1)
		assert.Equal(t, mac, addrs[0])
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, errSystem)

		service := wifi.New(mockHandle)
		_, err := service.GetAddresses()

		require.Error(t, err)
	})
}
