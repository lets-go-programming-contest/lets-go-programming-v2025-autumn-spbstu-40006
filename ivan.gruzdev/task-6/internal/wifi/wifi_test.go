package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

type stubHandle struct {
	res []*wifi.Interface
	err error
}

func (s stubHandle) Interfaces() ([]*wifi.Interface, error) {
	return s.res, s.err
}

func TestGetNames_Success(t *testing.T) {
	stub := stubHandle{
		res: []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
		},
	}

	service := New(stub)
	names, err := service.GetNames()

	assert.NoError(t, err)
	assert.Equal(t, []string{"wlan0", "wlan1"}, names)
}

func TestGetNames_Error(t *testing.T) {
	stub := stubHandle{
		res: nil,
		err: errors.New("fail"),
	}

	service := New(stub)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
}

func TestGetAddresses_Success(t *testing.T) {
	stub := stubHandle{
		res: []*wifi.Interface{
			{HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22}},
			{HardwareAddr: net.HardwareAddr{0x33, 0x44, 0x55}},
		},
		err: nil,
	}

	service := New(stub)
	addrs, err := service.GetAddresses()

	assert.NoError(t, err)
	assert.Equal(t, []net.HardwareAddr{
		{0x00, 0x11, 0x22},
		{0x33, 0x44, 0x55},
	}, addrs)
}

func TestGetAddresses_Error(t *testing.T) {
	stub := stubHandle{
		res: nil,
		err: errors.New("fail"),
	}

	service := New(stub)
	addrs, err := service.GetAddresses()

	assert.Error(t, err)
	assert.Nil(t, addrs)
}
