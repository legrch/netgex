package main

import (
	"context"
	"github.com/legrch/netgex/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"

	"github.com/legrch/netgex/examples/advanced/service"
)

func main() {
	// Create a context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create a logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create service registrars
	greeterService := service.NewGreeterService()

	// Create CORS options
	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token"},
	}

	// Create entrypoint options
	opts := []server.Option{
		server.WithLogger(logger),
		server.WithServices(greeterService),
		server.WithCloseTimeout(5 * time.Second),
		server.WithGRPCAddress(":50051"),
		server.WithHTTPAddress(":8080"),
		server.WithReflection(true),
		server.WithHealthCheck(true),
		server.WithGatewayCORS(corsOptions),
		server.WithSwaggerDir("./api/swagger"),
		server.WithSwaggerBasePath("/api"),
	}

	// Create entrypoint
	ep := server.NewServer(opts...)

	// Run the entrypoint
	logger.Info("starting application")
	if err := ep.Run(ctx); err != nil {
		logger.Error("entrypoint error", "error", err)
		os.Exit(1)
	}
}
