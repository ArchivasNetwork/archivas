package metrics

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server provides a lightweight HTTP server for Prometheus metrics
type Server struct {
	srv *http.Server
}

// NewServer creates a new metrics server
func NewServer(addr string) *Server {
	if addr == "" {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return &Server{
		srv: &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

// Start starts the metrics server in a goroutine
func (s *Server) Start() {
	if s == nil {
		return
	}
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("metrics server error: %v", err)
		}
	}()
	log.Printf("📊 Metrics server listening on %s/metrics", s.srv.Addr)
}

// Stop gracefully shuts down the metrics server
func (s *Server) Stop(ctx context.Context) error {
	if s == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}
