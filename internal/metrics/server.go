package metrics

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server represents a server for exposing Prometheus metrics
type Server struct {
	logger       *slog.Logger
	server       *http.Server
	closeTimeout time.Duration
}

// NewServer creates a new metrics server
func NewServer(logger *slog.Logger, address string, closeTimeout time.Duration) *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{
		logger:       logger,
		server:       server,
		closeTimeout: closeTimeout,
	}
}

// PreRun prepares the metrics server
func (*Server) PreRun(_ context.Context) error {
	// Register application metrics
	RegisterAppMetrics()
	return nil
}

// Run starts the metrics server
func (m *Server) Run(_ context.Context) error {
	m.logger.Info("starting metrics server", "address", m.server.Addr)
	if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("metrics server error: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the metrics server
func (m *Server) Shutdown(ctx context.Context) error {
	m.logger.Info("shutting down metrics server")

	shutdownCtx, cancel := context.WithTimeout(ctx, m.closeTimeout)
	defer cancel()

	if err := m.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("metrics server shutdown error: %w", err)
	}

	return nil
}

// AppVersion is a gauge for tracking application version
var AppVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "app",
	Name:      "version",
	Help:      "Application version",
}, []string{"version"})

// RegisterAppMetrics registers application metrics with Prometheus
func RegisterAppMetrics() {
	prometheus.MustRegister(AppVersion)
}

// UnregisterAppMetrics unregisters application metrics from Prometheus
func UnregisterAppMetrics() {
	prometheus.Unregister(AppVersion)
}

// SetAppVersion sets the application version metric
func SetAppVersion(version string) {
	AppVersion.WithLabelValues(version).Set(1)
}
