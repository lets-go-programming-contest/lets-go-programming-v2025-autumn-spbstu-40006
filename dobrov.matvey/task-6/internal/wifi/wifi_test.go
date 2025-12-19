package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	myWiFi "github.com/netwite/task-6/internal/wifi"
	"github.com/stretchr/testify/suite"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

type mockWiFiHandle struct {
	interfaces []*wifi.Interface
	err        error
}

func (m *mockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	return m.interfaces, m.err
}

type WiFiServiceTestSuite struct {
	suite.Suite
}

func (s *WiFiServiceTestSuite) TestNew() {
	mockWiFi := &mockWiFiHandle{}
	service := myWiFi.New(mockWiFi)
	s.Equal(mockWiFi, service.WiFi)
}

func (s *WiFiServiceTestSuite) TestGetAddresses_Success() {
	expectedAddrs := []net.HardwareAddr{
		s.parseMAC("00:11:22:33:44:55"),
		s.parseMAC("aa:bb:cc:dd:ee:ff"),
	}

	mockInterfaces := []*wifi.Interface{
		{
			Index:        1,
			Name:         "wlan0",
			HardwareAddr: expectedAddrs[0],
		},
		{
			Index:        2,
			Name:         "eth0",
			HardwareAddr: expectedAddrs[1],
		},
	}

	mockWiFi := &mockWiFiHandle{
		interfaces: mockInterfaces,
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetAddresses()

	s.Require().NoError(err)
	s.Equal(expectedAddrs, result)
}

func (s *WiFiServiceTestSuite) TestGetAddresses_EmptyResult() {
	mockWiFi := &mockWiFiHandle{
		interfaces: []*wifi.Interface{},
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetAddresses()

	s.Require().NoError(err)
	s.Empty(result)
}

func (s *WiFiServiceTestSuite) TestGetAddresses_InterfacesError() {
	testError := errors.New("interfaces error") //nolint:err113
	mockWiFi := &mockWiFiHandle{
		err: testError,
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetAddresses()

	s.Require().Error(err)
	s.Require().ErrorContains(err, "getting interfaces")
	s.Contains(err.Error(), testError.Error())
	s.Nil(result)
}

func (s *WiFiServiceTestSuite) TestGetAddresses_WithNilHardwareAddr() {
	mockInterfaces := []*wifi.Interface{
		{
			Index:        1,
			Name:         "wlan0",
			HardwareAddr: nil,
		},
		{
			Index:        2,
			Name:         "eth0",
			HardwareAddr: s.parseMAC("aa:bb:cc:dd:ee:ff"),
		},
	}

	expectedAddrs := []net.HardwareAddr{
		nil,
		s.parseMAC("aa:bb:cc:dd:ee:ff"),
	}

	mockWiFi := &mockWiFiHandle{
		interfaces: mockInterfaces,
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetAddresses()

	s.Require().NoError(err)
	s.Equal(expectedAddrs, result)
}

func (s *WiFiServiceTestSuite) TestGetNames_Success() {
	expectedNames := []string{"wlan0", "eth0", "wlan1"}

	mockInterfaces := []*wifi.Interface{
		{
			Index:        1,
			Name:         "wlan0",
			HardwareAddr: s.parseMAC("00:11:22:33:44:55"),
		},
		{
			Index:        2,
			Name:         "eth0",
			HardwareAddr: s.parseMAC("aa:bb:cc:dd:ee:ff"),
		},
		{
			Index:        3,
			Name:         "wlan1",
			HardwareAddr: s.parseMAC("11:22:33:44:55:66"),
		},
	}

	mockWiFi := &mockWiFiHandle{
		interfaces: mockInterfaces,
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetNames()

	s.Require().NoError(err)
	s.Equal(expectedNames, result)
}

func (s *WiFiServiceTestSuite) TestGetNames_EmptyResult() {
	mockWiFi := &mockWiFiHandle{
		interfaces: []*wifi.Interface{},
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetNames()

	s.Require().NoError(err)
	s.Empty(result)
}

func (s *WiFiServiceTestSuite) TestGetNames_InterfacesError() {
	testError := errors.New("interfaces error") //nolint:err113
	mockWiFi := &mockWiFiHandle{
		err: testError,
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetNames()

	s.Require().Error(err)
	s.Require().ErrorContains(err, "getting interfaces")
	s.Contains(err.Error(), testError.Error())
	s.Nil(result)
}

func (s *WiFiServiceTestSuite) TestGetNames_WithEmptyName() {
	mockInterfaces := []*wifi.Interface{
		{
			Index:        1,
			Name:         "",
			HardwareAddr: s.parseMAC("00:11:22:33:44:55"),
		},
		{
			Index:        2,
			Name:         "eth0",
			HardwareAddr: s.parseMAC("aa:bb:cc:dd:ee:ff"),
		},
	}

	expectedNames := []string{"", "eth0"}

	mockWiFi := &mockWiFiHandle{
		interfaces: mockInterfaces,
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetNames()

	s.Require().NoError(err)
	s.Equal(expectedNames, result)
}

func (s *WiFiServiceTestSuite) TestGetNames_SingleInterface() {
	mockInterfaces := []*wifi.Interface{
		{
			Index:        1,
			Name:         "wlan0",
			HardwareAddr: s.parseMAC("00:11:22:33:44:55"),
		},
	}

	expectedNames := []string{"wlan0"}

	mockWiFi := &mockWiFiHandle{
		interfaces: mockInterfaces,
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetNames()

	s.Require().NoError(err)
	s.Equal(expectedNames, result)
}

func (s *WiFiServiceTestSuite) TestGetNames_SpecialCharacterNames() {
	mockInterfaces := []*wifi.Interface{
		{
			Index:        1,
			Name:         "wlan-0",
			HardwareAddr: s.parseMAC("00:11:22:33:44:55"),
		},
		{
			Index:        2,
			Name:         "eth_1",
			HardwareAddr: s.parseMAC("aa:bb:cc:dd:ee:ff"),
		},
	}

	expectedNames := []string{"wlan-0", "eth_1"}

	mockWiFi := &mockWiFiHandle{
		interfaces: mockInterfaces,
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetNames()

	s.Require().NoError(err)
	s.Equal(expectedNames, result)
}

func (s *WiFiServiceTestSuite) TestGetAddresses_SingleInterface() {
	expectedAddr := s.parseMAC("00:11:22:33:44:55")

	mockInterfaces := []*wifi.Interface{
		{
			Index:        1,
			Name:         "wlan0",
			HardwareAddr: expectedAddr,
		},
	}

	mockWiFi := &mockWiFiHandle{
		interfaces: mockInterfaces,
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetAddresses()

	s.Require().NoError(err)
	s.Equal([]net.HardwareAddr{expectedAddr}, result)
}

func (s *WiFiServiceTestSuite) TestGetAddresses_ZeroMACAddress() {
	zeroMAC := s.parseMAC("00:00:00:00:00:00")

	mockInterfaces := []*wifi.Interface{
		{
			Index:        1,
			Name:         "wlan0",
			HardwareAddr: zeroMAC,
		},
	}

	mockWiFi := &mockWiFiHandle{
		interfaces: mockInterfaces,
	}
	service := myWiFi.New(mockWiFi)

	result, err := service.GetAddresses()

	s.Require().NoError(err)
	s.Equal([]net.HardwareAddr{zeroMAC}, result)
}

func (s *WiFiServiceTestSuite) parseMAC(macStr string) net.HardwareAddr {
	hwAddr, err := net.ParseMAC(macStr)
	s.Require().NoError(err, "failed to parse MAC address: %s", macStr)

	return hwAddr
}

func TestWiFiServiceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WiFiServiceTestSuite))
}
