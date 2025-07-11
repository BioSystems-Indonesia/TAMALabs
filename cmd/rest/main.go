package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"time"

	"github.com/energye/systray"
	"github.com/oibacidem/lims-hl-seven/internal/app"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/util"
	"github.com/oibacidem/lims-hl-seven/pkg/logger"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
	"github.com/tarm/serial"
)

var (
	flagDev      = flag.Bool("development", false, "development mode")
	flagLogLevel = flag.String("log-level", string(constant.LogLevelInfo), "log level: debug, info, warn, error")
)

// version is set at build time

var version = ""

//go:embed trayicon.ico
var trayicon []byte

func main() {
	config := &serial.Config{
		Name:        "COM6", // Ganti jika pakai Linux: "/dev/ttyS6"
		Baud:        115200,
		Size:        8,
		StopBits:    serial.Stop1,
		Parity:      serial.ParityNone,
		ReadTimeout: 0,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		os.Exit(1) // Diam-diam keluar jika gagal
	}

	flag.Parse()

	if *flagDev {
		os.Setenv(constant.ENVKey, string(constant.EnvDevelopment))
	} else {
		os.Setenv(constant.ENVKey, string(constant.EnvProduction))
	}

	os.Setenv(constant.ENVLogLevel, string(validateLogLevel(*flagLogLevel)))
	os.Setenv(constant.ENVVersion, version)

	provideGlobalLog()

	server := app.InitRestApp()
	serial := app.InitSerialNCC3300App()

	serial.Handle(port)
	go openb()
	go opensystray(server)
	server.Serve()
}

func openb() {
	if version == "" || util.IsDevelopment() {
		return
	}

	time.Sleep(3 * time.Second)
	openbrowser("http://127.0.0.1:8322")
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
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		slog.Error("error opening browser", "error", err)
	}
}

func opensystray(server server.RestServer) {
	systray.Run(func() {
		systray.SetTitle("LIMS HL Seven")
		systray.SetTooltip("LIMS HL Seven")
		systray.SetIcon(trayicon)

		systray.AddMenuItem("Open Browser", "Open Browser").Click(func() {
			openbrowser("http://127.0.0.1:8322")
		})

		systray.AddMenuItem("Quit", "Stop Server").Click(func() {
			if err := server.Stop(); err != nil {
				slog.Error("Error stopping server:", "err", err)
			} else {
				slog.Info("Server stopped successfully")
			}
			systray.Quit()
		})
	}, func() {})
}

func validateLogLevel(logLevel string) constant.LogLevel {
	if !slices.Contains(constant.ValidLogLevels, constant.LogLevel(logLevel)) {
		panic(fmt.Sprintf("invalid log level: %s", logLevel))
	}

	return constant.LogLevel(logLevel)
}

func provideGlobalLog() {
	l := logger.NewFileLogger(logger.Options{})

	slog.SetDefault(l)
	slog.Info("version", "version", version)
}
