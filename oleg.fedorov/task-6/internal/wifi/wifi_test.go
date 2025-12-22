package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	wifi_pkg "github.com/dizey5k/task-6/internal/wifi"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errMockInterfaces = errors.New("mock interfaces error")
)

type TestWiFiHandler struct {
	mock.Mock
}

func (t *TestWiFiHandler) Interfaces() ([]*wifi.Interface, error) {
	args := t.Called()

	if args.Get(0) == nil {
		err := args.Error(1)
		if err != nil {
			return []*wifi.Interface{}, fmt.Errorf("mock error: %w", err)
		}

		return []*wifi.Interface{}, nil
	}

	interfaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, fmt.Errorf("type assertion failed")
	}

	err := args.Error(1)
	if err != nil {
		return interfaces, fmt.Errorf("mock error: %w", err)
	}

	return interfaces, nil
}

func createTestInterface(name string, addr net.HardwareAddr) *wifi.Interface {
	return &wifi.Interface{
		Name:         name,
		HardwareAddr: addr,
	}
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMock     func(*TestWiFiHandler)
		expectedAddrs []net.HardwareAddr
		expectError   bool
	}{
		{
			name: "success getting address",
			setupMock: func(m *TestWiFiHandler) {
				interfaces := []*wifi.Interface{
					createTestInterface("wlan0", net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}),
					createTestInterface("wlan1", net.HardwareAddr{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}),
				}
				m.On("Interfaces").Return(interfaces, nil)
			},
			expectedAddrs: []net.HardwareAddr{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
			},
			expectError: false,
		},
		{
			name: "empty list interfaces",
			setupMock: func(m *TestWiFiHandler) {
				m.On("Interfaces").Return([]*wifi.Interface{}, nil)
			},
			expectedAddrs: []net.HardwareAddr{},
			expectError:   false,
		},
		{
			name: "err while getting interfaces",
			setupMock: func(m *TestWiFiHandler) {
				m.On("Interfaces").Return([]*wifi.Interface{}, errMockInterfaces)
			},
			expectedAddrs: nil,
			expectError:   true,
		},
		{
			name: "interface with zero address",
			setupMock: func(m *TestWiFiHandler) {
				interfaces := []*wifi.Interface{
					createTestInterface("wlan0", nil),
					createTestInterface("wlan1", net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}),
				}
				m.On("Interfaces").Return(interfaces, nil)
			},
			expectedAddrs: []net.HardwareAddr{
				nil,
				{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockHandler := &TestWiFiHandler{}
			tt.setupMock(mockHandler)

			service := wifi_pkg.New(mockHandler)
			addrs, err := service.GetAddresses()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedAddrs, addrs)
			}

			mockHandler.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMock     func(*TestWiFiHandler)
		expectedNames []string
		expectError   bool
	}{
		{
			name: "success getting name interfaces",
			setupMock: func(m *TestWiFiHandler) {
				interfaces := []*wifi.Interface{
					createTestInterface("wlan0", nil),
					createTestInterface("eth1", nil),
					createTestInterface("wifi2", nil),
				}
				m.On("Interfaces").Return(interfaces, nil)
			},
			expectedNames: []string{"wlan0", "eth1", "wifi2"},
			expectError:   false,
		},
		{
			name: "one interface",
			setupMock: func(m *TestWiFiHandler) {
				interfaces := []*wifi.Interface{
					createTestInterface("single", nil),
				}
				m.On("Interfaces").Return(interfaces, nil)
			},
			expectedNames: []string{"single"},
			expectError:   false,
		},
		{
			name: "err in Interfaces",
			setupMock: func(m *TestWiFiHandler) {
				m.On("Interfaces").Return([]*wifi.Interface{}, errMockInterfaces)
			},
			expectedNames: nil,
			expectError:   true,
		},
		{
			name: "interfaces with same name",
			setupMock: func(m *TestWiFiHandler) {
				interfaces := []*wifi.Interface{
					createTestInterface("wlan0", nil),
					createTestInterface("wlan0", nil),
					createTestInterface("eth1", nil),
				}
				m.On("Interfaces").Return(interfaces, nil)
			},
			expectedNames: []string{"wlan0", "wlan0", "eth1"},
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockHandler := &TestWiFiHandler{}
			tt.setupMock(mockHandler)

			service := wifi_pkg.New(mockHandler)
			names, err := service.GetNames()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedNames, names)
			}

			mockHandler.AssertExpectations(t)
		})
	}
}

func TestWiFiService_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("nil in return value", func(t *testing.T) {
		t.Parallel()

		mockHandler := &TestWiFiHandler{}
		mockHandler.On("Interfaces").Return(nil, nil)

		service := wifi_pkg.New(mockHandler)

		addrs, err := service.GetAddresses()
		require.NoError(t, err)
		assert.NotNil(t, addrs)
		assert.Empty(t, addrs)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.NotNil(t, names)
		assert.Empty(t, names)

		mockHandler.AssertExpectations(t)
	})

	t.Run("type assertion failed", func(t *testing.T) {
		t.Parallel()

		mockHandler := &TestWiFiHandler{}
		mockHandler.On("Interfaces").Return("not a slice", nil).Twice()

		service := wifi_pkg.New(mockHandler)

		addrs, err := service.GetAddresses()
		require.Error(t, err)
		assert.Nil(t, addrs)
		assert.Contains(t, err.Error(), "type assertion failed")

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "type assertion failed")

		mockHandler.AssertExpectations(t)
	})

	t.Run("service with zero handler", func(t *testing.T) {
		t.Parallel()

		service := wifi_pkg.WiFiService{}

		addrs, err := service.GetAddresses()
		require.Error(t, err)
		assert.Nil(t, addrs)

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
	})
}
