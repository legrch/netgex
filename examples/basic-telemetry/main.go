package main

import (
	"context"
	"log"
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

	// Create basic configuration
	cfg := config.NewConfig()

	// Enable pprof for basic profiling (doesn't require external services)
	cfg.Telemetry.Profiling.Enabled = true
	cfg.Telemetry.Profiling.Backend = "pprof"

	// Enable Prometheus metrics
	cfg.Telemetry.Metrics.Enabled = true
	cfg.Telemetry.Metrics.Backend = "prometheus"

	// In a production environment, you might want to enable tracing too:
	// cfg.Telemetry.Tracing.Enabled = true
	// cfg.Telemetry.Tracing.Backend = "jaeger"
	// cfg.Telemetry.Tracing.Endpoint = "http://localhost:14268/api/traces"

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

	// Simulate some load to generate metrics/profiles
	go generateLoad()

	// Wait for termination signal
	<-ctx.Done()
	log.Println("Shutting down...")
}

func generateLoad() {
	// Simple function to generate some CPU load for profiling
	for {
		// Do some meaningless calculations to generate CPU load
		x := 0
		for i := 0; i < 1000000; i++ {
			x += i
		}
		time.Sleep(100 * time.Millisecond)
	}
}
