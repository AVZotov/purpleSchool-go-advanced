package logger

import (
	"log/slog"
	"os"
)

func NewLogger(devEnv string) *slog.Logger {
	switch devEnv {
	case "dev":
		return slog.New(slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		return slog.New(slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	default:
		return slog.New(slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}
