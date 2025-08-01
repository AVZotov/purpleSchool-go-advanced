package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	pkgLogger "order_simple/pkg/logger"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get(pkgLogger.RequestIDHeader)

		pkgLogger.Logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     r.Method,
			"path":       r.URL.Path,
			"query":      r.URL.RawQuery,
			"user_agent": r.UserAgent(),
			"ip":         getClientIP(r),
			"type":       pkgLogger.HttpRequestStart,
		}).Info("HTTP request started")

		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)

		pkgLogger.Logger.WithFields(logrus.Fields{
			"request_id":     requestID,
			"method":         r.Method,
			"path":           r.URL.Path,
			"status_code":    wrapper.StatusCode,
			"duration_ms":    duration.Milliseconds(),
			"user_agent":     r.UserAgent(),
			"ip":             getClientIP(r),
			"type":           pkgLogger.HttpRequestEnd,
			"content_length": r.ContentLength,
		}).Info("HTTP request completed")
	})
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get(pkgLogger.RequestIPHeader)
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
