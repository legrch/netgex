package telemetry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// setupMetrics configures metrics collection based on the provided configuration
func (s *Service) setupMetrics(ctx context.Context) error {
	cfg := s.config.Telemetry.Metrics

	if !cfg.Enabled {
		s.logger.Info("metrics collection is disabled")
		return nil
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(s.config.ServiceName),
			semconv.ServiceVersion(s.config.ServiceVersion),
			attribute.String("environment", s.config.Environment),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	switch cfg.Backend {
	case "prometheus":
		// Register HTTP handler for Prometheus metrics
		http.Handle(cfg.Path, promhttp.Handler())
		s.logger.Info("initialized Prometheus metrics", "path", cfg.Path)

		// We don't need a meter provider to clean up in this case
		return nil

	case "otlp":
		// Create OTLP metrics exporter
		opts := []otlpmetrichttp.Option{
			otlpmetrichttp.WithEndpoint(cfg.Endpoint),
		}

		if cfg.Insecure {
			opts = append(opts, otlpmetrichttp.WithInsecure())
		}

		exp, err := otlpmetrichttp.New(ctx, opts...)
		if err != nil {
			return fmt.Errorf("failed to create OTLP metric exporter: %w", err)
		}

		// Create MeterProvider
		mp := metric.NewMeterProvider(
			metric.WithReader(metric.NewPeriodicReader(exp)),
			metric.WithResource(res),
		)

		// Set global MeterProvider
		otel.SetMeterProvider(mp)
		s.meter = mp
		s.logger.Info("initialized OTLP metrics exporter", "endpoint", cfg.Endpoint)

	default:
		return fmt.Errorf("unsupported metrics backend: %s", cfg.Backend)
	}

	s.logger.Info("metrics initialized successfully", "backend", cfg.Backend)
	return nil
}

// RegisterMetrics registers common application metrics
func (s *Service) RegisterMetrics() {
	// Only register when using Prometheus
	if s.config.Telemetry.Metrics.Backend != "prometheus" {
		return
	}

	// Create a custom registry
	registry := prometheus.NewRegistry()

	// Example metric definitions
	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: s.config.Telemetry.Metrics.Namespace,
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: s.config.Telemetry.Metrics.Namespace,
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"method", "path"},
	)

	grpcRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: s.config.Telemetry.Metrics.Namespace,
			Name:      "grpc_requests_total",
			Help:      "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	grpcRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: s.config.Telemetry.Metrics.Namespace,
			Name:      "grpc_request_duration_seconds",
			Help:      "Duration of gRPC requests in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"method"},
	)

	// Register with registry
	registry.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		grpcRequestsTotal,
		grpcRequestDuration,
	)

	// Use our registry as the default prometheus registry
	http.Handle(s.config.Telemetry.Metrics.Path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
}
