package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	pkgHttp "order_api_cart/pkg/http"
	pkgLogger "order_api_cart/pkg/logger"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := r.Context()

		pkgLogger.InfoWithRequestID(ctx, "HTTP request received", logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"query":      r.URL.RawQuery,
			"user_agent": r.UserAgent(),
			"ip":         getClientIP(r),
		})

		wrapper := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)

		pkgLogger.InfoWithRequestID(ctx, "HTTP request received", logrus.Fields{
			"method":         r.Method,
			"path":           r.URL.Path,
			"status_code":    wrapper.statusCode,
			"duration_ms":    duration.Milliseconds(),
			"user_agent":     r.UserAgent(),
			"ip":             getClientIP(r),
			"content_length": r.ContentLength,
		})
	})
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get(pkgHttp.RequestIPHeader)
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
