package system

import (
	"link_shortener/pkg/resp"
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
			"version":   "0.2.0",
		}

		resp.Json(w, http.StatusOK, response)
	}
}
