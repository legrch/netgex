package main

import (
	"context"
	"fmt"
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

	log.Println("Starting Grafana Cloud example app...")

	// Create basic configuration
	cfg := config.NewConfig()

	// Configure service info
	cfg.ServiceName = getEnv("SERVICE_NAME", "netgex-cloud-example")
	cfg.ServiceVersion = getEnv("SERVICE_VERSION", "1.0.0")
	cfg.Environment = getEnv("ENVIRONMENT", "production")

	// Configure server addresses - using different ports to avoid conflicts
	cfg.GRPCAddress = getEnv("GRPC_ADDRESS", ":50053")      // Changed from :9090
	cfg.HTTPAddress = getEnv("HTTP_ADDRESS", ":8083")       // Changed from :8080
	cfg.MetricsAddress = getEnv("METRICS_ADDRESS", ":9094") // Changed from :9091
	cfg.PprofAddress = getEnv("PPROF_ADDRESS", ":6063")     // Changed from :6060

	// Configure telemetry - Grafana Cloud setup
	apiKey := getEnv("GRAFANA_CLOUD_API_KEY", "")
	if apiKey == "" {
		log.Println("Warning: GRAFANA_CLOUD_API_KEY environment variable not set")
		log.Println("Telemetry data will not be sent to Grafana Cloud")
	}

	// Tracing with OTLP to Grafana Cloud
	cfg.Telemetry.Tracing.Enabled = true
	cfg.Telemetry.Tracing.Backend = "otlp"
	cfg.Telemetry.Tracing.Endpoint = getEnv("TRACING_ENDPOINT", "tempo-us-central1.grafana.net:4318")
	cfg.Telemetry.Tracing.Insecure = false
	cfg.Telemetry.Tracing.SampleRate = 0.1 // 10% sampling for production

	// We need to set authorization in environment variables for OTLP
	os.Setenv("OTEL_EXPORTER_OTLP_HEADERS", fmt.Sprintf("Authorization=Basic %s", apiKey))

	// Metrics with OTLP to Grafana Cloud
	cfg.Telemetry.Metrics.Enabled = true
	cfg.Telemetry.Metrics.Backend = "otlp"
	cfg.Telemetry.Metrics.Endpoint = getEnv("METRICS_ENDPOINT", "prometheus-us-central1.grafana.net:4318")
	cfg.Telemetry.Metrics.Insecure = false

	// The same environment variable is used for both traces and metrics

	// Profiling with Pyroscope to Grafana Cloud
	cfg.Telemetry.Profiling.Enabled = true
	cfg.Telemetry.Profiling.Backend = "pyroscope"
	cfg.Telemetry.Profiling.Endpoint = getEnv("PROFILING_ENDPOINT", "profiles-prod-006.grafana.net:4040")

	// Pyroscope auth token is added at runtime by the telemetry service

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

	// Simulate workload
	go simulateProduction()

	log.Println("Server started successfully")
	log.Printf("HTTP endpoint: http://localhost%s", cfg.HTTPAddress)
	log.Printf("gRPC endpoint: localhost%s", cfg.GRPCAddress)
	log.Printf("Metrics endpoint: http://localhost%s/metrics", cfg.MetricsAddress)
	log.Printf("pprof endpoint: http://localhost%s/debug/pprof", cfg.PprofAddress)
	log.Println("Grafana Cloud data should be visible at your Grafana Cloud instance")
	log.Println("Local Grafana: http://localhost:3000")

	// Wait for termination signal
	<-ctx.Done()
	log.Println("Shutting down...")
}

func simulateProduction() {
	// Simulate production workload
	for {
		// Simulate varied traffic patterns
		simulateWebTraffic()

		// Simulate background jobs
		simulateBackgroundJobs()

		// Create periodic CPU spikes for profiling
		if rand.Intn(100) < 5 { // 5% chance of CPU spike
			generateCPUSpike()
		}

		time.Sleep(time.Duration(10+rand.Intn(40)) * time.Millisecond)
	}
}

func simulateWebTraffic() {
	// Simulate web requests with varying latency
	// 80% fast requests, 15% medium, 5% slow
	r := rand.Intn(100)
	switch {
	case r < 80:
		time.Sleep(time.Duration(5+rand.Intn(20)) * time.Millisecond)
	case r < 95:
		time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
	default:
		time.Sleep(time.Duration(200+rand.Intn(300)) * time.Millisecond)
	}
}

func simulateBackgroundJobs() {
	// Simulate periodic background jobs
	if rand.Intn(100) < 10 { // 10% chance to run background job
		time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
	}
}

func generateCPUSpike() {
	// Create CPU spike for profiling
	x := 0
	for i := 0; i < 5000000; i++ {
		x += i * i
	}
}

// Helper to get environment variables with default fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
