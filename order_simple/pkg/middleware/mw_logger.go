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

		pkgLogger.InfoWithRequestID(r, "HTTP request start", logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"query":      r.URL.RawQuery,
			"user_agent": r.UserAgent(),
			"ip":         pkgLogger.GetClientIP(r),
			"type":       pkgLogger.HttpRequestStart,
		})

		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		pkgLogger.InfoWithRequestID(r, "HTTP request completed", logrus.Fields{
			"method":         r.Method,
			"path":           r.URL.Path,
			"status_code":    wrapper.StatusCode,
			"duration_ms":    time.Since(start).Milliseconds(),
			"user_agent":     r.UserAgent(),
			"ip":             pkgLogger.GetClientIP(r),
			"type":           pkgLogger.HttpRequestEnd,
			"content_length": r.ContentLength,
		})
	})
}
