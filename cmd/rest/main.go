package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"slices"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/app"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/logger"
)

var (
	flagDev      = flag.Bool("development", false, "development mode")
	flagLogLevel = flag.String("log-level", string(constant.LogLevelInfo), "log level: debug, info, warn, error")
)

// version is set at build time

var version = ""

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
		showErrorMessage("Cannot open LIS", fmt.Sprintf("%v", err))
		os.Exit(1)
	}
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
