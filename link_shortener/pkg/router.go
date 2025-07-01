package pkg

import (
	"net/http"
)

type Config interface {
	GetEmailConfig() map[string]string
}

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	return router
}
