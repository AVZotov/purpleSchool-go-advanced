package errors

import "errors"

// =============================================================================
// AUTHENTICATION & AUTHORIZATION ERRORS - Token, Signature, Claims, Session, SMS/Code errors
// =============================================================================

var (
	ErrMissingAuthHeader = errors.New("missing authentication header")

	ErrInvalidToken = errors.New("invalid token")
	ErrMissingToken = errors.New("token missing")

	ErrInvalidSignature = errors.New("invalid signature")
	ErrInvalidAlgorithm = errors.New("invalid signing algorithm")

	ErrMissingClaims = errors.New("missing required claims")
	ErrInvalidClaims = errors.New("invalid claims format")

	ErrInvalidSession  = errors.New("invalid session")
	ErrSessionNotFound = errors.New("session not found")

	ErrInvalidCode = errors.New("invalid verification code")
)

// =============================================================================
// DATABASE ERRORS - Connection, Migration, Query, Constraint errors
// =============================================================================

var (
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrDatabaseTimeout    = errors.New("database operation timeout")

	ErrMigrationFailed = errors.New("database migration failed")

	ErrRecordNotFound    = errors.New("record not found")
	ErrRecordExists      = errors.New("record already exists")
	ErrQueryFailed       = errors.New("database query failed")
	ErrTransactionFailed = errors.New("database transaction failed")

	ErrUniqueViolation     = errors.New("unique constraint violation")
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")
)

// =============================================================================
// VALIDATION ERRORS - General & Phone validation and General format errors
// =============================================================================

var (
	ErrValidation   = errors.New("validation failed")
	ErrInvalidInput = errors.New("invalid input data")

	ErrInvalidPhone  = errors.New("invalid phone number")
	ErrPhoneFormat   = errors.New("phone number format invalid")
	ErrPhoneRequired = errors.New("phone number required")

	ErrInvalidFormat     = errors.New("invalid data format")
	ErrInvalidLength     = errors.New("invalid data length")
	ErrInvalidCharacters = errors.New("invalid characters")
)

// =============================================================================
// HTTP & JSON ERRORS - JSON processing, HTTP request and response errors
// =============================================================================

var (
	ErrEncodingJSON = errors.New("encoding JSON failed")
	ErrDecodingJSON = errors.New("decoding JSON failed")

	ErrInvalidRequest = errors.New("invalid HTTP request")
	ErrMissingHeader  = errors.New("missing required header")

	ErrResponseFailed  = errors.New("HTTP response failed")
	ErrInvalidResponse = errors.New("invalid response format")
)

// =============================================================================
// BUSINESS LOGIC ERRORS (DOMAIN SPECIFIC) - User, Order, Product errors
// =============================================================================

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")

	ErrOrderNotFound = errors.New("order not found")
	ErrOrderAccess   = errors.New("access denied to order")

	ErrProductNotFound = errors.New("product not found")
)

// =============================================================================
// EXTERNAL SERVICE ERRORS - JWT, SessionID
// =============================================================================

var (
	ErrCreatingToken = errors.New("token creation failed")

	ErrGeneratingSessionID = errors.New("sessionID generation failed")
)

// =============================================================================
// SYSTEM & INFRASTRUCTURE ERRORS - Configuration
// =============================================================================

var (
	ErrConfigMissing      = errors.New("configuration missing")
	ErrConfigInvalid      = errors.New("configuration invalid")
	ErrEnvironmentInvalid = errors.New("invalid environment")
)

func IsAuthenticationError(err error) bool {
	authErrors := []error{
		ErrInvalidToken, ErrMissingToken, ErrSessionNotFound,
		ErrInvalidSignature, ErrInvalidAlgorithm,
		ErrMissingClaims, ErrInvalidClaims, ErrMissingAuthHeader,
		ErrInvalidSession, ErrInvalidCode,
	}

	for _, authErr := range authErrors {
		if errors.Is(err, authErr) {
			return true
		}
	}
	return false
}

func IsValidationError(err error) bool {
	validationErrors := []error{
		ErrValidation, ErrInvalidInput,
		ErrInvalidPhone, ErrPhoneFormat, ErrPhoneRequired,
		ErrInvalidFormat, ErrInvalidLength, ErrInvalidCharacters,
	}

	for _, valErr := range validationErrors {
		if errors.Is(err, valErr) {
			return true
		}
	}
	return false
}

func GetHTTPStatusCode(err error) int {
	switch {
	// 400 Bad Request
	case IsValidationError(err):
		return 400
	case errors.Is(err, ErrInvalidRequest):
		return 400

	// 401 Unauthorized
	case IsAuthenticationError(err):
		return 401

	// 403 Forbidden
	case errors.Is(err, ErrOrderAccess):
		return 403

	// 404 Not Found
	case errors.Is(err, ErrRecordNotFound):
		return 404
	case errors.Is(err, ErrUserNotFound):
		return 404
	case errors.Is(err, ErrOrderNotFound):
		return 404

	// 409 Conflict
	case errors.Is(err, ErrUserExists):
		return 409
	case errors.Is(err, ErrUniqueViolation):
		return 409

	// 500 Internal Server Error
	default:
		return 500
	}
}
