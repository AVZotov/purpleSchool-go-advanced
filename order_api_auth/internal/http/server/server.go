package server

import (
	"net/http"
	"order/pkg/middleware"
)

type Server struct {
	Port    string
	Handler http.Handler
}

func New(port string, router *http.ServeMux) *Server {
	return &Server{
		Port:    ":" + port,
		Handler: middleware.Logger(router),
	}
}

func (s *Server) ListenAndServe() error {
	server := &http.Server{
		Addr:    s.Port,
		Handler: s.Handler,
	}
	return server.ListenAndServe()
}
