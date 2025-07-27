package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrMissingClaims    = errors.New("missing required claims")
	ErrInvalidAlgorithm = errors.New("invalid signing algorithm")
)

func Create(secret, phone string) (string, error) {
	if secret == "" || phone == "" {
		return "", errors.New("secret and phone can not be empty")
	}

	jwtSecret := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": phone,
	})

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func ParseValidate(tokenString, secret string) (string, error) {
	if tokenString == "" && secret == "" {
		return "", ErrInvalidToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w %v", ErrInvalidAlgorithm, token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return "", ErrInvalidSignature
		}
		return "", fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrMissingClaims
	}

	phoneInterface, exists := claims["phone"]
	if !exists {
		return "", ErrMissingClaims
	}

	phone, ok := phoneInterface.(string)
	if !ok || phone == "" {
		return "", ErrMissingClaims
	}

	return phone, nil
}
