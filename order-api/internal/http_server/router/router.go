package router

import (
	"net/http"
)

func New() *http.ServeMux {
	router := http.NewServeMux()
	return router
}
