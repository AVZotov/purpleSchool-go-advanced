package api

import "net/http"

type Server struct {
	Addr   string
	Router *http.ServeMux
}

func NewServer(port string) *Server {
	return &Server{
		Addr:   ":" + port,
		Router: newRouter(),
	}
}

func (s *Server) ListenAndServe() error {
	server := &http.Server{
		Addr:    s.Addr,
		Handler: s.Router,
	}
	return server.ListenAndServe()
}
