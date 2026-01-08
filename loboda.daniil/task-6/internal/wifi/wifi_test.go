package wifi

import (
	"errors"
	"net"
	"reflect"
	"testing"

	"github.com/mdlayher/wifi"
)

func TestWiFiService_GetAddresses_Success(t *testing.T) {
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

	svc := New(m)
	addrs, err := svc.GetAddresses()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []net.HardwareAddr{hw1, hw2}
	if !reflect.DeepEqual(addrs, want) {
		t.Fatalf("addresses mismatch: got %v want %v", addrs, want)
	}
	if m.Calls != 1 {
		t.Fatalf("Interfaces should be called once, got %d", m.Calls)
	}
}

func TestWiFiService_GetAddresses_Error(t *testing.T) {
	boom := errors.New("boom")
	m := &MockWiFiHandle{InterfacesFunc: func() ([]*wifi.Interface, error) { return nil, boom }}

	svc := New(m)
	_, err := svc.GetAddresses()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}

func TestWiFiService_GetNames_Success(t *testing.T) {
	m := &MockWiFiHandle{
		InterfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{Name: "wlan0"},
				{Name: "wlan1"},
			}, nil
		},
	}

	svc := New(m)
	names, err := svc.GetNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"wlan0", "wlan1"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("names mismatch: got %v want %v", names, want)
	}
}

func TestWiFiService_GetNames_Error(t *testing.T) {
	boom := errors.New("boom")
	m := &MockWiFiHandle{InterfacesFunc: func() ([]*wifi.Interface, error) { return nil, boom }}

	svc := New(m)
	_, err := svc.GetNames()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}
