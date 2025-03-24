package telemetry

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/legrch/netgex/config"
)

// Service represents the telemetry service which handles tracing, metrics, logging, and profiling
type Service struct {
	logger         *slog.Logger
	config         *config.Config
	tracerProvider interface{ Shutdown(context.Context) error }
	meterProvider  interface{ Shutdown(context.Context) error }
	profiler       interface{ Stop() error }
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

	// Initialize each component based on configuration
	if err := s.setupTracing(ctx); err != nil {
		return fmt.Errorf("failed to set up tracing: %w", err)
	}

	if err := s.setupMetrics(ctx); err != nil {
		return fmt.Errorf("failed to set up metrics: %w", err)
	}

	if err := s.setupLogging(ctx); err != nil {
		return fmt.Errorf("failed to set up logging: %w", err)
	}

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
	if s.tracerProvider != nil {
		if err := s.tracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("trace provider shutdown: %w", err))
		}
	}

	// Shutdown metrics
	if s.meterProvider != nil {
		if err := s.meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("meter provider shutdown: %w", err))
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
