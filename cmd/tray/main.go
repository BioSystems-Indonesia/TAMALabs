package main

import (
	_ "embed"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/logger"
	"github.com/energye/systray"
)

//go:embed trayicon.ico
var trayicon []byte

var serverCmd *exec.Cmd

func main() {
	runtime.LockOSThread()

	if os.Getenv(constant.ENVLogLevel) == "" {
		os.Setenv(constant.ENVLogLevel, string(constant.LogLevelInfo))
	}

	provideGlobalLog()

	systray.Run(onReady, func() {})
}

func onReady() {
	systray.SetTitle("TAMALabs")
	systray.SetTooltip("TAMALabs")
	systray.SetIcon(trayicon)

	mStatus := systray.AddMenuItem("Status: Checking...", "Server status")
	systray.AddSeparator()

	mRun := systray.AddMenuItem("Run", "Start server")
	mRestart := systray.AddMenuItem("Restart", "Restart server")
	mStop := systray.AddMenuItem("Stop", "Stop server")
	systray.AddSeparator()

	mOpen := systray.AddMenuItem("Open Browser", "Open Browser")
	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Exit tray")

	// status updater
	go func() {
		for {
			status := "Not Running"
			if isServerRunning() {
				status = "Running"
			}
			mStatus.SetTitle("Status: " + status)
			time.Sleep(2 * time.Second)
		}
	}()

	mOpen.Click(func() {
		openbrowser("http://127.0.0.1:8322")
	})

	mRun.Click(func() {
		startServer()
	})

	mRestart.Click(func() {
		stopServer()
		time.Sleep(500 * time.Millisecond)
		startServer()
	})

	mStop.Click(func() {
		stopServer()
	})

	mQuit.Click(func() {
		stopServer()
		systray.Quit()
	})
}

func startServer() {
	if serverCmd != nil && serverCmd.Process != nil {
		slog.Info("Server already running")
		return
	}

	// Resolve TAMALabs executable path relative to this tray executable.
	// This avoids failures when the tray's working directory is different
	// (for example when launched from Startup shortcuts).
	var exePath string
	if runtime.GOOS == "windows" {
		// Prefer the TAMALabs.exe next to the tray executable
		if exe, err := os.Executable(); err == nil {
			exePath = filepath.Join(filepath.Dir(exe), "TAMALabs.exe")
		} else {
			exePath = "./TAMALabs.exe"
		}
	} else {
		if exe, err := os.Executable(); err == nil {
			exePath = filepath.Join(filepath.Dir(exe), "TAMALabs")
		} else {
			exePath = "./TAMALabs"
		}
	}

	// Log and validate the resolved executable path to help diagnose installer/startup issues
	slog.Info("Resolved server executable path", "exePath", exePath)
	if info, err := os.Stat(exePath); err != nil {
		if os.IsNotExist(err) {
			slog.Error("TAMALabs executable not found", "path", exePath, "err", err)
			return
		}
		slog.Error("Failed to stat TAMALabs executable", "path", exePath, "err", err)
		return
	} else if info.IsDir() {
		slog.Error("Resolved server path is a directory, expected executable", "path", exePath)
		return
	}

	cmd := exec.Command(exePath)
	// ensure the child process has its working directory set to the
	// application folder so relative file reads (like .env) work as
	// expected when started from shortcuts.
	cmd.Dir = filepath.Dir(exePath)

	// Avoid flashing a console window when starting the server from the tray.
	// - Do not attach Stdout/Stderr to the tray process's stdio (that can cause
	//   a console to appear).
	// - On Windows set HideWindow so no console is shown for the child.
	// Redirect child's output to a rotating file under logs so we can inspect it.
	if logFile, err := os.OpenFile(filepath.Join(filepath.Dir(exePath), "logs", "server.stdout.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
		cmd.Stdout = logFile
		cmd.Stderr = logFile
	}

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	if err := cmd.Start(); err != nil {
		slog.Error("Failed to start server", "err", err)
		return
	}
	serverCmd = cmd
	slog.Info("Server started", "pid", cmd.Process.Pid)
}

func stopServer() {
	slog.Info("Stopping TAMALabs.exe via taskkill /IM")

	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/IM", "TAMALabs.exe", "/F", "/T")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		if err := cmd.Run(); err != nil {
			slog.Error("taskkill failed", "err", err)
		} else {
			slog.Info("taskkill succeeded")
		}
	} else {
		cmd := exec.Command("pkill", "TAMALabs")
		if err := cmd.Run(); err != nil {
			slog.Error("pkill failed", "err", err)
		} else {
			slog.Info("pkill succeeded")
		}
	}

	serverCmd = nil
}

func isServerRunning() bool {
	client := &http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Get("http://127.0.0.1:8322/")
	if err == nil {
		resp.Body.Close()
		return true
	}
	return false
}

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		slog.Error("Failed to open browser", "err", err)
	}
}

func provideGlobalLog() {
	// Prefer a dedicated tray log file under the application's logs folder so
	// tray logs don't mix with the main server logs.
	if exe, err := os.Executable(); err == nil {
		logsDir := filepath.Join(filepath.Dir(exe), "logs")
		// Ensure logs directory exists (installer usually creates it, but be safe)
		_ = os.MkdirAll(logsDir, 0o755)
		trayLogPath := filepath.Join(logsDir, "tray.txt")
		l := logger.NewFileLogger(logger.Options{Filename: trayLogPath})
		slog.SetDefault(l)
		slog.Info("Tray logger initialized", "path", trayLogPath)
		return
	}

	// Fallback: use default logger behavior
	l := logger.NewFileLogger(logger.Options{})
	slog.SetDefault(l)
}
