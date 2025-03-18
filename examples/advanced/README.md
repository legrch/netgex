# Bootstrap Example

This example demonstrates how to use the bootstrap components:

- `service.Registrar` interface
- `grpc.Server` implementation
- `gateway.Server` implementation
- `entrypoint.Entrypoint` implementation

## Features

- gRPC server with reflection
- HTTP/REST gateway
- Swagger UI
- CORS support
- Health checks
- Graceful shutdown

## Running the Example

```bash
go run main.go
```

## Testing the Example

### HTTP Endpoints

```bash
# Test the HTTP endpoint with path parameter
curl http://localhost:8080/api/v1/hello/John

# Test the HTTP endpoint with query parameter
curl http://localhost:8080/api/v1/hello?name=Jane
```

### Health Check

```bash
# Test the health check endpoint
curl http://localhost:8080/health
```

### Swagger UI

Open the Swagger UI in your browser:

```
http://localhost:8080/swagger/
```

## Implementation Details

### Service Registrar

The `GreeterService` implements the `service.Registrar` interface:

```go
type Registrar interface {
    // RegisterGRPC registers the gRPC service with the gRPC server
    RegisterGRPC(*grpc.Server)

    // RegisterHTTP registers the HTTP/REST handlers with the gateway mux
    RegisterHTTP(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}
```

This interface allows services to register themselves with both gRPC and HTTP/REST servers.

### Entrypoint

The `Entrypoint` is responsible for:

1. Creating and configuring the gRPC server
2. Creating and configuring the HTTP/REST gateway
3. Running all processes
4. Handling graceful shutdown

It uses the options pattern for configuration, allowing for flexible setup. 