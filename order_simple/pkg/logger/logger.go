package logger

import (
	"github.com/sirupsen/logrus"
	"net/http"
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

func LogWithRequestID(r *http.Request, level logrus.Level, message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	requestID := GetRequestID(r)
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

func InfoWithRequestID(r *http.Request, message string, fields logrus.Fields) {
	LogWithRequestID(r, logrus.InfoLevel, message, fields)
}

func ErrorWithRequestID(r *http.Request, message string, fields logrus.Fields) {
	LogWithRequestID(r, logrus.ErrorLevel, message, fields)
}

func WarnWithRequestID(r *http.Request, message string, fields logrus.Fields) {
	LogWithRequestID(r, logrus.WarnLevel, message, fields)
}

func GetRequestID(r *http.Request) string {
	return r.Header.Get(RequestIDHeader)
}

func GetClientIP(r *http.Request) string {
	forwarded := r.Header.Get(RequestIPHeader)
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
