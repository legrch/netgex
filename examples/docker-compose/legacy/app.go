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
	"github.com/legrch/netgex/service"
)

func main() {
	// Create a context that listens for termination signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("Starting legacy example app with Jaeger and Prometheus...")

	// Create basic configuration
	cfg := config.NewConfig()

	// Configure service info
	cfg.ServiceName = getEnv("SERVICE_NAME", "netgex-legacy-example")
	cfg.ServiceVersion = getEnv("SERVICE_VERSION", "1.0.0")
	cfg.Environment = getEnv("ENVIRONMENT", "development")

	// Configure server addresses - using different ports to avoid conflicts
	cfg.GRPCAddress = getEnv("GRPC_ADDRESS", ":50052")      // Changed from :9090
	cfg.HTTPAddress = getEnv("HTTP_ADDRESS", ":8082")       // Changed from :8080
	cfg.MetricsAddress = getEnv("METRICS_ADDRESS", ":9093") // Changed from :9091
	cfg.PprofAddress = getEnv("PPROF_ADDRESS", ":6062")     // Changed from :6060

	// Configure telemetry - Legacy setup
	// Tracing with Jaeger - use localhost instead of container names
	cfg.Telemetry.Tracing.Enabled = true
	cfg.Telemetry.Tracing.Backend = "jaeger"
	cfg.Telemetry.Tracing.Endpoint = getEnv("TRACING_ENDPOINT", "http://localhost:14268/api/traces")
	cfg.Telemetry.Tracing.SampleRate = 1.0 // 100% sampling for demo

	// Metrics with Prometheus
	cfg.Telemetry.Metrics.Enabled = true
	cfg.Telemetry.Metrics.Backend = "prometheus"
	cfg.Telemetry.Metrics.Path = "/metrics"

	// Profiling with pprof
	cfg.Telemetry.Profiling.Enabled = true
	cfg.Telemetry.Profiling.Backend = "pprof"

	// Create the server with telemetry enabled
	srv := server.NewServer(
		server.WithConfig(cfg),
		server.WithTelemetry(),
	)

	// // Add our example service (in a real app)
	// exampleService := &ExampleService{}
	// srv.WithServices(exampleService)

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
	log.Printf("Pprof endpoint: http://localhost%s/debug/pprof", cfg.PprofAddress)
	log.Println("Jaeger UI: http://localhost:16686")
	log.Println("Prometheus UI: http://localhost:9090")
	log.Println("Grafana: http://localhost:3000")

	// Wait for termination signal
	<-ctx.Done()
	log.Println("Shutting down...")
}

// ExampleService demonstrates service registration
type ExampleService struct{}

// Register registers the service with the server
func (s *ExampleService) Register(server service.Registrar) {
	log.Println("Registering example service")
	// In a real application, you would register your gRPC service here
}

func generateLoad() {
	// Simple function to generate some CPU load and randomized behavior
	for {
		// Random delay to simulate variable processing time
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		// Do some meaningless calculations to generate CPU load
		x := 0
		for i := 0; i < 1000000; i++ {
			x += i
		}
	}
}

// Helper to get environment variables with default fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
