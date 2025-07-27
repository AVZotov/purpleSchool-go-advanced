package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	pkgCtx "order_api_auth/pkg/context"
	pkgHeader "order_api_auth/pkg/http"
	pkgJWT "order_api_auth/pkg/jwt"
	pkgLogger "order_api_auth/pkg/logger"
	"strings"
)

type AuthErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code"`
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := extractBearerToken(r)
			if err != nil {
				writeAuthError(w, "missing_token", "Authorization header is required")
			}

			phone, err := pkgJWT.ParseValidate(token, jwtSecret)
			if err != nil {
				handleJWTError(w, r, err)
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

	authHeader := r.Header.Get(pkgHeader.AuthorizationHeader)
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", errors.New("invalid authorization format")
	}

	token := strings.TrimPrefix(authHeader, BearerPrefix)
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}

func handleJWTError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()

	switch {
	case errors.Is(err, pkgJWT.ErrInvalidSignature):
		pkgLogger.ErrorWithRequestID(ctx, "Invalid token signature", logrus.Fields{"error": err.Error()})
		writeAuthError(w, "invalid_signature", "Invalid token signature")
	case errors.Is(err, pkgJWT.ErrInvalidAlgorithm):
		pkgLogger.ErrorWithRequestID(ctx, "Invalid signing algorithm", logrus.Fields{"error": err.Error()})
		writeAuthError(w, "invalid_algorithm", "Invalid token algorithm")
	case errors.Is(err, pkgJWT.ErrMissingClaims):
		pkgLogger.ErrorWithRequestID(ctx, "Missing required claims", logrus.Fields{"error": err.Error()})
		writeAuthError(w, "missing_claims", "Token missing required claims")
	default:
		pkgLogger.ErrorWithRequestID(ctx, "Token validation failed", logrus.Fields{"error": err.Error()})
		writeAuthError(w, "invalid_token", "Invalid token")
	}
}

func writeAuthError(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	response := AuthErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}

	_ = json.NewEncoder(w).Encode(response)
}
