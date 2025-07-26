package server

import (
	"net/http"
)

type Server struct {
	Port    string
	Handler http.Handler
}

func New(port string, router http.Handler) *Server {
	return &Server{
		Port:    ":" + port,
		Handler: router,
	}
}

func (s *Server) ListenAndServe() error {
	server := &http.Server{
		Addr:    s.Port,
		Handler: s.Handler,
	}
	return server.ListenAndServe()
}
