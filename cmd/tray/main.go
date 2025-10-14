package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"syscall"

	"github.com/energye/systray"
)

//go:embed trayicon.ico
var trayicon []byte

type ServiceStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func main() {
	systray.Run(onReady, nil)
}

func onReady() {
	systray.SetTitle("TAMALabs")
	systray.SetTooltip("TAMALabs Service Controller")
	systray.SetIcon(trayicon)

	mStatus := systray.AddMenuItem("Service Status: Checking...", "Current service status")
	mStatus.Disable()
	systray.AddSeparator()

	mStart := systray.AddMenuItem("Start Service", "Start the TAMALabs service")
	mStop := systray.AddMenuItem("Stop Service", "Stop the TAMALabs service")
	mRestart := systray.AddMenuItem("Restart Service", "Restart the TAMALabs service")
	systray.AddSeparator()

	mOpen := systray.AddMenuItem("Open Web Interface", "Open TAMALabs in browser")
	systray.AddSeparator()
	mExit := systray.AddMenuItem("Exit Tray", "Exit this tray app")

	// Update status periodically
	go func() {
		for {
			status := getServiceStatus()
			mStatus.SetTitle(fmt.Sprintf("Service Status: %s", status))

			// Enable/disable buttons based on status
			switch strings.ToLower(status) {
			case "running":
				mStart.Disable()
				mStop.Enable()
				mRestart.Enable()
			case "stopped":
				mStart.Enable()
				mStop.Disable()
				mRestart.Disable()
			case "not installed":
				mStart.Disable()
				mStop.Disable()
				mRestart.Disable()
				mStatus.SetTitle("Service Status: Not Installed (Run installer)")
			default:
				// Unknown status - enable all for manual control
				mStart.Enable()
				mStop.Enable()
				mRestart.Enable()
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// Event handlers
	mStart.Click(func() {
		go func() {
			mStart.SetTitle("Starting...")
			mStart.Disable()

			success := controlService("start")

			// Reset button title
			mStart.SetTitle("Start Service")

			if success {
				showNotification("TAMALabs service started successfully")
			} else {
				showError("Failed to start service")
			}

			// Force status update
			time.Sleep(1 * time.Second)
		}()
	})

	mStop.Click(func() {
		go func() {
			mStop.SetTitle("Stopping...")
			mStop.Disable()

			success := controlService("stop")

			// Reset button title
			mStop.SetTitle("Stop Service")

			if success {
				showNotification("TAMALabs service stopped successfully")
			} else {
				showError("Failed to stop service")
			}

			// Force status update
			time.Sleep(1 * time.Second)
		}()
	})

	mRestart.Click(func() {
		go func() {
			mRestart.SetTitle("Restarting...")
			mRestart.Disable()

			success := controlService("restart")

			// Reset button title
			mRestart.SetTitle("Restart Service")

			if success {
				showNotification("TAMALabs service restarted successfully")
			} else {
				showError("Failed to restart service")
			}

			// Force status update
			time.Sleep(1 * time.Second)
		}()
	})

	mOpen.Click(func() {
		openbrowser("http://127.0.0.1:8322")
	})

	mExit.Click(func() {
		systray.Quit()
	})
}

func controlService(action string) bool {
	exePath, err := os.Executable()
	if err != nil {
		showError(fmt.Sprintf("Failed to get executable path: %v", err))
		return false
	}

	helperPath := filepath.Join(filepath.Dir(exePath), "service-helper.exe")

	// Check if helper exists
	if _, err := os.Stat(helperPath); os.IsNotExist(err) {
		showError("Service helper not found. Please reinstall TAMALabs.")
		return false
	}

	cmd := exec.Command(helperPath, action, "TAMALabs")
	output, err := cmd.CombinedOutput()

	if err != nil {
		showError(fmt.Sprintf("Failed to %s service: %v\nOutput: %s", action, err, string(output)))
		return false
	}

	return true
}

func getServiceStatus() string {
	// First try to check if web interface is responding (primary check)
	// Check root path instead of /ping since /ping might not exist
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://127.0.0.1:8322/")
	if err == nil {
		defer resp.Body.Close()
		// Any response code means service is running (200, 404, etc.)
		return "Running"
	}

	// Fallback to direct NSSM check using correct path
	exePath, err := os.Executable()
	if err != nil {
		return "Unknown Path"
	}

	nssmPath := filepath.Join(filepath.Dir(exePath), "nssm.exe")

	// Check if NSSM exists
	if _, err := os.Stat(nssmPath); os.IsNotExist(err) {
		// Try system PATH as fallback
		nssmPath = "nssm"
	}

       cmd := exec.Command(nssmPath, "status", "TAMALabs")
       // Hide window on Windows
       if runtime.GOOS == "windows" {
	       cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
       }
       output, err := cmd.CombinedOutput()

	if err != nil {
		// If NSSM fails, service might not be installed yet
		return "Not Installed"
	}

	status := strings.TrimSpace(string(output))
	if strings.Contains(status, "SERVICE_RUNNING") {
		return "Running"
	} else if strings.Contains(status, "SERVICE_STOPPED") {
		return "Stopped"
	} else if strings.Contains(status, "SERVICE_PAUSED") {
		return "Paused"
	} else if strings.Contains(status, "service does not exist") {
		return "Not Installed"
	}

	return "Unknown"
}

func showNotification(message string) {
	// Simple notification via systray tooltip
	systray.SetTooltip(fmt.Sprintf("TAMALabs - %s", message))

	// Reset tooltip after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		systray.SetTooltip("TAMALabs Service Controller")
	}()
}

func showError(message string) {
	fmt.Printf("Error: %s\n", message)
	systray.SetTooltip(fmt.Sprintf("TAMALabs - Error: %s", message))

	// Reset tooltip after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		systray.SetTooltip("TAMALabs Service Controller")
	}()
}

func openbrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}
