package pkg

import (
	"net/http"
)

type Config interface {
	GetGmailSecrets() *map[string]string
}

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	return router
}
