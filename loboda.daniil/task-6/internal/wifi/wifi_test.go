package wifi_test

import (
	"errors"
	"net"
	"reflect"
	"testing"

	wifisvc "loboda.daniil/task-6/internal/wifi"

	"github.com/mdlayher/wifi"
)

var errBoom = errors.New("boom")

func TestWiFiService_GetAddresses_Success(t *testing.T) {
	t.Parallel()

	hw1, _ := net.ParseMAC("aa:bb:cc:dd:ee:01")
	hw2, _ := net.ParseMAC("aa:bb:cc:dd:ee:02")

	m := &MockWiFiHandle{
		InterfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{Name: "wlan0", HardwareAddr: hw1},
				{Name: "wlan1", HardwareAddr: hw2},
			}, nil
		},
	}

	svc := wifisvc.New(m)

	got, err := svc.GetAddresses()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []net.HardwareAddr{hw1, hw2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestWiFiService_GetAddresses_Error(t *testing.T) {
	t.Parallel()

	m := &MockWiFiHandle{
		InterfacesFunc: func() ([]*wifi.Interface, error) {
			return nil, errBoom
		},
	}

	svc := wifisvc.New(m)
	_, err := svc.GetAddresses()

	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, errBoom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}

func TestWiFiService_GetNames_Success(t *testing.T) {
	t.Parallel()

	m := &MockWiFiHandle{
		InterfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{Name: "wlan0"},
				{Name: "wlan1"},
			}, nil
		},
	}

	svc := wifisvc.New(m)

	got, err := svc.GetNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"wlan0", "wlan1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestWiFiService_GetNames_Error(t *testing.T) {
	t.Parallel()

	m := &MockWiFiHandle{
		InterfacesFunc: func() ([]*wifi.Interface, error) {
			return nil, errBoom
		},
	}

	svc := wifisvc.New(m)
	_, err := svc.GetNames()

	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, errBoom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}
