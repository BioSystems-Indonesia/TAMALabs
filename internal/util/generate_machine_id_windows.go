//go:build windows
// +build windows

package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// windowsGetSystemUUID uses "wmic csproduct get uuid" (best-effort)
func windowsGetSystemUUID() (string, error) {
	cmd := exec.Command("wmic", "csproduct", "get", "uuid")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		// fallback to wmic csproduct get uuid /format:list
		cmd2 := exec.Command("wmic", "csproduct", "get", "uuid", "/format:list")
		cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out2, _ := cmd2.Output()
		if len(out2) == 0 {
			return "", err
		}
		out = out2
	}
	// output contains header and value; get last non-empty line
	lines := bytes.Split(out, []byte("\n"))
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(string(lines[i]))
		if line == "" || strings.Contains(strings.ToLower(line), "uuid") {
			continue
		}
		return line, nil
	}
	return "", fmt.Errorf("uuid not found in wmic output")
}

func windowsGetBiosSerial() (string, error) {
	cmd := exec.Command("wmic", "bios", "get", "serialnumber")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := bytes.Split(out, []byte("\n"))
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(string(lines[i]))
		if line == "" || strings.Contains(strings.ToLower(line), "serialnumber") {
			continue
		}
		return line, nil
	}
	return "", fmt.Errorf("bios serial not found")
}
