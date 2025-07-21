package middleware

import (
	"crypto/rand"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	pkgLogger "order_simple/pkg/logger"
	"strings"
	"time"
)

const RequestIDHeader = "X-Request-ID"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestID string

		requestID = r.Header.Get(RequestIDHeader)

		if requestID == "" {
			requestID = generateRequestID()
		}

		requestID = sanitizeRequestID(requestID)

		r.Header.Set(RequestIDHeader, requestID)

		w.Header().Set(RequestIDHeader, requestID)

		pkgLogger.Logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     r.Method,
			"path":       r.URL.Path,
			"type":       "request_id_set",
		}).Debug("Request ID assigned")

		next.ServeHTTP(w, r)
	})
}

func generateRequestID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("req_%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("req_%x", bytes)
}

func sanitizeRequestID(requestID string) string {
	requestID = strings.ReplaceAll(requestID, "\n", "")
	requestID = strings.ReplaceAll(requestID, "\r", "")
	requestID = strings.TrimSpace(requestID)

	if len(requestID) > 64 {
		requestID = requestID[:64]
	}

	return requestID
}
