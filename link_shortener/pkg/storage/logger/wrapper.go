package logger

import (
	t "link_shortener/pkg/storage/types"
	"log/slog"
)

type StorageLoggerWrapper struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) t.Logger {
	return &StorageLoggerWrapper{logger: logger}
}

func (w *StorageLoggerWrapper) Debug(msg string, args ...any) {
	w.logger.Debug(msg, args...)
}

func (w *StorageLoggerWrapper) Info(msg string, args ...any) {
	w.logger.Info(msg, args...)
}

func (w *StorageLoggerWrapper) Warn(msg string, args ...any) {
	w.logger.Warn(msg, args...)
}

func (w *StorageLoggerWrapper) Error(msg string, args ...any) {
	w.logger.Error(msg, args...)
}

func (w *StorageLoggerWrapper) With(args ...any) t.Logger {
	return &StorageLoggerWrapper{logger: w.logger.With(args...)}
}
