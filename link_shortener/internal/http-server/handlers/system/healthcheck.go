package system

import (
	"link_shortener"
	"link_shortener/internal/http-server/handlers/base"
	"link_shortener/pkg/logger"
	"net/http"
	"time"
)

const (
	HealthV1 = "/api/v1/health"
	Health   = "/health"
)

type Handler struct {
	base.Handler
}

func New(mux *http.ServeMux, logger logger.Logger) {
	handler := &Handler{
		Handler: base.Handler{Logger: logger},
	}

	handler.registerRoutes(mux)
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET "+HealthV1, h.health)
	mux.HandleFunc("GET "+Health, h.health)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	response := map[string]interface{}{
		"status":    "OK",
		"service":   link_shortener.AppName,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   link_shortener.Version,
		"buildDate": link_shortener.BuildDate,
	}
	h.WriteJSON(w, http.StatusOK, response)
}
