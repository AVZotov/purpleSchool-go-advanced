package middleware

import (
	"context"
	"net/http"
	pkgCtx "order_api_auth/pkg/context"
	pkgHeader "order_api_auth/pkg/http"
	"order_api_auth/pkg/jwt"
	"strings"
)

const Prefix = "Bearer "

func AuthMiddleware(next http.Handler, jwtSecret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(pkgHeader.AuthorizationHeader)
		if !strings.HasPrefix(authHeader, Prefix) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		authToken := strings.TrimPrefix(authHeader, Prefix)
		if authToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		isValid, phone := jwt.ParsePhone(authToken, jwtSecret)
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), pkgCtx.CtxUserPhone, phone)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
