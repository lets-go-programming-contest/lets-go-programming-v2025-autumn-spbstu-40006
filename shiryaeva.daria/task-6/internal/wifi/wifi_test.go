package wifi_test

import (
    "io"
    "net"
    "testing"

    "github.com/mdlayher/wifi"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    wifisvc "github.com/Dora-shi/task-6/internal/wifi"
)

type wiFiHandleMock struct {
    mock.Mock
}

func (m *wiFiHandleMock) Interfaces() ([]*wifi.Interface, error) {
    args := m.Called()

    var out []*wifi.Interface
    v0 := args.Get(0)
    if v0 != nil {
        if typed, ok := v0.([]*wifi.Interface); ok {
            out = typed
        }
    }

    var err error
    v1 := args.Get(1)
    if v1 != nil {
        if typedErr, ok := v1.(error); ok {
            err = typedErr
        }
    }

    return out, err
}

func TestWiFiService_New(t *testing.T) {
    t.Parallel()

    m := &wiFiHandleMock{}
    service := wifisvc.New(m)

    assert.NotNil(t, service, "service should not be nil")
    assert.Equal(t, m, service.WiFi, "WiFi handle should be set correctly")
}

func TestWiFiService_GetAddresses_OK(t *testing.T) {
    t.Parallel()

    m := &wiFiHandleMock{}

    hwAddr1 := net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
    hwAddr2 := net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}

    ifaces := []*wifi.Interface{
        {
            Name:         "wlan0",
            HardwareAddr: hwAddr1,
        },
        {
            Name:         "wlan1",
            HardwareAddr: hwAddr2,
        },
    }

    m.On("Interfaces").Return(ifaces, nil).Once()

    service := wifisvc.New(m)

    got, err := service.GetAddresses()
    require.NoError(t, err, "GetAddresses should not fail")

    assert.Len(t, got, 2, "should return 2 addresses")
    assert.Equal(t, hwAddr1, got[0], "first address should match")
    assert.Equal(t, hwAddr2, got[1], "second address should match")

    m.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_EmptyList(t *testing.T) {
    t.Parallel()

    m := &wiFiHandleMock{}
    m.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()

    service := wifisvc.New(m)

    got, err := service.GetAddresses()
    require.NoError(t, err, "GetAddresses should not fail with empty list")

    assert.Empty(t, got, "should return empty slice")

    m.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_NilHardwareAddr(t *testing.T) {
    t.Parallel()

    m := &wiFiHandleMock{}

    hwAddr := net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}

    ifaces := []*wifi.Interface{
        {
            Name:         "wlan0",
            HardwareAddr: nil,
        },
        {
            Name:         "wlan1",
            HardwareAddr: hwAddr,
        },
    }

    m.On("Interfaces").Return(ifaces, nil).Once()

    service := wifisvc.New(m)

    got, err := service.GetAddresses()
    require.NoError(t, err, "GetAddresses should not fail with nil addresses")

    assert.Len(t, got, 2, "should return 2 addresses")
    assert.Nil(t, got[0], "first address should be nil")
    assert.Equal(t, hwAddr, got[1], "second address should match")

    m.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_Error(t *testing.T) {
    t.Parallel()

    m := &wiFiHandleMock{}
    m.On("Interfaces").Return(nil, io.EOF).Once()

    service := wifisvc.New(m)

    _, err := service.GetAddresses()
    require.Error(t, err, "should return error")
    assert.Contains(t, err.Error(), "getting interfaces:", "error should be wrapped")

    m.AssertExpectations(t)
}

func TestWiFiService_GetNames_OK(t *testing.T) {
    t.Parallel()

    m := &wiFiHandleMock{}


    ifaces := []*wifi.Interface{
        {Name: "wlan0", HardwareAddr: nil},
        {Name: "wlan1", HardwareAddr: nil},
        {Name: "eth0", HardwareAddr: nil},
    }

    m.On("Interfaces").Return(ifaces, nil).Once()

    service := wifisvc.New(m)

    got, err := service.GetNames()
    require.NoError(t, err, "GetNames should not fail")

    assert.Equal(t, []string{"wlan0", "wlan1", "eth0"}, got, "names should match")

    m.AssertExpectations(t)
}

func TestWiFiService_GetNames_EmptyList(t *testing.T) {
    t.Parallel()

    m := &wiFiHandleMock{}
    m.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()

    service := wifisvc.New(m)

    got, err := service.GetNames()
    require.NoError(t, err, "GetNames should not fail with empty list")

    assert.Empty(t, got, "should return empty slice")

    m.AssertExpectations(t)
}


func TestWiFiService_GetNames_Error(t *testing.T) {
    t.Parallel()

    m := &wiFiHandleMock{}
    m.On("Interfaces").Return(nil, io.EOF).Once()

    service := wifisvc.New(m)

    _, err := service.GetNames()
    require.Error(t, err, "should return error")
    assert.Contains(t, err.Error(), "getting interfaces:", "error should be wrapped")

    m.AssertExpectations(t)
}

func TestWiFiService_AddressesAndNamesConsistency(t *testing.T) {
    t.Parallel()

    m := &wiFiHandleMock{}

    hwAddr1 := net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
    hwAddr2 := net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}

    ifaces := []*wifi.Interface{
        {
            Name:         "wlan0",
            HardwareAddr: hwAddr1,
        },
        {
            Name:         "wlan1",
            HardwareAddr: hwAddr2,
        },
    }

    m.On("Interfaces").Return(ifaces, nil).Twice()

    service := wifisvc.New(m)

    addresses, err := service.GetAddresses()
    require.NoError(t, err, "GetAddresses should not fail")
    assert.Len(t, addresses, 2, "should return 2 addresses")

    names, err := service.GetNames()
    require.NoError(t, err, "GetNames should not fail")
    assert.Len(t, names, 2, "should return 2 names")

    assert.Equal(t, "wlan0", names[0], "first name should be wlan0")
    assert.Equal(t, "wlan1", names[1], "second name should be wlan1")

    m.AssertExpectations(t)
}
