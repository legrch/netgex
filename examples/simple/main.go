package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/legrch/netgex/server"
)

// ExampleService is a simple implementation adapter for demonstration purposes
type ExampleService struct{}

// RegisterGRPC registers the gRPC service with the server (deprecated)
func (s *ExampleService) RegisterGRPC(server *grpc.Server) {
	// In a real application, you would register your gRPC service here
	// Example: pb.RegisterExampleServiceServer(server, s)
}

// RegisterServer registers the gRPC service with the server
func (s *ExampleService) RegisterServer(server *grpc.Server) {
	// In a real application, you would register your gRPC service here
	// Example: pb.RegisterExampleServiceServer(server, s)
}

// RegisterGateway registers the gRPC-Gateway handlers
func (s *ExampleService) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	// In a real application, you would register your gRPC-Gateway handlers here
	// Example: return pb.RegisterExampleServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)

	// For demonstration, we'll register a simple health check handler
	err := mux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK"}`))
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to register health handler: %v", err)
	}

	// And a simple hello world handler
	err = mux.HandlePath("GET", "/hello/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		name := pathParams["name"]
		if name == "" {
			name = "world"
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Hello, ` + name + `!"}`))
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to register hello handler: %v", err)
	}

	return nil
}

// RegisterHTTP registers the HTTP/REST handlers with the gateway mux
func (s *ExampleService) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return s.RegisterGateway(ctx, mux, endpoint, opts)
}

func main() {
	// Create a context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create a logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create a service adapter
	service := &ExampleService{}

	// Create entrypoint options
	opts := []server.Option{
		server.WithLogger(logger),
		server.WithGRPCAddress(":50051"),
		server.WithHTTPAddress(":8000"),
		server.WithMetricsAddress(":9090"),
		server.WithPprofAddress(":6060"),
		server.WithCloseTimeout(5 * time.Second),
		server.WithReflection(true),
		server.WithServices(service),
	}

	// Create entrypoint
	ep := server.NewServer(opts...)

	// Log startup message
	logger.Info("starting example service",
		"grpc_address", ":50051",
		"http_address", ":8000",
		"metrics_address", ":9090",
	)

	// Run the entrypoint
	if err := ep.Run(ctx); err != nil {
		logger.Error("entrypoint error", "error", err)
		os.Exit(1)
	}

	// Log shutdown message
	logger.Info("example service shutdown complete")
}
