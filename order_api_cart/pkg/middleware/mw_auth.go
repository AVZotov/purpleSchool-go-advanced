package middleware

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	pkgCtx "order_api_cart/pkg/context"
	pkgErr "order_api_cart/pkg/errors"
	pkgHdr "order_api_cart/pkg/http"
	pkgJWT "order_api_cart/pkg/jwt"
	pkgLogger "order_api_cart/pkg/logger"
	"strings"
)

type AuthErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := extractBearerToken(r)
			if err != nil {
				errCode := pkgErr.GetStatusCode(err)
				errMessage := http.StatusText(errCode)
				pkgLogger.ErrorWithRequestID(ctx, errMessage, logrus.Fields{
					"error": err.Error(),
				})
				writeAuthError(w, errCode, errMessage)
			}

			phone, err := pkgJWT.ParseValidate(token, jwtSecret)
			if err != nil {
				handleJWTError(ctx, w, err)
				return
			}

			pkgLogger.InfoWithRequestID(ctx, "User authenticated successfully", logrus.Fields{
				"phone": phone,
			})

			ctx = context.WithValue(r.Context(), pkgCtx.CtxUserPhone, phone)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}

}

func extractBearerToken(r *http.Request) (string, error) {
	const BearerPrefix = "Bearer "

	authHeader := r.Header.Get(pkgHdr.AuthorizationHeader)
	if authHeader == "" {
		return "", pkgErr.ErrInvalidAuth
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", pkgErr.ErrInvalidAuth
	}

	token := strings.TrimPrefix(authHeader, BearerPrefix)
	if token == "" {
		return "", pkgErr.ErrInvalidAuth
	}

	return token, nil
}

func handleJWTError(ctx context.Context, w http.ResponseWriter, err error) {
	errCode := pkgErr.GetStatusCode(err)
	pkgLogger.ErrorWithRequestID(ctx, http.StatusText(errCode), logrus.Fields{"error": err.Error()})
	writeAuthError(w, errCode, http.StatusText(errCode))
}

func writeAuthError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := AuthErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}

	_ = json.NewEncoder(w).Encode(response)
}
