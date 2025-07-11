package utils

import (
	"fmt"
	"link_shortener/pkg/logger"
	"runtime"
	"strings"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// GetContext returns caller function name in format package/function
func GetContext() string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	funcName := runtime.FuncForPC(pc).Name()
	if lastDot := strings.LastIndex(funcName, "."); lastDot >= 0 {
		funcName = funcName[lastDot+1:]
	}

	fileName := file
	if lastSlash := strings.LastIndex(fileName, "/"); lastSlash >= 0 {
		fileName = fileName[lastSlash+1:]
	}

	return fmt.Sprintf("%s:%d %s", fileName, line, funcName)
}

// LogContext short macros to return wrapped logging info with path to caller
// package.function
func _(logger logger.Logger, level string, msg string, args ...any) {
	context := GetContext()
	allArgs := append([]any{"context", context}, args...)

	switch level {
	case LevelDebug:
		logger.Debug(msg, allArgs...)
	case LevelInfo:
		logger.Info(msg, allArgs...)
	case LevelWarn:
		logger.Warn(msg, allArgs...)
	case LevelError:
		logger.Error(msg, allArgs...)
	default:
		logger.Info(msg, allArgs...)
	}
}
