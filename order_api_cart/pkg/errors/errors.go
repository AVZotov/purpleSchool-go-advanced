package errors

import "errors"

// Authentication errors
var (
	ErrInvalidAuth     = errors.New("authentication failed")
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidCode     = errors.New("invalid verification code")
)

// Database errors
var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrRecordExists      = errors.New("record already exists")
	ErrQueryFailed       = errors.New("database query failed")
	ErrTransactionFailed = errors.New("database transaction failed")
)

// Validation errors
var (
	ErrValidation   = errors.New("validation failed")
	ErrInvalidInput = errors.New("invalid input data")
	ErrInvalidPhone = errors.New("invalid phone number")
)

// HTTP base errors
var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrEncodingJSON   = errors.New("json encoding failed")
	ErrDecodingJSON   = errors.New("json decoding failed")
)

// Business logic errors
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrOrderNotFound   = errors.New("order not found")
	ErrOrderAccess     = errors.New("access denied to order")
	ErrProductNotFound = errors.New("product not found")
)

// System errors
var (
	ErrConfigMissing      = errors.New("configuration missing")
	ErrServiceUnavailable = errors.New("external service unavailable")
)

func IsAuthError(err error) bool {
	return errors.Is(err, ErrInvalidAuth) ||
		errors.Is(err, ErrSessionNotFound) ||
		errors.Is(err, ErrInvalidCode)
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrRecordNotFound) ||
		errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrOrderNotFound) ||
		errors.Is(err, ErrProductNotFound)
}

func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation) ||
		errors.Is(err, ErrInvalidInput) ||
		errors.Is(err, ErrInvalidPhone)
}

func IsConflictError(err error) bool {
	return errors.Is(err, ErrRecordExists) ||
		errors.Is(err, ErrUserExists)
}

// GetStatusCode is a function returning error code mapped to appErrors
func GetStatusCode(err error) int {
	switch {
	case IsValidationError(err):
		return 400
	case errors.Is(err, ErrInvalidRequest):
		return 400
	case IsAuthError(err):
		return 401
	case errors.Is(err, ErrOrderAccess):
		return 403
	case IsNotFoundError(err):
		return 404
	case IsConflictError(err):
		return 409
	case errors.Is(err, ErrServiceUnavailable):
		return 503
	default:
		return 500
	}
}
