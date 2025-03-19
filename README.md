# NetGeX - Go Network and gRPC Extensions

[![Go Reference](https://pkg.go.dev/badge/github.com/legrch/server.svg)](https://pkg.go.dev/github.com/legrch/netgex)
[![Go Report Card](https://goreportcard.com/badge/github.com/legrch/netgex)](https://goreportcard.com/report/github.com/legrch/netgex)
[![License](https://img.shields.io/github/license/legrch/netgex)](LICENSE)
[![Release](https://img.shields.io/github/v/release/legrch/netgex)](https://github.com/legrch/netgex/releases)

NetGeX is a comprehensive toolkit for building robust gRPC and HTTP/REST services in Go with integrated support for metrics, profiling, health checks, and graceful shutdown management.

## Overview

This package offers a unified approach to create and manage various server components:
- gRPC servers with reflection and health checks
- HTTP/REST gateway for gRPC services
- Metrics servers for Prometheus
- Profiling servers with pprof
- Graceful shutdown management
- Environment-based configuration

## Package Structure

- `grpc/` - gRPC server implementation
- `gateway/` - HTTP/REST gateway server implementation
- `service/` - Service registration interfaces
- `metrics/` - Metrics server for Prometheus
- `pprof/` - Profiling server
- `examples/` - Example implementations

## Usage

### Prerequisites
- Go 1.18+
- Protocol buffer compiler (for gRPC services)
- gRPC Gateway (for HTTP/REST support)

### Installation

```bash
go get -u github.com/legrch/netgex
```

### Configuration

The server package uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `GRPC_ADDRESS` | gRPC server address | `:9090` |
| `HTTP_ADDRESS` | HTTP/REST gateway address | `:8080` |
| `METRICS_ADDRESS` | Metrics server address | `:9091` |
| `PPROF_ADDRESS` | pprof server address | `:6060` |
| `REFLECTION_ENABLED` | Enable gRPC reflection | `true` |
| `CLOSE_TIMEOUT` | Timeout for graceful shutdown | `10s` |
| `SWAGGER_DIR` | Directory containing swagger files | `./api` |
| `SWAGGER_BASE_PATH` | Base path for swagger UI | `/` |
| `PROJECT_NAME` | Application name for splash screen | `Service` |
| `VERSION` | Application version for splash screen | `dev` |

### Components

#### Service Registrar

The `service.Registrar` interface defines the contract for services that can register with both gRPC and HTTP/REST servers:

```go
type Registrar interface {
	// RegisterGRPC registers the gRPC service with the gRPC server
	RegisterGRPC(*grpc.Server)

	// RegisterHTTP registers the HTTP/REST handlers with the gateway mux
	RegisterHTTP(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}
```

#### gRPC Server

The `grpc.Server` provides a configurable gRPC server with support for:

- Multiple service registrations
- Custom interceptors
- Health checks
- gRPC reflection
- Graceful shutdown

```go
server := grpc.NewServer(
	logger,
	5*time.Second,
	":50051",
	grpc.WithServices(myService),
	grpc.WithUnaryInterceptors(
		grpc_recovery.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
	),
	grpc.WithReflection(true),
	grpc.WithHealthCheck(true),
)

if err := server.Run(ctx); err != nil {
	log.Fatal(err)
}
```

#### Gateway Server

The `gateway.Server` provides a configurable HTTP/REST gateway server with support for:

- Multiple service registrations
- Custom header matchers
- CORS configuration
- Swagger UI
- Health checks
- Graceful shutdown

```go
gateway := gateway.NewServer(
	logger,
	5*time.Second,
	":50051",  // gRPC server address
	":8080",   // HTTP server address
	gateway.WithServices(myService),
	gateway.WithCORS(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}),
	gateway.WithSwagger("./api/swagger", "/api/v1"),
)

if err := gateway.Run(ctx); err != nil {
	log.Fatal(err)
}
```

#### Main Server

The `server.Server` provides a unified way to initialize and run your application with all components:

```go
// Create server options
opts := []server.Option{
	server.WithLogger(logger),
	server.WithServices(myService),
	server.WithCloseTimeout(5 * time.Second),
	server.WithGRPCAddress(":50051"),
	server.WithHTTPAddress(":8080"),
	server.WithReflection(true),
	server.WithHealthCheck(true),
	server.WithCORS(&corsOptions),
	server.WithSwaggerDir("./api/swagger"),
	server.WithSwaggerBasePath("/api/v1"),
}

// Create server
srv := server.NewServer(opts...)

// Run the server
if err := srv.Run(ctx); err != nil {
	logger.Error("server error", "error", err)
	os.Exit(1)
}
```

### Complete Example

Here's a complete example of using the server package:

```go
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/legrch/netgex/server"
	"github.com/legrch/netgex/service"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/runtime"
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
	myService := NewMyServiceRegistrar()

	// Create CORS options
	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
	}

	// Create server options
	opts := []server.Option{
		server.WithLogger(logger),
		server.WithServices(myService),
		server.WithCloseTimeout(5 * time.Second),
		server.WithGRPCAddress(":50051"),
		server.WithHTTPAddress(":8080"),
		server.WithReflection(true),
		server.WithHealthCheck(true),
		server.WithCORS(&corsOptions),
		server.WithSwaggerDir("./api/swagger"),
		server.WithSwaggerBasePath("/api/v1"),
	}

	// Create server
	srv := server.NewServer(opts...)

	// Run the server
	if err := srv.Run(ctx); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

// NewMyServiceRegistrar creates a new service registrar
func NewMyServiceRegistrar() service.Registrar {
	// Implement your service registrar here
	return &myServiceRegistrar{}
}

// myServiceRegistrar implements the service.Registrar interface
type myServiceRegistrar struct{}

// RegisterGRPC registers the gRPC service with the gRPC server
func (s *myServiceRegistrar) RegisterGRPC(server *grpc.Server) {
	// Register your gRPC service here
	// Example: pb.RegisterMyServiceServer(server, s.service)
}

// RegisterHTTP registers the HTTP/REST handlers with the gateway mux
func (s *myServiceRegistrar) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	// Register your HTTP/REST handlers here
	// Example: return pb.RegisterMyServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	return nil
}
```

## Configuration Options

The server package provides the following configuration options:

### Basic Options
- `WithLogger(logger *slog.Logger)` - Sets the logger for the server
- `WithCloseTimeout(timeout time.Duration)` - Sets the timeout for graceful shutdown
- `WithGRPCAddress(address string)` - Sets the gRPC server address
- `WithHTTPAddress(address string)` - Sets the HTTP server address
- `WithMetricsAddress(address string)` - Sets the metrics server address
- `WithPprofAddress(address string)` - Sets the pprof server address
- `WithSwaggerDir(dir string)` - Sets the directory containing swagger files
- `WithSwaggerBasePath(path string)` - Sets the base path for swagger UI
- `WithReflection(enabled bool)` - Enables or disables gRPC reflection
- `WithHealthCheck(enabled bool)` - Enables or disables health checks
- `WithServices(registrars ...service.Registrar)` - Sets the service registrars
- `WithProcesses(processes ...Process)` - Adds additional processes to the server
- `WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor)` - Sets the unary interceptors for the gRPC server
- `WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor)` - Sets the stream interceptors for the gRPC server
- `WithCORS(options *cors.Options)` - Enables CORS with the specified options
- `WithAppName(name string)` - Sets the application name for the splash screen
- `WithAppVersion(version string)` - Sets the application version for the splash screen

### JSON Options
- `WithJSONConfig(config *gateway.JSONConfig)` - Sets the JSON configuration for the gateway
- `WithJSONConfigFromEnv()` - Configures the JSON options from environment variables
- `WithJSONUseProtoNames(useProtoNames bool)` - Enables or disables using proto field names
- `WithJSONEmitUnpopulated(emitUnpopulated bool)` - Enables or disables emitting unpopulated fields
- `WithJSONUseEnumNumbers(useEnumNumbers bool)` - Enables or disables using enum numbers
- `WithJSONAllowPartial(allowPartial bool)` - Enables or disables allowing partial messages
- `WithJSONMultiline(multiline bool)` - Enables or disables multiline output
- `WithJSONIndent(indent string)` - Sets the indentation string for multiline output

## Custom Processes

You can add custom processes to the server by implementing the `Process` interface:

```go
type Process interface {
	PreRun(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
```

Then add your process to the server:

```go
server.WithProcesses(yourCustomProcess)
```

## Examples

See the `examples/` directory for complete examples of how to use the server package:

- `simple/` - Basic gRPC and gateway server setup
- `advanced/` - Advanced configuration with interceptors, Swagger, and more
- `json/` - Examples of JSON configuration for the gateway

## Related Documentation

### Official Documentation
- [gRPC](https://grpc.io/docs/) - gRPC documentation
- [gRPC Gateway](https://grpc-ecosystem.github.io/grpc-gateway/) - gRPC Gateway documentation
- [Protocol Buffers](https://protobuf.dev/) - Protocol Buffers documentation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.