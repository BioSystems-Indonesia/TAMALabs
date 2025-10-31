package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"slices"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/app"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/logger"
)

var (
	flagDev      = flag.Bool("development", false, "development mode")
	flagLogLevel = flag.String("log-level", string(constant.LogLevelInfo), "log level: debug, info, warn, error")
)

var version = "" // diset saat build

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

	slog.Info("Starting REST server")
	server.Serve()
}

func showErrorOnPanic() {
	if err := recover(); err != nil {
		slog.Error("Error on startup", "err", err)
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
