package system

import "net/http"

func NewHealthCheckHandler(router *http.ServeMux) {
	router.HandleFunc("GET /health", health())
}

func health() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
