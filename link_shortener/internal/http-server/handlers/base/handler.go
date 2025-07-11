package base

import (
	"encoding/json"
	"link_shortener/pkg/errors"
	"link_shortener/pkg/logger"
	"net/http"
)

type Handler struct {
	Logger logger.Logger
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (h *Handler) WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: status < 400,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.Error(errors.NewJsonParseError("").Error())
	}
}

func (h *Handler) WriteError(w http.ResponseWriter, err error) {
	var appErr errors.AppError
	var status int
	var errorInfo ErrorInfo

	if errors.AsAppError(err, &appErr) {
		status = appErr.Status
		errorInfo = ErrorInfo{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		}
	} else {
		status = http.StatusInternalServerError
		errorInfo = ErrorInfo{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
		}
	}

	h.Logger.Debug("app error struct:", appErr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: false,
		Error:   &errorInfo,
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		h.Logger.Error("Failed to encode error response", "error", encodeErr)
	}

	h.Logger.Error("HTTP error", "status", status, "error", err)
}

func (h *Handler) ParseJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.NewJsonParseError("Invalid JSON format")
	}
	return nil
}
