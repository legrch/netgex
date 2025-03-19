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
- JSON marshaling customization
- CORS configuration

## Package Structure

- `server/` - Main server implementation
- `service/` - Service registration interfaces
- `config/` - Configuration utilities
- `splash/` - Terminal startup display
- `internal/` - Internal implementation details:
  - `grpc/` - gRPC server implementation
  - `gateway/` - HTTP/REST gateway server implementation
  - `metrics/` - Metrics server for Prometheus
  - `pprof/` - Profiling server
  - `pyroscope/` - Continuous profiling
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
| `LOG_LEVEL` | Logging level | `info` |
| `GRPC_ADDRESS` | gRPC server address | `:9090` |
| `HTTP_ADDRESS` | HTTP/REST gateway address | `:8080` |
| `METRICS_ADDRESS` | Metrics server address | `:9091` |
| `PPROF_ADDRESS` | pprof server address | `:6060` |
| `REFLECTION_ENABLED` | Enable gRPC reflection | `true` |
| `HEALTH_CHECK_ENABLED` | Enable health checks | `true` |
| `CLOSE_TIMEOUT` | Timeout for graceful shutdown | `10s` |
| `SWAGGER_ENABLED` | Enable Swagger UI | `true` |
| `SWAGGER_DIR` | Directory containing swagger files | `./api` |
| `SWAGGER_BASE_PATH` | Base path for swagger UI | `/` |

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
	server.WithGatewayCORS(corsOptions),
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
		server.WithGatewayCORS(corsOptions),
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
- `WithConfig(config *config.Config)` - Sets the configuration for the server
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

### Server Options
- `WithGRPCServerOptions(options ...grpc.ServerOption)` - Sets additional options for the gRPC server
- `WithGRPCUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor)` - Sets the unary interceptors for the gRPC server
- `WithGRPCStreamInterceptors(interceptors ...grpc.StreamServerInterceptor)` - Sets the stream interceptors for the gRPC server
- `WithGatewayMuxOptions(options ...runtime.ServeMuxOption)` - Sets the ServeMux options for the gateway server
- `WithGatewayCORS(options cors.Options)` - Enables CORS with the specified options for the gateway

### JSON Options
The gateway server supports customizable JSON marshaling through `runtime.ServeMuxOption`:

```go
server.WithGatewayMuxOptions(
	runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
			UseEnumNumbers:  false,
			AllowPartial:    true,
			Multiline:       true,
			Indent:          "  ",
		},
	}),
)
```

Options include:
- `UseProtoNames` - Use proto field names instead of JSON names
- `EmitUnpopulated` - Include unpopulated fields in output
- `UseEnumNumbers` - Use enum numbers instead of string names
- `AllowPartial` - Allow partial messages
- `Multiline` - Format output with multiple lines
- `Indent` - Set indentation for multiline output

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
- `advanced/` - Advanced configuration with service registration, Swagger, and more
- `json/` - Examples of JSON configuration for the gateway
- `config/` - Environment-based configuration examples

## Development

### Task Runner

This project uses [Task](https://taskfile.dev) as a task runner instead of Make. To get started:

1. Install Task: `go install github.com/go-task/task/v3/cmd/task@latest`
2. Run `task` to see available commands

Common tasks:
- `task lint` - Run linters
- `task test` - Run tests
- `task test:coverage` - Run tests with coverage report
- `task mock` - Generate mocks

### Release Process

The project follows a structured release process:

1. Prepare a release: `task prepare-release -- 1.2.3`
2. Update the CHANGELOG.md
3. Create the release: `task release`

For detailed information about creating releases, see [docs/RELEASING.md](docs/RELEASING.md).

## Related Documentation

### Official Documentation
- [gRPC](https://grpc.io/docs/) - gRPC documentation
- [gRPC Gateway](https://grpc-ecosystem.github.io/grpc-gateway/) - gRPC Gateway documentation
- [Protocol Buffers](https://protobuf.dev/) - Protocol Buffers documentation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.