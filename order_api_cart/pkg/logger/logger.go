package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	pkgTp "order_api_cart/pkg/types"
	"os"
	"time"
)

const (
	RequestIdField = string(pkgTp.CtxRequestId)
	UserPhoneField = string(pkgTp.CtxUserPhone)
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

func LogWithContext(ctx context.Context, level logrus.Level, message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	if requestID := ctx.Value(pkgTp.CtxRequestId); requestID != nil {
		if id, ok := requestID.(string); ok && id != "" {
			fields[RequestIdField] = id
		}
	}

	if userPhone := ctx.Value(pkgTp.CtxUserPhone); userPhone != nil {
		if phone, ok := userPhone.(string); ok && phone != "" {
			fields[UserPhoneField] = phone
		}
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
	LogWithContext(ctx, logrus.InfoLevel, message, fields)
}

func ErrorWithRequestID(ctx context.Context, message string, fields logrus.Fields) {
	LogWithContext(ctx, logrus.ErrorLevel, message, fields)
}
