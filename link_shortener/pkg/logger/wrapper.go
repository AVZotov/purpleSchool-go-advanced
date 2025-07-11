package logger

import (
	"log/slog"
)

type Wrapper struct {
	logger *slog.Logger
}

func NewWrapper(logger *slog.Logger) Logger {
	return &Wrapper{logger: logger}
}

func (w *Wrapper) Debug(msg string, args ...any) {
	w.logger.Debug(msg, args...)
}

func (w *Wrapper) Info(msg string, args ...any) {
	w.logger.Info(msg, args...)
}

func (w *Wrapper) Warn(msg string, args ...any) {
	w.logger.Warn(msg, args...)
}

func (w *Wrapper) Error(msg string, args ...any) {
	w.logger.Error(msg, args...)
}

func (w *Wrapper) With(args ...any) Logger {
	return &Wrapper{logger: w.logger.With(args...)}
}
