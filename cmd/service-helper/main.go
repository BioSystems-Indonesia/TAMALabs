package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	// Create a log file for debugging
	logFile, _ := os.OpenFile(filepath.Join(filepath.Dir(getExePath()), "service-helper.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if logFile != nil {
		defer logFile.Close()
		logFile.WriteString(fmt.Sprintf("\n=== %s ===\n", time.Now().Format("2006-01-02 15:04:05")))
		logFile.WriteString(fmt.Sprintf("Args: %v\n", os.Args))
	}

	if len(os.Args) < 3 {
		msg := "Usage: service-helper.exe <action> <service-name>\nActions: start, stop, restart, status"
		fmt.Println(msg)
		if logFile != nil {
			logFile.WriteString("ERROR: " + msg + "\n")
		}
		os.Exit(1)
	}

	action := os.Args[1]
	serviceName := os.Args[2]

	if logFile != nil {
		logFile.WriteString(fmt.Sprintf("Action: %s, Service: %s\n", action, serviceName))
	}

	switch action {
	case "start":
		startService(serviceName, logFile)
	case "stop":
		stopService(serviceName, logFile)
	case "restart":
		restartService(serviceName, logFile)
	case "status":
		getServiceStatus(serviceName, logFile)
	default:
		msg := fmt.Sprintf("Unknown action: %s", action)
		fmt.Println(msg)
		if logFile != nil {
			logFile.WriteString("ERROR: " + msg + "\n")
		}
		os.Exit(1)
	}
}

func getExePath() string {
	exePath, err := os.Executable()
	if err != nil {
		return "."
	}
	return exePath
}

func getNSSMPath() string {
	// Try to find NSSM in the same directory as this executable
	exePath := getExePath()
	nssmPath := filepath.Join(filepath.Dir(exePath), "nssm.exe")
	if _, err := os.Stat(nssmPath); err == nil {
		return nssmPath
	}

	// Fallback to system PATH
	return "nssm"
}

func runNSSMCommand(args []string, logFile *os.File) (string, error) {
	nssmPath := getNSSMPath()
	cmd := exec.Command(nssmPath, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	
	if logFile != nil {
		logFile.WriteString(fmt.Sprintf("Executing: %s %v\n", nssmPath, args))
	}
	
	output, err := cmd.CombinedOutput()
	
	if logFile != nil {
		logFile.WriteString(fmt.Sprintf("Output: %s\n", string(output)))
		if err != nil {
			logFile.WriteString(fmt.Sprintf("Error: %v\n", err))
		}
	}
	
	return string(output), err
}

func startService(serviceName string, logFile *os.File) {
	output, err := runNSSMCommand([]string{"start", serviceName}, logFile)

	if err != nil {
		msg := fmt.Sprintf("Error starting service: %v\nOutput: %s", err, output)
		fmt.Println(msg)
		if logFile != nil {
			logFile.WriteString("RESULT: FAILED - " + msg + "\n")
		}
		os.Exit(1)
	}

	msg := fmt.Sprintf("Service %s started successfully", serviceName)
	fmt.Println(msg)
	if logFile != nil {
		logFile.WriteString("RESULT: SUCCESS - " + msg + "\n")
	}
}

func stopService(serviceName string, logFile *os.File) {
	output, err := runNSSMCommand([]string{"stop", serviceName}, logFile)

	if err != nil {
		msg := fmt.Sprintf("Error stopping service: %v\nOutput: %s", err, output)
		fmt.Println(msg)
		if logFile != nil {
			logFile.WriteString("RESULT: FAILED - " + msg + "\n")
		}
		os.Exit(1)
	}

	msg := fmt.Sprintf("Service %s stopped successfully", serviceName)
	fmt.Println(msg)
	if logFile != nil {
		logFile.WriteString("RESULT: SUCCESS - " + msg + "\n")
	}
}

func restartService(serviceName string, logFile *os.File) {
	output, err := runNSSMCommand([]string{"restart", serviceName}, logFile)

	if err != nil {
		msg := fmt.Sprintf("Error restarting service: %v\nOutput: %s", err, output)
		fmt.Println(msg)
		if logFile != nil {
			logFile.WriteString("RESULT: FAILED - " + msg + "\n")
		}
		os.Exit(1)
	}

	msg := fmt.Sprintf("Service %s restarted successfully", serviceName)
	fmt.Println(msg)
	if logFile != nil {
		logFile.WriteString("RESULT: SUCCESS - " + msg + "\n")
	}
}

func getServiceStatus(serviceName string, logFile *os.File) {
	output, err := runNSSMCommand([]string{"status", serviceName}, logFile)

	if err != nil {
		msg := fmt.Sprintf("Error getting service status: %v", err)
		fmt.Println(msg)
		if logFile != nil {
			logFile.WriteString("RESULT: FAILED - " + msg + "\n")
		}
		os.Exit(1)
	}

	msg := fmt.Sprintf("Service %s status: %s", serviceName, output)
	fmt.Println(msg)
	if logFile != nil {
		logFile.WriteString("RESULT: SUCCESS - " + msg + "\n")
	}
}