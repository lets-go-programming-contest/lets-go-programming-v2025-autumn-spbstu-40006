package wifi

import "net"

// Interface is a minimal stand-in for github.com/mdlayher/wifi.Interface.
// It includes only the fields used by the training exercise.
type Interface struct {
	Name         string
	HardwareAddr net.HardwareAddr
}
