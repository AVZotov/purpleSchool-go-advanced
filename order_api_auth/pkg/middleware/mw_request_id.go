package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	pkgCtx "order_api_auth/pkg/context"
	pkgHttp "order_api_auth/pkg/http"
	pkgLogger "order_api_auth/pkg/logger"
	"strings"
	"time"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestID string

		requestID = r.Header.Get(pkgHttp.RequestIDHeader)

		if requestID == "" {
			requestID = generateRequestID()
		}

		requestID = sanitizeRequestID(requestID)

		r.Header.Set(pkgHttp.RequestIDHeader, requestID)

		w.Header().Set(pkgHttp.RequestIDHeader, requestID)

		ctx := context.WithValue(r.Context(), pkgCtx.CtxRequestId, requestID)

		r = r.WithContext(ctx)

		pkgLogger.Logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     r.Method,
			"path":       r.URL.Path,
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
