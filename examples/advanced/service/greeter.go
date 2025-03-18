package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GreeterService implements the Greeter service
type GreeterService struct{}

// NewGreeterService creates a new GreeterService
func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

// RegisterGRPC registers the gRPC service with the gRPC server
func (s *GreeterService) RegisterGRPC(server *grpc.Server) {
	// In a real implementation, you would register the generated gRPC service
	// Example: pb.RegisterGreeterServiceServer(server, s)
}

// RegisterHTTP registers the HTTP/REST handlers with the gateway mux
func (s *GreeterService) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	// In a real implementation, you would register the generated HTTP handlers
	// Example: return pb.RegisterGreeterServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)

	// Register the endpoint defined in the proto file
	err := mux.HandlePath("GET", "/api/v1/greeter/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		name := pathParams["name"]
		if name == "" {
			name = "World"
		}

		response := map[string]interface{}{
			"message":   fmt.Sprintf("Hello, %s!", name),
			"timestamp": time.Now().Unix(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	if err != nil {
		return err
	}

	// Keep the existing handlers for backward compatibility
	err = mux.HandlePath("GET", "/api/v1/hello/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		name := pathParams["name"]
		if name == "" {
			name = "World"
		}

		response := map[string]interface{}{
			"message":   fmt.Sprintf("Hello, %s!", name),
			"timestamp": time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	if err != nil {
		return err
	}

	// Add a query parameter version
	return mux.HandlePath("GET", "/api/v1/hello", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}

		response := map[string]interface{}{
			"message":   fmt.Sprintf("Hello, %s!", name),
			"timestamp": time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

// SayHello is the implementation of the SayHello RPC method
func (s *GreeterService) SayHello(ctx context.Context, req interface{}) (interface{}, error) {
	// This would be called in a real gRPC implementation
	// For now, we don't need this since we're using HandlePath directly
	return nil, status.Error(codes.Unimplemented, "method SayHello not implemented")
}
