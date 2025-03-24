package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/pyroscope-io/client/pyroscope"
)

// setupProfiling configures continuous profiling based on the provided configuration
func (s *Service) setupProfiling(ctx context.Context) error {
	cfg := s.config.Telemetry.Profiling

	if !cfg.Enabled {
		s.logger.Info("continuous profiling is disabled")
		return nil
	}

	switch cfg.Backend {
	case "pyroscope", "phlare":
		// Configure Pyroscope profiler
		profileTypes := parseProfileTypes(cfg.Types)

		profiler, err := pyroscope.Start(pyroscope.Config{
			ApplicationName: s.config.ServiceName,
			ServerAddress:   cfg.Endpoint,
			Logger:          newPyroscopeLogger(s.logger),
			ProfileTypes:    profileTypes,
			Tags: map[string]string{
				"version":     s.config.ServiceVersion,
				"environment": s.config.Environment,
			},
			SampleRate: uint32(cfg.SampleRate * 100), // Convert to sampling frequency
		})

		if err != nil {
			return fmt.Errorf("failed to start pyroscope profiler: %w", err)
		}

		s.profiler = profiler
		s.logger.Info("pyroscope profiler initialized",
			"endpoint", cfg.Endpoint,
			"types", cfg.Types)

	case "otlp":
		s.logger.Info("OTLP profiling is not fully supported yet")
		// Future: implement OTLP profiles when standard is more mature
		return nil

	default:
		return fmt.Errorf("unsupported profiling backend: %s", cfg.Backend)
	}

	s.logger.Info("continuous profiling initialized successfully", "backend", cfg.Backend)
	return nil
}

// newPyroscopeLogger creates a logger adapter for Pyroscope that uses slog
func newPyroscopeLogger(logger *slog.Logger) pyroscope.Logger {
	return &pyroscopeLoggerAdapter{logger: logger}
}

// pyroscopeLoggerAdapter adapts slog for use with Pyroscope
type pyroscopeLoggerAdapter struct {
	logger *slog.Logger
}

func (l *pyroscopeLoggerAdapter) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *pyroscopeLoggerAdapter) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *pyroscopeLoggerAdapter) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// parseProfileTypes converts a comma-separated string of profile types to Pyroscope types
func parseProfileTypes(types string) []pyroscope.ProfileType {
	var profileTypes []pyroscope.ProfileType

	typeMap := map[string]pyroscope.ProfileType{
		"cpu":       pyroscope.ProfileCPU,
		"heap":      pyroscope.ProfileAllocObjects,
		"alloc":     pyroscope.ProfileAllocSpace,
		"goroutine": pyroscope.ProfileGoroutines,
		"mutex":     pyroscope.ProfileMutexCount,
		"block":     pyroscope.ProfileBlockCount,
	}

	for _, t := range strings.Split(types, ",") {
		if pt, ok := typeMap[strings.TrimSpace(t)]; ok {
			profileTypes = append(profileTypes, pt)
		}
	}

	// Default to CPU if no valid types
	if len(profileTypes) == 0 {
		profileTypes = append(profileTypes, pyroscope.ProfileCPU)
	}

	return profileTypes
}
