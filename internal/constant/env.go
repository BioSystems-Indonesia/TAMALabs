package constant

import (
	"fmt"
	"log/slog"
)

const (
	ENVKey      = "LIMS_ENV"
	ENVLogLevel = "LIMS_LOG_LEVEL"
	ENVVersion  = "LIMS_VERSION"
)

type Env string

const (
	EnvDevelopment Env = "development"
	EnvProduction  Env = "production"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

var ValidLogLevels = []LogLevel{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError}

func (l LogLevel) GetSlogLevel() slog.Level {
	switch l {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		panic(fmt.Sprintf("invalid log level: %s", l))
	}
}
