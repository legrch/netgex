# Bootstrap Example

This example demonstrates how to use the `entrypoint` package to create a simple gRPC and HTTP service.

## Overview

The example implements a simple service with the following endpoints:

- HTTP: `GET /health` - Returns a health check response
- HTTP: `GET /hello/{name}` - Returns a greeting message

The service is configured to run with:
- gRPC server on port 50051
- HTTP gateway on port 8000
- Metrics server on port 9090
- pprof server on port 6060

## Running the Example

To run the example:

```bash
go run main.go
```

## Testing the Endpoints

### HTTP Endpoints

Test the hello endpoint:

```bash
curl http://localhost:8000/hello/world
```

Expected response:
```json
{"message":"Hello, world!"}
```

Test the health endpoint:

```bash
curl http://localhost:8000/health
```

Expected response:
```json
{"status":"OK"}
```

### Metrics Endpoint

View the metrics:

```bash
curl http://localhost:9090/metrics
```

### pprof Endpoint

View the pprof index:

```bash
curl http://localhost:6060/debug/pprof/
```

## Notes

This is a simplified example that doesn't use actual protobuf definitions. In a real application, you would:

1. Define your service using Protocol Buffers
2. Generate the gRPC and gRPC-Gateway code
3. Implement the generated service interfaces
4. Register the services with the gRPC server and gateway

The example simulates these steps to demonstrate the usage of the `entrypoint` package without requiring the protobuf toolchain. 