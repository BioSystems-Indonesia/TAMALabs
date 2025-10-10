package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/app"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/util"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/logger"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/server"
	"github.com/energye/systray"
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
	defer showErrorOnPanic()

	flag.Parse()

	if *flagDev {
		os.Setenv(constant.ENVKey, string(constant.EnvDevelopment))
	} else {
		os.Setenv(constant.ENVKey, string(constant.EnvProduction))
	}

	os.Setenv(constant.ENVLogLevel, string(validateLogLevel(*flagLogLevel)))
	os.Setenv(constant.ENVVersion, version)

	provideGlobalLog()

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	server := app.InitRestApp()
	go func() {
		slog.Info("Initializing Khanza Canal Handler...")
		startCanalHandler()
	}()

	go openb()
	go opensystray(server)
	server.Serve()
}

func startCanalHandler() {
	time.Sleep(5 * time.Second)

	canalHandler := app.InitCanalHandler()

	if canalHandler == nil {
		slog.Error("Failed to create Canal Handler - dependency injection failed")
		return
	}

	slog.Info("Canal Handler initialized successfully with all dependencies")
	canalHandler.StartCanalHandler()
}

func showErrorOnPanic() {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case error:
			slog.Error("Error on startup", slog.String("error", e.Error()))
		default:
			slog.Error("Error on startup", slog.String("error", fmt.Sprintf("%v", err)))
		}
		// showErrorMessage("Cannot open LIS", fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}

func openb() {
	if version == "" || util.IsDevelopment() {
		return
	}

	time.Sleep(3 * time.Second)

	service := app.InitService()
	err := service.Check()
	if err != nil {
		openbrowser("http://127.0.0.1:8322/license")
	}

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
