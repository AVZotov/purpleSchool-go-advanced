package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"-"`
}

func (e AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	ErrValidation = AppError{
		Code:    "VALIDATION_ERROR",
		Message: "Invalid input data",
		Status:  http.StatusBadRequest,
	}

	ErrStructValidation = AppError{
		Code:    "VALIDATION_ERROR",
		Message: "Object validation failed",
		Status:  http.StatusInternalServerError,
	}

	ErrNotFound = AppError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
		Status:  http.StatusNotFound,
	}

	ErrInternal = AppError{
		Code:    "INTERNAL_ERROR",
		Message: "Internal server error",
		Status:  http.StatusInternalServerError,
	}

	ErrEmailSending = AppError{
		Code:    "EMAIL_SENDING_ERROR",
		Message: "Failed to send email",
		Status:  http.StatusInternalServerError,
	}

	ErrStorageOperation = AppError{
		Code:    "STORAGE_ERROR",
		Message: "Storage operation failed",
		Status:  http.StatusInternalServerError,
	}

	ErrJsonParse = AppError{
		Code:    "JSON_PARSE_ERROR",
		Message: "Json parse failed",
		Status:  http.StatusBadRequest,
	}
)

func NewValidationError(details string) AppError {
	err := ErrValidation
	err.Details = details
	return err
}

func NewStructValidationError(details string) AppError {
	err := ErrStructValidation
	err.Details = details
	return err
}

func NewNotFoundError(details string) AppError {
	err := ErrNotFound
	err.Details = details
	return err
}

func NewInternalError(details string) AppError {
	err := ErrInternal
	err.Details = details
	return err
}

func NewEmailSendingError(details string) AppError {
	err := ErrEmailSending
	err.Details = details
	return err
}

func NewStorageError(details string) AppError {
	err := ErrStorageOperation
	err.Details = details
	return err
}

func NewJsonParseError(details string) AppError {
	err := ErrJsonParse
	err.Details = details
	return err
}

func Wrap(message string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

func IsAppError(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr)
}

func AsAppError(err error) (AppError, bool) {
	var appErr AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}
