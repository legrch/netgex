package service

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Registrar is an interface for gRPC service implementations that can register
// themselves with both gRPC and HTTP/REST gateway servers
type Registrar interface {
	// RegisterGRPC registers the gRPC service with the gRPC server
	RegisterGRPC(*grpc.Server)

	// RegisterHTTP registers the HTTP/REST handlers with the gateway mux
	RegisterHTTP(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}
