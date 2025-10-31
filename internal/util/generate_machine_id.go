package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// GenerateMachineID returns deterministic machine id like "MACHINE-<first-20-hex-chars>"
// It tries multiple sources (BIOS/CSProduct UUID, MAC addresses, hostname, /etc/machine-id)
// and hashes them. Works on Windows, Linux, macOS (best-effort).
// Uses stable hardware-based identifiers first, fallback to less stable ones.
func GenerateMachineID() (string, error) {
	var parts []string

	// 1) MOST STABLE: Platform-specific hardware UUIDs (highest priority)
	if runtime.GOOS == "windows" {
		if uuid, err := windowsGetSystemUUID(); err == nil && uuid != "" {
			parts = append(parts, "uuid:"+uuid)
		}
		// BIOS serial is very stable
		if sn, err := windowsGetBiosSerial(); err == nil && sn != "" {
			parts = append(parts, "bios:"+sn)
		}
	} else {
		// Linux/macOS: try /etc/machine-id, /var/lib/dbus/machine-id, or ioreg on mac
		if mid, err := readMachineIDFile(); err == nil && mid != "" {
			parts = append(parts, "mid:"+mid)
		}
		if runtime.GOOS == "darwin" {
			if uuid, err := macGetHardwareUUID(); err == nil && uuid != "" {
				parts = append(parts, "uuid:"+uuid)
			}
		} else {
			// linux: try dmidecode or hostnamectl (best-effort)
			if uuid, err := linuxGetDmiUUID(); err == nil && uuid != "" {
				parts = append(parts, "uuid:"+uuid)
			}
		}
	}

	// 2) STABLE: Get only physical MAC addresses (sorted for consistency)
	if macs, err := getStableMACs(); err == nil && len(macs) > 0 {
		parts = append(parts, "macs:"+strings.Join(macs, ","))
	}

	// 3) FALLBACK: Only add hostname if no hardware identifiers found
	if len(parts) == 0 {
		if hn, err := os.Hostname(); err == nil && hn != "" {
			parts = append(parts, "host:"+hn)
		}
	}

	// If nothing found (very unlikely), return error
	if len(parts) == 0 {
		return "", fmt.Errorf("no machine identifiers available")
	}

	combined := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(combined))
	hexStr := hex.EncodeToString(hash[:])
	// return first 20 chars to keep it reasonably short but unique
	id := fmt.Sprintf("MACHINE-%s", strings.ToUpper(hexStr)[:20])
	return id, nil
}

// GenerateMachineIDDebug returns machine ID with debug information
func GenerateMachineIDDebug() (string, map[string]string, error) {
	debug := make(map[string]string)
	var parts []string

	// 1) MOST STABLE: Platform-specific hardware UUIDs (highest priority)
	if runtime.GOOS == "windows" {
		if uuid, err := windowsGetSystemUUID(); err == nil && uuid != "" {
			parts = append(parts, "uuid:"+uuid)
			debug["windows_uuid"] = uuid
		} else {
			debug["windows_uuid_error"] = err.Error()
		}

		if sn, err := windowsGetBiosSerial(); err == nil && sn != "" {
			parts = append(parts, "bios:"+sn)
			debug["bios_serial"] = sn
		} else {
			debug["bios_serial_error"] = err.Error()
		}
	} else {
		if mid, err := readMachineIDFile(); err == nil && mid != "" {
			parts = append(parts, "mid:"+mid)
			debug["machine_id"] = mid
		} else {
			debug["machine_id_error"] = err.Error()
		}

		if runtime.GOOS == "darwin" {
			if uuid, err := macGetHardwareUUID(); err == nil && uuid != "" {
				parts = append(parts, "uuid:"+uuid)
				debug["mac_uuid"] = uuid
			} else {
				debug["mac_uuid_error"] = err.Error()
			}
		} else {
			if uuid, err := linuxGetDmiUUID(); err == nil && uuid != "" {
				parts = append(parts, "uuid:"+uuid)
				debug["dmi_uuid"] = uuid
			} else {
				debug["dmi_uuid_error"] = err.Error()
			}
		}
	}

	// 2) STABLE: Get only physical MAC addresses (sorted for consistency)
	if macs, err := getStableMACs(); err == nil && len(macs) > 0 {
		parts = append(parts, "macs:"+strings.Join(macs, ","))
		debug["stable_macs"] = strings.Join(macs, ",")
	} else {
		debug["stable_macs_error"] = err.Error()
	}

	// 3) FALLBACK: Only add hostname if no hardware identifiers found
	if len(parts) == 0 {
		if hn, err := os.Hostname(); err == nil && hn != "" {
			parts = append(parts, "host:"+hn)
			debug["hostname"] = hn
		}
	}

	debug["combined_parts"] = strings.Join(parts, "|")

	// If nothing found (very unlikely), return error
	if len(parts) == 0 {
		return "", debug, fmt.Errorf("no machine identifiers available")
	}

	combined := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(combined))
	hexStr := hex.EncodeToString(hash[:])
	id := fmt.Sprintf("MACHINE-%s", strings.ToUpper(hexStr)[:20])

	debug["final_id"] = id
	return id, debug, nil
}

