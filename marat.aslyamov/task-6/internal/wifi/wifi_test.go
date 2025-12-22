package wifi_test

import (
	"errors"
	"net"
	"testing"

	service "github.com/IvanIgnashin7D/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (_m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	ifaces := ret.Get(0)
	if ifaces == nil {
		return nil, ret.Error(1)
	}

	return ifaces.([]*wifi.Interface), ret.Error(1)
}

var (
	errInterface  = errors.New("interface error")
	errPermission = errors.New("permission denied")
)

func TestWiFiService_New(t *testing.T) {
	t.Parallel()

	mockHandle := &MockWiFiHandle{}
	svc := service.New(mockHandle)

	assert.NotNil(t, svc)
	assert.Same(t, mockHandle, svc.WiFi)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	mac1, _ := net.ParseMAC("aa:bb:cc:00:00:01")
	mac2, _ := net.ParseMAC("aa:bb:cc:00:00:02")
	mac3, _ := net.ParseMAC("aa:bb:cc:00:00:03")

	testCases := []struct {
		name        string
		mockSetup   func(*MockWiFiHandle)
		wantAddrs   []net.HardwareAddr
		wantError   bool
		errorSubstr string
	}{
		{
			name: "multiple interfaces",
			mockSetup: func(m *MockWiFiHandle) {
				ifaces := []*wifi.Interface{
					{Name: "wlan0", HardwareAddr: mac1},
					{Name: "wlan1", HardwareAddr: mac2},
					{Name: "wlan2", HardwareAddr: mac3},
				}
				m.On("Interfaces").Return(ifaces, nil).Once()
			},
			wantAddrs: []net.HardwareAddr{mac1, mac2, mac3},
		},
		{
			name: "single interface",
			mockSetup: func(m *MockWiFiHandle) {
				ifaces := []*wifi.Interface{
					{Name: "wlan0", HardwareAddr: mac1},
				}
				m.On("Interfaces").Return(ifaces, nil).Once()
			},
			wantAddrs: []net.HardwareAddr{mac1},
		},
		{
			name: "empty result",
			mockSetup: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()
			},
			wantAddrs: []net.HardwareAddr{},
		},
		{
			name: "interface with nil hardware address",
			mockSetup: func(m *MockWiFiHandle) {
				ifaces := []*wifi.Interface{
					{Name: "wlan0", HardwareAddr: mac1},
					{Name: "wlan1", HardwareAddr: nil},
					{Name: "wlan2", HardwareAddr: mac2},
				}
				m.On("Interfaces").Return(ifaces, nil).Once()
			},
			wantAddrs: []net.HardwareAddr{mac1, nil, mac2},
		},
		{
			name: "error from interface",
			mockSetup: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface(nil), errInterface).Once()
			},
			wantError:   true,
			errorSubstr: "getting interfaces",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &MockWiFiHandle{}
			tc.mockSetup(mockHandle)
			svc := service.New(mockHandle)

			addrs, err := svc.GetAddresses()

			if tc.wantError {
				require.Error(t, err)
				assert.Nil(t, addrs)

				if tc.errorSubstr != "" {
					assert.Contains(t, err.Error(), tc.errorSubstr)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantAddrs, addrs)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	mac, _ := net.ParseMAC("aa:bb:cc:00:00:01")

	testCases := []struct {
		name        string
		mockSetup   func(*MockWiFiHandle)
		wantNames   []string
		wantError   bool
		errorSubstr string
	}{
		{
			name: "multiple interface names",
			mockSetup: func(m *MockWiFiHandle) {
				ifaces := []*wifi.Interface{
					{Name: "wlp3s0", HardwareAddr: mac},
					{Name: "wlan0", HardwareAddr: mac},
					{Name: "eth1", HardwareAddr: mac},
				}
				m.On("Interfaces").Return(ifaces, nil).Once()
			},
			wantNames: []string{"wlp3s0", "wlan0", "eth1"},
		},
		{
			name: "single interface name",
			mockSetup: func(m *MockWiFiHandle) {
				ifaces := []*wifi.Interface{
					{Name: "wlan0", HardwareAddr: mac},
				}
				m.On("Interfaces").Return(ifaces, nil).Once()
			},
			wantNames: []string{"wlan0"},
		},
		{
			name: "empty name allowed",
			mockSetup: func(m *MockWiFiHandle) {
				ifaces := []*wifi.Interface{
					{Name: "", HardwareAddr: mac},
					{Name: "wlan0", HardwareAddr: mac},
				}
				m.On("Interfaces").Return(ifaces, nil).Once()
			},
			wantNames: []string{"", "wlan0"},
		},
		{
			name: "empty result",
			mockSetup: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()
			},
			wantNames: []string{},
		},
		{
			name: "error from interface",
			mockSetup: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface(nil), errPermission).Once()
			},
			wantError:   true,
			errorSubstr: "getting interfaces",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &MockWiFiHandle{}
			tc.mockSetup(mockHandle)
			svc := service.New(mockHandle)

			names, err := svc.GetNames()

			if tc.wantError {
				require.Error(t, err)
				assert.Nil(t, names)

				if tc.errorSubstr != "" {
					assert.Contains(t, err.Error(), tc.errorSubstr)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantNames, names)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}
