package pyroscope

import (
	"context"
	"log/slog"
)

// Profiler represents a profiler using Pyroscope
type Profiler struct {
	logger *slog.Logger
}

// NewProfiler creates a new Pyroscope profiler
func NewProfiler(logger *slog.Logger) *Profiler {
	return &Profiler{
		logger: logger,
	}
}

// PreRun prepares the Pyroscope profiler
func (p *Profiler) PreRun(_ context.Context) error {
	p.logger.Info("initializing pyroscope profiler")
	// In a real implementation, this would initialize the Pyroscope client
	return nil
}

// Run starts the Pyroscope profiler
func (*Profiler) Run(_ context.Context) error {
	// Pyroscope runs in the background, so this just returns nil
	return nil
}

// Shutdown stops the Pyroscope profiler
func (p *Profiler) Shutdown(_ context.Context) error {
	p.logger.Info("shutting down pyroscope profiler")
	// In a real implementation, this would stop the Pyroscope client
	return nil
}
