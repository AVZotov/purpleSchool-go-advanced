package base

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	pkgLogger "order_api_auth/pkg/logger"
)

type Handler struct{}

type ErrorInfo struct {
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

func (h *Handler) WriteJSON(r *http.Request, w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: status < 400,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		pkgLogger.ErrorWithRequestID(r, "error encoding json", logrus.Fields{
			"error": err.Error(),
		})
		http.Error(w, "error encoding json", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) WriteError(r *http.Request, w http.ResponseWriter, statusCode int, err error) {
	var errorInfo ErrorInfo
	status := statusCode
	errorInfo = ErrorInfo{
		Code:    http.StatusText(statusCode),
		Details: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: false,
		Error:   &errorInfo,
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		pkgLogger.ErrorWithRequestID(r, "error encoding json", logrus.Fields{
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ParseJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		pkgLogger.ErrorWithRequestID(r, "error decoding json", logrus.Fields{
			"error": err.Error(),
		})
		return fmt.Errorf("could not parse request body: %w", err)
	}

	return nil
}
