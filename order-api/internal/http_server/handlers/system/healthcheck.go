package system

import (
	"net/http"
	"order/internal/http_server/handlers/base"
	"time"
)

const (
	HealthV1 = "/api/v1/health"
)

type Handler struct {
	base.Handler
}

func New(mux *http.ServeMux) {
	handler := &Handler{
		Handler: base.Handler{},
	}

	handler.registerRoutes(mux)
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET "+HealthV1, h.health)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	response := map[string]interface{}{
		"status":    "OK",
		"service":   "order-api",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	h.WriteJSON(w, http.StatusOK, response)
}
