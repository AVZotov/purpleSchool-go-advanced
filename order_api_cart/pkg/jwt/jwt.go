package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	pkgErr "order_api_cart/pkg/errors"
)

func Create(secret, phone string) (string, error) {
	if secret == "" {
		return "", pkgErr.ErrConfigMissing
	}
	if phone == "" {
		return "", pkgErr.ErrInvalidPhone
	}

	jwtSecret := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": phone,
	})

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("%w %v", pkgErr.ErrServiceUnavailable)
	}

	return signedToken, nil
}

func ParseValidate(tokenString, secret string) (string, error) {
	if tokenString == "" || secret == "" {
		return "", pkgErr.ErrInvalidAuth
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, pkgErr.ErrInvalidAuth
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("token validation failed: %w", pkgErr.ErrInvalidAuth)
	}

	if !token.Valid {
		return "", pkgErr.ErrInvalidAuth
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", pkgErr.ErrInvalidAuth
	}

	phoneInterface, exists := claims["phone"]
	if !exists {
		return "", pkgErr.ErrInvalidAuth
	}

	phone, ok := phoneInterface.(string)
	if !ok || phone == "" {
		return "", pkgErr.ErrInvalidAuth
	}

	return phone, nil
}
