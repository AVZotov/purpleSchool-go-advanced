package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := log.New()
		logger.SetFormatter(&log.JSONFormatter{})
		logger.SetOutput(os.Stdout)
		logger.WithFields(log.Fields{
			"time":       time.Now().Format("2006-01-02 15:04:05"),
			"status":     "",
			"method":     r.Method,
			"host":       r.Host,
			"url":        r.URL.String(),
			"referer":    r.Referer(),
			"user_agent": r.UserAgent(),
			"latency":    time.Since(start),
		})

		logger.Info("Message")
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, r)
		logger.WithFields(log.Fields{})
	})
}
