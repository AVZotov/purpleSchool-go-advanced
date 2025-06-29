package api

import (
	"net/http"
)

func newRouter() *http.ServeMux {
	router := http.NewServeMux()
	registerHandlers(router)
	return router
}

func registerHandlers(router *http.ServeMux) {
	newRandHandler(router)
}
