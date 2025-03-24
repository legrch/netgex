package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/legrch/netgex/config"
	"github.com/legrch/netgex/server"
)

func main() {
	// Create a context that listens for termination signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("Starting OTLP local example app with OpenTelemetry, Tempo, Prometheus and Pyroscope...")

	// Create basic configuration
	cfg := config.NewConfig()

	// Configure service info
	cfg.ServiceName = getEnv("SERVICE_NAME", "netgex-otlp-example")
	cfg.ServiceVersion = getEnv("SERVICE_VERSION", "1.0.0")
	cfg.Environment = getEnv("ENVIRONMENT", "development")

	// Configure server addresses
	cfg.GRPCAddress = getEnv("GRPC_ADDRESS", ":9090")
	cfg.HTTPAddress = getEnv("HTTP_ADDRESS", ":8080")
	cfg.MetricsAddress = getEnv("METRICS_ADDRESS", ":9091")

	// Configure telemetry - OTLP setup
	// Tracing with OTLP
	cfg.Telemetry.Tracing.Enabled = true
	cfg.Telemetry.Tracing.Backend = "otlp"
	cfg.Telemetry.Tracing.Endpoint = getEnv("TRACING_ENDPOINT", "otel-collector:4318")
	cfg.Telemetry.Tracing.SampleRate = 1.0 // 100% sampling for demo

	// Metrics with Prometheus
	cfg.Telemetry.Metrics.Enabled = true
	cfg.Telemetry.Metrics.Backend = "prometheus"
	cfg.Telemetry.Metrics.Path = "/metrics"

	// Profiling with Pyroscope
	cfg.Telemetry.Profiling.Enabled = true
	cfg.Telemetry.Profiling.Backend = "pyroscope"
	cfg.Telemetry.Profiling.Endpoint = getEnv("PROFILING_ENDPOINT", "http://pyroscope:4040")

	// Create the server with telemetry enabled
	srv := server.NewServer(
		server.WithConfig(cfg),
		server.WithTelemetry(),
	)

	// Start the server in a goroutine
	go func() {
		log.Println("Starting server...")
		if err := srv.Run(ctx); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Simulate some load to generate metrics/traces/profiles
	go generateLoad()

	log.Println("Server started successfully")
	log.Printf("HTTP endpoint: http://localhost%s", cfg.HTTPAddress)
	log.Printf("gRPC endpoint: localhost%s", cfg.GRPCAddress)
	log.Printf("Metrics endpoint: http://localhost%s/metrics", cfg.MetricsAddress)
	log.Println("Grafana: http://localhost:3000")
	log.Println("Prometheus: http://localhost:9090")
	log.Println("Pyroscope: http://localhost:4040")
	log.Println("Tempo (via Grafana): http://localhost:3000/explore (select Tempo)")

	// Wait for termination signal
	<-ctx.Done()
	log.Println("Shutting down...")
}

func generateLoad() {
	// Generate varied load patterns
	for {
		// Do some CPU-intensive work to generate metrics
		generateCPULoad()

		// Simulate slow operations for tracing
		simulateSlowOperations()

		// Random delay between operations
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	}
}

func generateCPULoad() {
	// Do CPU-intensive calculations
	x := 0
	for i := 0; i < 1000000; i++ {
		x += i * i
	}
}

func simulateSlowOperations() {
	// Simulate database query
	time.Sleep(time.Duration(10+rand.Intn(30)) * time.Millisecond)

	// Simulate external API call
	if rand.Intn(10) > 7 { // 30% chance of slow request
		time.Sleep(time.Duration(100+rand.Intn(150)) * time.Millisecond)
	} else {
		time.Sleep(time.Duration(20+rand.Intn(40)) * time.Millisecond)
	}
}

// Helper to get environment variables with default fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
