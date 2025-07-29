package errors

import "errors"

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrMissingClaims    = errors.New("missing required claims")
	ErrInvalidAlgorithm = errors.New("invalid signing algorithm")
)
