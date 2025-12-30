package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mywifi "github.com/MoneyprogerISG/task-6/internal/wifi"
)

type stubHandle struct {
	res []*wifi.Interface
	err error
}

func (s stubHandle) Interfaces() ([]*wifi.Interface, error) {
	return s.res, s.err
}

var errFail = errors.New("fail")

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	stub := stubHandle{
		res: []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
		},
	}

	service := mywifi.New(stub)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"wlan0", "wlan1"}, names)
}

func TestGetNames_Error(t *testing.T) {
	t.Parallel()

	stub := stubHandle{
		res: nil,
		err: errFail,
	}

	service := mywifi.New(stub)
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
}

func TestGetAddresses_Success(t *testing.T) {
	t.Parallel()

	stub := stubHandle{
		res: []*wifi.Interface{
			{HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22}},
			{HardwareAddr: net.HardwareAddr{0x33, 0x44, 0x55}},
		},
		err: nil,
	}

	service := mywifi.New(stub)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	assert.Equal(t, []net.HardwareAddr{
		{0x00, 0x11, 0x22},
		{0x33, 0x44, 0x55},
	}, addrs)
}

func TestGetAddresses_Error(t *testing.T) {
	t.Parallel()

	stub := stubHandle{
		res: nil,
		err: errFail,
	}

	service := mywifi.New(stub)
	addrs, err := service.GetAddresses()

	require.Error(t, err)
	assert.Nil(t, addrs)
}
