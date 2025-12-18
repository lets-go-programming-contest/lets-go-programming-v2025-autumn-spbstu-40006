package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifipkg "github.com/zaouiomar/task-6/internal/wifi"
)

var (
	errSystem     = errors.New("system error")
	errPermission = errors.New("permission error")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	mac1 := net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	mac2 := net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	mac3 := net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	mac4 := net.HardwareAddr{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}

	tests := []struct {
		name        string
		setupMock   func(*MockWiFiHandle)
		expectAddrs []net.HardwareAddr
		expectErr   bool
		errContains string
	}{
		{
			name: "success two interfaces",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{
					{HardwareAddr: mac1},
					{HardwareAddr: mac2},
				}, nil).Once()
			},
			expectAddrs: []net.HardwareAddr{mac1, mac2},
		},
		{
			name: "success four interfaces",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{
					{HardwareAddr: mac1},
					{HardwareAddr: mac2},
					{HardwareAddr: mac3},
					{HardwareAddr: mac4},
				}, nil).Once()
			},
			expectAddrs: []net.HardwareAddr{mac1, mac2, mac3, mac4},
		},
		{
			name: "empty interfaces",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()
			},
			expectAddrs: []net.HardwareAddr{},
		},
		{
			name: "interfaces error",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errSystem).Once()
			},
			expectErr:   true,
			errContains: "getting interfaces:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := new(MockWiFiHandle)
			tc.setupMock(mockHandle)

			svc := wifipkg.New(mockHandle)
			addrs, err := svc.GetAddresses()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
				assert.Nil(t, addrs)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectAddrs, addrs)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	mac1 := net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	mac2 := net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}

	tests := []struct {
		name        string
		setupMock   func(*MockWiFiHandle)
		expectNames []string
		expectErr   bool
		errContains string
	}{
		{
			name: "success two names",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{
					{Name: "wlan0", HardwareAddr: mac1},
					{Name: "eth0", HardwareAddr: mac2},
				}, nil).Once()
			},
			expectNames: []string{"wlan0", "eth0"},
		},
		{
			name: "empty list",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()
			},
			expectNames: []string{},
		},
		{
			name: "interface with empty name",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{
					{Name: ""},
					{Name: "wlan1"},
				}, nil).Once()
			},
			expectNames: []string{"", "wlan1"},
		},
		{
			name: "error from Interfaces",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errPermission).Once()
			},
			expectErr:   true,
			errContains: "getting interfaces:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := new(MockWiFiHandle)
			tc.setupMock(mockHandle)

			svc := wifipkg.New(mockHandle)
			names, err := svc.GetNames()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectNames, names)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}
