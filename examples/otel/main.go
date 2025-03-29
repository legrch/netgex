package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/legrch/netgex/server"
)

func main() {
	// Create context that listens for signals
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create a custom logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Method 1: Configuration via options (most concise)
	srv := server.NewServer(
		// Basic configuration
		server.WithLogger(logger),

		// Enable telemetry with OpenTelemetry unified approach
		server.WithTelemetry(),
		server.WithOTEL("otel-collector:4318", true), // endpoint, insecure

		// Configure profiling separately
		server.WithProfilingBackend("pprof", ""),
	)

	// Method 2: Configuration via config object
	/*
		cfg := config.NewConfig()

		// Configure service info
		cfg.ServiceName = "my-service"
		cfg.ServiceVersion = "1.0.0"
		cfg.Environment = "production"

		// Configure unified OpenTelemetry
		cfg.Telemetry.OTEL.Enabled = true
		cfg.Telemetry.OTEL.Endpoint = "otel-collector:4318"
		cfg.Telemetry.OTEL.Insecure = true
		cfg.Telemetry.OTEL.Protocol = "http"
		cfg.Telemetry.OTEL.Headers = "Authorization=Basic YOUR_API_KEY"
		cfg.Telemetry.OTEL.TracesEnabled = true
		cfg.Telemetry.OTEL.MetricsEnabled = true
		cfg.Telemetry.OTEL.SampleRate = 0.1 // 10% sampling in production

		// Configure profiling separately
		cfg.Telemetry.Profiling.Enabled = true
		cfg.Telemetry.Profiling.Backend = "pprof"

		srv := server.NewServer(
			server.WithConfig(cfg),
			server.WithTelemetry(),
		)
	*/

	// Method 3: Configuration via environment variables
	/*
		// These would be set in the environment:
		// export OTEL_ENABLED=true
		// export OTEL_ENDPOINT=otel-collector:4318
		// export OTEL_INSECURE=true
		// export OTEL_PROTOCOL=http
		// export OTEL_HEADERS="Authorization=Basic YOUR_API_KEY"
		// export OTEL_TRACES_ENABLED=true
		// export OTEL_METRICS_ENABLED=true
		// export OTEL_SAMPLE_RATE=0.1
		// export PROFILING_ENABLED=true
		// export PROFILING_BACKEND=pprof

		// Load config from environment
		cfg, err := config.LoadFromEnv("")
		if err != nil {
			logger.Error("failed to load config", "error", err)
			os.Exit(1)
		}

		srv := server.NewServer(
			server.WithConfig(cfg),
			server.WithTelemetry(),
		)
	*/

	// Register your service implementations
	// srv.AddService(&myservice.ServiceImpl{})

	// Run the server
	if err := srv.Run(ctx); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
