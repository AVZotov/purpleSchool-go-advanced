package base

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	pkgErrors "order_api_cart/pkg/errors"
	pkgLogger "order_api_cart/pkg/logger"
)

type Handler struct{}

func (h *Handler) WriteJSON(ctx context.Context, w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: status < 400,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		pkgLogger.ErrorWithRequestID(ctx, pkgErrors.ErrEncodingJSON.Error(), logrus.Fields{
			"error": err.Error(),
		})
		http.Error(w, pkgErrors.ErrEncodingJSON.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) WriteError(ctx context.Context, w http.ResponseWriter, statusCode int, err error) {
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
		pkgLogger.ErrorWithRequestID(ctx, pkgErrors.ErrEncodingJSON.Error(), logrus.Fields{
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ParseJSON(ctx context.Context, r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		pkgLogger.ErrorWithRequestID(ctx, pkgErrors.ErrDecodingJSON.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return fmt.Errorf("%w %v", pkgErrors.ErrDecodingJSON, err)
	}

	return nil
}
