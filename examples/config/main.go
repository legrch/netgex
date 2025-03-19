package main

import (
	"context"
	"github.com/legrch/netgex/config"
	"github.com/legrch/netgex/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Set up logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Method 1: Load configuration from environment variables

	cfg := &config.Config{
		LogLevel:           "debug",
		CloseTimeout:       10 * time.Second,
		AppName:            "Service",
		AppVersion:         "dev",
		GRPCAddress:        ":9090",
		HTTPAddress:        ":8080",
		MetricsAddress:     ":9091",
		PprofAddress:       ":6060",
		ReflectionEnabled:  true,
		HealthCheckEnabled: true,
		SwaggerDir:         "./api/swagger",
		SwaggerBasePath:    "/api",
	}

	// Set logger (not loaded from env because it's a complex object)
	// cfg.Logger = logger

	// Create server with config
	server := server.NewServer(
		server.WithConfig(cfg),
		// netgex.WithLogger(logger),
	)

	// Method 2: Configure directly (alternative approach)
	/*
		server := netgex.NewServer(
			netgex.WithLogger(logger),
			netgex.WithGRPCAddress(":9090"),
			netgex.WithHTTPAddress(":8080"),
			netgex.WithAppName("MyService"),
			netgex.WithAppVersion("1.0.0"),
		)
	*/

	// Create a context that will be canceled on signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Run the server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run(ctx)
	}()

	// Wait for signal or error
	select {
	case <-sigCh:
		logger.Info("received shutdown signal")
		cancel()
	case err := <-errCh:
		logger.Error("server error", "error", err)
	}

	// Wait for graceful shutdown
	select {
	case <-time.After(20 * time.Second):
		logger.Error("shutdown timeout exceeded")
	case err := <-errCh:
		if err != nil {
			logger.Error("shutdown error", "error", err)
		} else {
			logger.Info("server shutdown successfully")
		}
	}
}
