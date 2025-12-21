package ingestion

import (
	"context"
	"net/http"
)

type Server struct {
	server  *http.Server
	mux     *http.ServeMux
	handler *Handler
}

func NewServer(cfg *Config, handler *Handler) *Server {
	mux := http.NewServeMux()

	return &Server{
		handler: handler,
		mux:     mux,
		server: &http.Server{
			Addr:    cfg.BindAddress,
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	s.registerRoutes()
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	s.mux.HandleFunc("/ws", s.handler.HandleWS)
}
