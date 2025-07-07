package router

import (
	"net/http"
)

type Config interface {
	GetEmailSecrets() *map[string]string
}

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	return router
}