/* ---------- helper functions ---------- */

// getStableMACs returns stable physical MAC addresses (sorted for consistency)
func getStableMACs() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var macs []string
	for _, itf := range ifaces {
		// Skip loopback, virtual, and temporary interfaces
		if itf.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Skip virtual interfaces (common patterns)
		name := strings.ToLower(itf.Name)
		if strings.Contains(name, "virtual") ||
			strings.Contains(name, "vmware") ||
			strings.Contains(name, "vbox") ||
			strings.Contains(name, "docker") ||
			strings.Contains(name, "bridge") ||
			strings.HasPrefix(name, "veth") ||
			strings.HasPrefix(name, "tap") ||
			strings.HasPrefix(name, "tun") {
			continue
		}

		m := itf.HardwareAddr.String()
		if m == "" || m == "00:00:00:00:00:00" {
			continue
		}

		// Normalize MAC address
		clean := strings.ToUpper(strings.ReplaceAll(m, ":", ""))
		clean = strings.ReplaceAll(clean, "-", "")
		if clean != "" && len(clean) == 12 {
			macs = append(macs, clean)
		}
	}

	// Sort MACs for consistency
	if len(macs) > 1 {
		for i := 0; i < len(macs)-1; i++ {
			for j := i + 1; j < len(macs); j++ {
				if macs[i] > macs[j] {
					macs[i], macs[j] = macs[j], macs[i]
				}
			}
		}
	}

	// Return only first 2 MACs to avoid too much variability
	if len(macs) > 2 {
		macs = macs[:2]
	}

	return macs, nil
}

// windowsGetSystemUUID uses "wmic csproduct get uuid" (best-effort)
func readMachineIDFile() (string, error) {
	paths := []string{
		"/etc/machine-id",
		"/var/lib/dbus/machine-id",
	}
	for _, p := range paths {
		if data, err := os.ReadFile(p); err == nil {
			s := strings.TrimSpace(string(data))
			if s != "" {
				return s, nil
			}
		}
	}
	return "", fmt.Errorf("no machine-id file found")
}

// windowsGetSystemUUID uses "wmic csproduct get uuid" (best-effort)
func windowsGetSystemUUID() (string, error) {
	cmd := exec.Command("wmic", "csproduct", "get", "uuid")
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	out, err := cmd.Output()
	if err != nil {
		// fallback to wmic csproduct get uuid /format:list
		cmd2 := exec.Command("wmic", "csproduct", "get", "uuid", "/format:list")
		if runtime.GOOS == "windows" {
			cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		}
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
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
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

// linuxGetDmiUUID tries "cat /sys/class/dmi/id/product_uuid" or hostnamectl
func linuxGetDmiUUID() (string, error) {
	paths := []string{
		"/sys/class/dmi/id/product_uuid",
	}
	for _, p := range paths {
		if data, err := os.ReadFile(p); err == nil {
			s := strings.TrimSpace(string(data))
			if s != "" {
				return s, nil
			}
		}
	}
	// fallback to "hostnamectl"
	if out, err := exec.Command("hostnamectl", "status", "--pretty").Output(); err == nil {
		// less reliable: try "hostnamectl"
		txt := strings.TrimSpace(string(out))
		if txt != "" {
			return txt, nil
		}
	}
	return "", fmt.Errorf("dmi uuid not found")
}

// macGetHardwareUUID uses "ioreg -rd1 -c IOPlatformExpertDevice" parsing (best-effort)
func macGetHardwareUUID() (string, error) {
	out, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
	if err != nil {
		return "", err
	}
	// look for "IOPlatformUUID" = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
	o := string(out)
	for _, line := range strings.Split(o, "\n") {
		if strings.Contains(line, "IOPlatformUUID") {
			parts := strings.Split(line, "=")
			if len(parts) >= 2 {
				u := strings.TrimSpace(parts[1])
				u = strings.Trim(u, "\" ")
				return u, nil
			}
		}
	}
	return "", fmt.Errorf("IOPlatformUUID not found")
}
