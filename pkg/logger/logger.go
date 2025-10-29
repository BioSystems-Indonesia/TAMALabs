package logger

import (
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/util"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

const (
	logMaxSizeMB  = 10
	logMaxBackup  = 10
	logMaxAgeDays = 1
	logFilename   = "main.log"
)

func GetDefaultLogFile() string {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	defaultLogFile := path.Join(path.Dir(executable), "logs", logFilename)
	return defaultLogFile
}

func NewFileLogger(opts Options) *slog.Logger {
	opts = checkOptions(opts)

	filename := opts.Filename
	if !filepath.IsAbs(filename) {
		filename = GetDefaultLogFile()
	}

	logRotator := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    opts.MaxSize,
		MaxBackups: opts.MaxBackups,
		MaxAge:     opts.MaxAge,
		Compress:   true,
	}

	logLevel := constant.LogLevel(os.Getenv(constant.ENVLogLevel))
	writer := newWriter(logRotator)

	return slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
		AddSource: false,
		Level:     logLevel.GetSlogLevel(),
	}))
}

func newWriter(logRotator *lumberjack.Logger) io.Writer {
	if util.IsDevelopment() {
		return io.MultiWriter(os.Stdout, logRotator)
	}

	return logRotator
}

func checkOptions(opts Options) Options {
	if opts.MaxSize == 0 {
		opts.MaxSize = logMaxSizeMB
	}
	if opts.MaxBackups == 0 {
		opts.MaxBackups = logMaxBackup
	}
	if opts.MaxAge == 0 {
		opts.MaxAge = logMaxAgeDays
	}
	if opts.Filename == "" {
		opts.Filename = logFilename
	}

	return opts
}
