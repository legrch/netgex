package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/legrch/netgex/server"
	"github.com/legrch/netgex/service"
)

func main() {
	// Create context that listens for signals
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create a custom logger (optional)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Method 1: Configuration via options
	srv := server.NewServer(
		// Basic configuration
		server.WithLogger(logger),

		// Enable telemetry (required to enable observability)
		server.WithTelemetry(),

		// Configure specific backends
		server.WithTracingBackend("otlp", "otel-collector:4318"),
		server.WithMetricsBackend("prometheus", ""),
		server.WithProfilingBackend("pprof", ""),
	)

	// Method 2: Configuration via config object
	/*
		cfg := config.NewConfig()

		// Configure tracing
		cfg.Telemetry.Tracing.Enabled = true
		cfg.Telemetry.Tracing.Backend = "otlp"
		cfg.Telemetry.Tracing.Endpoint = "tempo-us-central1.grafana.net:4318"
		cfg.Telemetry.Tracing.Insecure = false
		cfg.Telemetry.Tracing.SampleRate = 0.1 // 10% sampling in production

		// Configure metrics
		cfg.Telemetry.Metrics.Enabled = true
		cfg.Telemetry.Metrics.Backend = "prometheus"
		cfg.Telemetry.Metrics.Path = "/metrics"

		// Configure profiling
		cfg.Telemetry.Profiling.Enabled = true
		cfg.Telemetry.Profiling.Backend = "pyroscope"
		cfg.Telemetry.Profiling.Endpoint = "http://pyroscope:4040"

		srv := server.NewServer(
			server.WithConfig(cfg),
			server.WithTelemetry(),
		)
	*/

	// Method 3: Configuration via environment variables
	/*
		// These would be set in the environment:
		// export TRACING_ENABLED=true
		// export TRACING_BACKEND=otlp
		// export TRACING_ENDPOINT=otel-collector:4318
		// export METRICS_ENABLED=true
		// export METRICS_BACKEND=prometheus
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

// Example for adding a service with traces and metrics
type MyServiceImpl struct {
	// ... your fields
}

func (s *MyServiceImpl) Register(server service.Registrar) {
	// Register your service with the server
}

// Example method with tracing/metrics
func (s *MyServiceImpl) SomeRPCMethod(ctx context.Context, req interface{}) (interface{}, error) {
	// Tracing is handled automatically through gRPC interceptors
	// Additional context or spans can be added:
	/*
		tracer := otel.Tracer("my-service")
		ctx, span := tracer.Start(ctx, "SomeRPCMethod.Process")
		defer span.End()

		// Add attributes to the span
		span.SetAttributes(attribute.String("user.id", "123"))
	*/

	// Do your work...

	// Metrics can be manually recorded:
	/*
		meter := otel.Meter("my-service")
		counter, _ := meter.Int64Counter("rpc.calls.count")
		counter.Add(ctx, 1)
	*/

	return nil, nil
}
