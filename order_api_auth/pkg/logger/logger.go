package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger *logrus.Logger

func Init() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetOutput(os.Stdout)
}

func LogWithRequestID(ctx context.Context, level logrus.Level, message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	requestID := ctx.Value("request_id")
	if requestID != "" {
		fields["request_id"] = requestID
	}

	entry := Logger.WithFields(fields)
	switch level {
	case logrus.DebugLevel:
		entry.Debug(message)
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.FatalLevel:
		entry.Fatal(message)
	default:
		entry.Info(message)
	}
}

func InfoWithRequestID(ctx context.Context, message string, fields logrus.Fields) {
	LogWithRequestID(ctx, logrus.InfoLevel, message, fields)
}

func ErrorWithRequestID(ctx context.Context, message string, fields logrus.Fields) {
	LogWithRequestID(ctx, logrus.ErrorLevel, message, fields)
}
