package main

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create a context that listens for termination signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Setup profiling
	setupPprof()

	// Setup metrics
	setupMetrics()

	// Simulate some load to generate metrics/profiles
	go generateLoad()

	log.Println("Server is running...")
	log.Println("Metrics: http://localhost:9091/metrics")
	log.Println("Profiling: http://localhost:6060/debug/pprof/")

	// Wait for termination signal
	<-ctx.Done()
	log.Println("Shutting down...")
}

func setupPprof() {
	// Create pprof HTTP server
	pprofMux := http.NewServeMux()
	pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
	pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Start pprof server
	go func() {
		log.Println("Starting pprof server on :6060")
		if err := http.ListenAndServe(":6060", pprofMux); err != nil {
			log.Printf("Failed to start pprof server: %v", err)
		}
	}()
}

func setupMetrics() {
	// Create metrics HTTP server
	metricsMux := http.NewServeMux()
	metricsMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// In a real app, you would use Prometheus registry to collect and serve metrics
		w.Write([]byte(`# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 10
# HELP netgex_example_counter Example counter
# TYPE netgex_example_counter counter
netgex_example_counter{label="test"} 42
`))
	})

	// Start metrics server
	go func() {
		log.Println("Starting metrics server on :9091")
		if err := http.ListenAndServe(":9091", metricsMux); err != nil {
			log.Printf("Failed to start metrics server: %v", err)
		}
	}()
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
