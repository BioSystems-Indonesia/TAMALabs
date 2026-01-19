//go:build !windows
// +build !windows

package util

import "fmt"

// stubs for non-windows platforms
func windowsGetSystemUUID() (string, error) {
	return "", fmt.Errorf("windows system UUID not available on this platform")
}

func windowsGetBiosSerial() (string, error) {
	return "", fmt.Errorf("windows bios serial not available on this platform")
}
