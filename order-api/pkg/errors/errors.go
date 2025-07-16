package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func (e AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	ErrJsonUnmarshal = AppError{
		Code:    "JSON_UNMARSHAL_ERROR",
		Message: "Failed to unmarshal json",
		Status:  http.StatusBadRequest,
	}

	ErrJsonMarshal = AppError{
		Code:    "JSON_MARSHAL_ERROR",
		Message: "Failed to marshal json",
		Status:  http.StatusBadRequest,
	}

	ErrInvalidId = AppError{
		Code:    "INVALID_ID",
		Message: "Invalid id",
		Status:  http.StatusBadRequest,
	}

	ErrNotFound = AppError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
		Status:  http.StatusNotFound,
	}

	ErrRecordNotCreated = AppError{
		Code:    "RECORD_NOT_CREATED",
		Message: "Record not created",
		Status:  http.StatusBadRequest,
	}
)

func NewJsonUnmarshalError(details string) AppError {
	err := ErrJsonUnmarshal
	err.Details = details
	return err
}

func NewJsonMarshalError(details string) AppError {
	err := ErrJsonMarshal
	err.Details = details
	return err
}

func NewInvalidIdError(details string) AppError {
	err := ErrInvalidId
	err.Details = details
	return err
}

func NewNotFoundError(details string) AppError {
	err := ErrNotFound
	err.Details = details
	return err
}

func NewRecordNotCreatedError(details string) AppError {
	err := ErrRecordNotCreated
	err.Details = details
	return err
}

func AsAppError(err error) (AppError, bool) {
	var appErr AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}

func Wrap(message string, err error) error {
	if err == nil {
		return nil
	}
	context := getCallerPath()
	return fmt.Errorf("%s: %s: %w", context, message, err)
}

func getCallerPath() string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	funcName := runtime.FuncForPC(pc).Name()
	if lastDot := strings.LastIndex(funcName, "."); lastDot >= 0 {
		funcName = funcName[lastDot+1:]
	}

	fileName := file
	if lastSlash := strings.LastIndex(fileName, "/"); lastSlash >= 0 {
		fileName = fileName[lastSlash+1:]
	}

	return fmt.Sprintf("%s:%d %s", fileName, line, funcName)
}
