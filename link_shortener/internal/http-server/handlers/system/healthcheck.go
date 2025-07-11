package system

import (
	"link_shortener"
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
			"service":   link_shortener.AppName,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   link_shortener.Version,
			"buildDate": link_shortener.BuildDate,
		}

		resp.Json(w, http.StatusOK, response)
	}
}
