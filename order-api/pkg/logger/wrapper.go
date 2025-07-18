package logger

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
)

const ProjectName = "order-api"

type SmartWrapper struct {
	logger *slog.Logger
}

// NewWrapper return [Logger] wrapped with path to caller function
func NewWrapper(logger *slog.Logger) Logger {
	return &SmartWrapper{logger: logger}
}

func (w *SmartWrapper) Debug(msg string, args ...any) {
	//context := w.getContext()
	//allArgs := append([]any{"context", context}, args...)
	w.logger.Debug(msg, args...)
}

func (w *SmartWrapper) Info(msg string, args ...any) {
	//context := w.getContext()
	//allArgs := append([]any{"context", context}, args...)
	w.logger.Info(msg, args...)
}

func (w *SmartWrapper) Warn(msg string, args ...any) {
	context := w.getContext()
	allArgs := append([]any{"context", context}, args...)
	w.logger.Warn(msg, allArgs...)
}

func (w *SmartWrapper) Error(msg string, args ...any) {
	context := w.getContext()
	allArgs := append([]any{"context", context}, args...)
	w.logger.Error(msg, allArgs...)
}

func (w *SmartWrapper) With(args ...any) Logger {
	return &SmartWrapper{logger: w.logger.With(args...)}
}

func (w *SmartWrapper) getContext() string {
	// skip=2 потому что: getContext -> Debug/Info/Error -> пользовательская функция
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}

	funcName := runtime.FuncForPC(pc).Name()
	funcName = w.cleanFuncName(funcName)

	relativePath := w.getRelativePath(file)

	return fmt.Sprintf("%s:%d %s", relativePath, line, funcName)
}

// cleanFuncName cleans full path and returning only caller function name
func (w *SmartWrapper) cleanFuncName(fullName string) string {
	if lastDot := strings.LastIndex(fullName, "."); lastDot >= 0 {
		fullName = fmt.Sprintf("caller=%s", fullName[lastDot+1:])
	}

	return fullName
}

func (w *SmartWrapper) getRelativePath(fullPath string) string {

	parts := strings.Split(fullPath, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == ProjectName {
			if i+1 < len(parts) {
				return strings.Join(parts[i+1:], "/")
			}
		}
	}

	return filepath.Base(fullPath)
}
