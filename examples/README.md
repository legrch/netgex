# NetGeX Examples

This directory contains examples demonstrating how to use the NetGeX package to create gRPC and HTTP services.

## Available Examples

### [Simple Example](./simple/)

A basic example showing how to use the NetGeX package to create a simple gRPC and HTTP service with minimal configuration.

Features:
- Basic HTTP endpoints
- Metrics server
- pprof profiler

### [Advanced Example](./advanced/)

A more complex example demonstrating advanced features of the NetGeX package.

Features:
- gRPC service with real proto files
- HTTP/REST gateway
- Swagger documentation
- Custom interceptors (logging and metrics)
- Prometheus metrics
- pprof profiler

### [JSON Example](./json/)

Examples of JSON configuration options for the gateway.

Features:
- Configuring JSON options from environment variables
- Setting JSON options programmatically

## Running the Examples

Each example can be run from its own directory:

```bash
# Simple example
cd simple
go run main.go

# Advanced example
cd advanced
go run main.go

# JSON example
cd json
go run main.go
```

## Testing the Examples

See the README.md in each example directory for detailed instructions on how to test the endpoints. 