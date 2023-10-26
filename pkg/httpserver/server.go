package httpserver

import (
	"context"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(handler http.Handler, opts ...Options) *Server {
	httpServer := &http.Server{
		Handler: handler,
	}

	server := &Server{
		server: httpServer,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Startup() error {
	return s.server.ListenAndServe()
}
