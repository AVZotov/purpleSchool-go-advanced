package base

import (
	"encoding/json"
	"errors"
	"net/http"
	pkgErrors "order/pkg/errors"
	pkgLogger "order/pkg/logger"
)

type Handler struct {
	Logger pkgLogger.Logger
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

func (h *Handler) WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: status < 400,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		response.Error = &ErrorInfo{
			Code:    http.StatusText(status),
			Message: err.Error(),
			Details: http.StatusText(status),
		}
	}
}

func (h *Handler) WriteError(w http.ResponseWriter, err error) {
	var status int
	var errorInfo ErrorInfo

	if appErr, ok := pkgErrors.AsAppError(err); ok {
		status = appErr.Status
		errorInfo = ErrorInfo{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		}
	} else {
		status = http.StatusInternalServerError
		errorInfo = ErrorInfo{
			Code:    http.StatusText(http.StatusInternalServerError),
			Message: "Internal Server Error",
			Details: err.Error(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: false,
		Error:   &errorInfo,
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.Error(pkgErrors.NewJsonMarshalError("").Error())
	}
	h.Logger.Error("HTTP error", "status", status, "error", err)
}

func (h *Handler) ParseJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New("invalid JSON format")
	}
	return nil
}
