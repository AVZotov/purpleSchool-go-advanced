package system

import (
	resp "link_shortener/internal/http-server/handlers/types"
	"net/http"
	"time"
)

func NewHealthCheckHandler(router *http.ServeMux) {
	router.HandleFunc("GET /health", health())
	router.HandleFunc("GET /api/v1/health", health())
}

func health() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status":    "OK",
			"service":   "link_shortener",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.2.1",
		}

		resp.Json(w, http.StatusOK, response)
	}
}
