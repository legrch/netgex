package telemetry

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/legrch/netgex/config"
)

// Service represents the telemetry service which handles tracing, metrics, logging, and profiling
type Service struct {
	logger *slog.Logger
	config *config.Config
	// tracer is `otlp.TracerProvider`, `jaeger.Tracer`, or none
	tracer interface{ Shutdown(context.Context) error }
	// meter is `otlp.MeterProvider`, or none
	meter interface{ Shutdown(context.Context) error }
	// profiler is `pyroscope.Profiler`, or none
	profiler interface{ Stop() error }
	// otelProvider is the unified OpenTelemetry provider if enabled
	otelProvider interface{ Shutdown(context.Context) error }
}

// NewService creates a new telemetry service
func NewService(logger *slog.Logger, config *config.Config) *Service {
	return &Service{
		logger: logger,
		config: config,
	}
}

// PreRun sets up telemetry before the server starts
func (s *Service) PreRun(ctx context.Context) error {
	s.logger.Info("initializing telemetry services")

	// Initialize logging first for better diagnostics
	if err := s.setupLogging(ctx); err != nil {
		return fmt.Errorf("failed to set up logging: %w", err)
	}

	// Check if OpenTelemetry unified configuration is enabled
	if s.config.Telemetry.OTEL.Enabled {
		// If OTEL is enabled, use it as the primary provider
		if err := s.setupOTEL(ctx); err != nil {
			return fmt.Errorf("failed to set up OpenTelemetry: %w", err)
		}
	} else {
		// Otherwise, initialize separate components based on configuration
		// Legacy tracing setup
		if err := s.setupTracing(ctx); err != nil {
			return fmt.Errorf("failed to set up tracing: %w", err)
		}

		// Legacy metrics setup
		if err := s.setupMetrics(ctx); err != nil {
			return fmt.Errorf("failed to set up metrics: %w", err)
		}
	}

	// Profiling is always set up separately
	if err := s.setupProfiling(ctx); err != nil {
		return fmt.Errorf("failed to set up profiling: %w", err)
	}

	return nil
}

// Run is a no-op for telemetry service as it's passive
func (s *Service) Run(ctx context.Context) error {
	// Telemetry Service doesn't need to run actively
	<-ctx.Done()
	return nil
}

// Shutdown gracefully terminates telemetry components
func (s *Service) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down telemetry services")

	var errs []error

	// Shutdown tracing
	if s.tracer != nil {
		if err := s.tracer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("trace provider shutdown: %w", err))
		}
	}

	// Shutdown metrics
	if s.meter != nil {
		if err := s.meter.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("meter provider shutdown: %w", err))
		}
	}

	// Shutdown unified OTEL provider if exists
	if s.otelProvider != nil {
		if err := s.otelProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("OTEL provider shutdown: %w", err))
		}
	}

	// Shutdown profiler
	if s.profiler != nil {
		if err := s.profiler.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("profiler shutdown: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("telemetry shutdown errors: %v", errs)
	}

	return nil
}
