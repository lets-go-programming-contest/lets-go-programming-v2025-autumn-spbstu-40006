package wifi_test

import "github.com/mdlayher/wifi"

type MockWiFiHandle struct {
	InterfacesFunc func() ([]*wifi.Interface, error)
	Calls          int
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	m.Calls++
	if m.InterfacesFunc != nil {
		return m.InterfacesFunc()
	}

	return nil, nil
}
