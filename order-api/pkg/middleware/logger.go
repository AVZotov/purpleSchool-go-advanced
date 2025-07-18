package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logger.SetLevel(logrus.InfoLevel)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.WithFields(logrus.Fields{
			"method":         r.Method,
			"url":            r.URL.String(),
			"remote_addr":    r.RemoteAddr,
			"user_agent":     r.UserAgent(),
			"referer":        r.Referer(),
			"host":           r.Host,
			"proto":          r.Proto,
			"content_length": r.ContentLength,
		}).Info("Incoming request")

		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)

		logger.WithFields(logrus.Fields{
			"method":         r.Method,
			"url":            r.URL.String(),
			"status":         wrapper.StatusCode,
			"status_text":    http.StatusText(wrapper.StatusCode),
			"duration_ms":    duration.Milliseconds(),
			"remote_addr":    r.RemoteAddr,
			"user_agent":     r.UserAgent(),
			"content_length": r.ContentLength,
		}).Info("Incoming response")
	})
}
