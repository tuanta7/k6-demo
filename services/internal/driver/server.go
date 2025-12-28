package driver

import (
	"context"
	"net/http"

	"github.com/tuanta7/k6noz/services/pkg/otelx"
)

type Server struct {
	mux        *http.ServeMux
	server     *http.Server
	handler    *Handler
	prometheus *otelx.PrometheusProvider
}

func NewServer(addr string, handler *Handler, prometheus *otelx.PrometheusProvider) *Server {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &Server{
		mux:        mux,
		server:     server,
		handler:    handler,
		prometheus: prometheus,
	}
}

func (s *Server) Run() error {
	s.mux.Handle("GET /metrics", s.prometheus.Handler())
	s.mux.Handle("GET /drivers/{id}", otelx.Handler(s.handler.GetDriverByID, "GetDriverByID"))
	s.mux.Handle("POST /ratings", otelx.Handler(s.handler.CreateNewRating, "CreateNewRating"))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
