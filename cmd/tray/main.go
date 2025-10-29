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
	"syscall"
	"time"
	"io"
	"log"

	"github.com/energye/systray"
)

//go:embed trayicon.ico
var trayicon []byte

type ServiceStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func main() {
	// initialize logging to logs/tray.log next to the executable
	lf, err := setupLogging()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
	} else {
		defer lf.Close()
	}

	systray.Run(onReady, nil)
}

// setupLogging creates a logs directory next to the executable and
// configures the standard logger to write to logs/tray.log (and stdout).
func setupLogging() (*os.File, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	base := filepath.Dir(exePath)
	logsDir := filepath.Join(base, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, err
	}
	logPath := filepath.Join(logsDir, "tray.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Logger initialized, writing to %s", logPath)
	return f, nil
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
	status := strings.TrimSpace(string(output))

	// Parse NSSM output regardless of the command exit code. NSSM may return
	// non-zero for certain states but still output a useful status string.
	if strings.Contains(strings.ToLower(status), "service_running") || strings.Contains(strings.ToLower(status), "running") {
		return "Running"
	}
	if strings.Contains(strings.ToLower(status), "service_stopped") || strings.Contains(strings.ToLower(status), "stopped") {
		return "Stopped"
	}
	if strings.Contains(strings.ToLower(status), "service_paused") || strings.Contains(strings.ToLower(status), "paused") {
		return "Paused"
	}
	if strings.Contains(strings.ToLower(status), "service does not exist") || strings.Contains(strings.ToLower(status), "not installed") {
		return "Not Installed"
	}

	// If we have an error but couldn't parse a known status, log the output for debugging and return Unknown
	if err != nil {
		log.Printf("nssm status error: %v, output: %s\n", err, status)

		// Fallback: try using sc.exe which is available on all Windows systems
		if runtime.GOOS == "windows" {
			scCmd := exec.Command("sc", "query", "TAMALabs")
			if runtime.GOOS == "windows" {
				scCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			}
			scOut, scErr := scCmd.CombinedOutput()
			scStatus := strings.ToLower(strings.TrimSpace(string(scOut)))
			if scErr == nil || len(scStatus) > 0 {
				if strings.Contains(scStatus, "running") || strings.Contains(scStatus, "state") && strings.Contains(scStatus, "running") {
					return "Running"
				}
				if strings.Contains(scStatus, "stopped") || strings.Contains(scStatus, "state") && strings.Contains(scStatus, "stopped") {
					return "Stopped"
				}
				if strings.Contains(scStatus, "does not exist") || strings.Contains(scStatus, "not found") {
					return "Not Installed"
				}
				// If SC returned something we don't understand, fall through to Unknown
				log.Printf("sc query output: %s, err: %v\n", scStatus, scErr)
			}
		}

		return "Unknown"
	}

	// Fallback to Unknown if nothing matched
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
	log.Printf("Error: %s", message)
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
