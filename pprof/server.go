package pprof

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	// #nosec G108 - pprof endpoints are intentionally exposed for debugging
	_ "net/http/pprof" // Register pprof handlers
	"time"
)

// Server represents a server for exposing pprof profiling endpoints
type Server struct {
	logger *slog.Logger
	server *http.Server
}

// NewServer creates a new pprof server
func NewServer(logger *slog.Logger, address string) *Server {
	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:              address,
			Handler:           http.DefaultServeMux, // DefaultServeMux has pprof handlers registered
			ReadHeaderTimeout: 5 * time.Second,      // Prevent Slowloris attacks
		},
	}
}

// PreRun prepares the pprof server
func (*Server) PreRun(_ context.Context) error {
	return nil
}

// Run starts the pprof server
func (p *Server) Run(_ context.Context) error {
	p.logger.Info("starting pprof server", "address", p.server.Addr)
	if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("pprof server error: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the pprof server
func (p *Server) Shutdown(ctx context.Context) error {
	p.logger.Info("shutting down pprof server")
	if err := p.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("pprof server shutdown error: %w", err)
	}
	return nil
}
