package telemetry

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// setupLogging configures structured logging based on the provided configuration
func (s *Service) setupLogging(ctx context.Context) error {
	cfg := s.config.Telemetry.Logging

	if !cfg.Enabled {
		s.logger.Info("structured logging is disabled")
		return nil
	}

	// If we already have a logger set by WithLogger, we respect it
	// but still add service attributes if not already present
	if s.logger != nil {
		// Add service context to logs if not already present
		s.logger = s.logger.With(
			"service", s.config.ServiceName,
			"version", s.config.ServiceVersion,
			"environment", s.config.Environment,
		)

		// Set as default logger if requested
		if cfg.Backend == "global" {
			slog.SetDefault(s.logger)
			s.logger.Info("setting logger as global default")
		}

		s.logger.Info("using existing logger provided via WithLogger")
		return nil
	}

	// Determine log level
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	var output io.Writer = os.Stdout

	// Configure file output if specified
	if cfg.FilePath != "" {
		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		output = file
	}

	// Create handler based on format
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{Level: level})
	case "text", "console":
		handler = slog.NewTextHandler(output, &slog.HandlerOptions{Level: level})
	default:
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{Level: level})
	}

	// Add service context to logs
	handler = handler.WithAttrs([]slog.Attr{
		slog.String("service", s.config.ServiceName),
		slog.String("version", s.config.ServiceVersion),
		slog.String("environment", s.config.Environment),
	})

	// Create logger
	logger := slog.New(handler)
	s.logger = logger

	// Set as global default logger if requested
	if cfg.Backend == "global" {
		slog.SetDefault(logger)
	}

	// If we're using a backend like OTLP, we'd configure it here
	if cfg.Backend == "otlp" && cfg.Endpoint != "" {
		// OTLP logging setup would go here when OpenTelemetry
		// fully supports the log bridge for Go
		s.logger.Info("OTLP logging is not fully supported yet in Go OTel SDK")
	}

	s.logger.Info("structured logging initialized successfully",
		"format", cfg.Format,
		"level", cfg.Level)
	return nil
}
