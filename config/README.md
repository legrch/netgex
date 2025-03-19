# Netgex Configuration Package

This package provides a simplified configuration system for the Netgex server framework using the [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig) library. It centralizes all configuration options in a single, easily accessible location.

## Features

- Environment variable support with sensible defaults
- Type-safe configuration via Go structs with struct tags
- JSON serialization support
- Configuration for all Netgex server components

## Usage

### Basic Usage with Environment Variables

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "github.com/legrch/netgex"
    "github.com/legrch/netgex/config"
)

func main() {
    // Load configuration from environment variables with prefix "NETGEX"
    cfg, err := config.LoadFromEnv("NETGEX")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Set logger (not loaded from env)
    cfg.Logger = slog.Default()
    
    // Create server with config
    server := netgex.NewServer(
        netgex.WithConfig(cfg),
    )
    
    // Run the server
    server.Run(context.Background())
}
```

### Custom Configuration

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "time"
    "github.com/legrch/netgex"
    "github.com/legrch/netgex/config"
)

func main() {
    // Create custom configuration
    cfg := &config.Config{
        Logger:            slog.Default(),
        CloseTimeout:      15 * time.Second,
        GRPCAddress:       ":9090",
        HTTPAddress:       ":8080",
        MetricsAddress:    ":9091",
        PprofAddress:      ":6060",
        SwaggerDir:        "./api",
        SwaggerBasePath:   "/",
        ReflectionEnabled: true,
        HealthCheckEnabled: true,
        AppName:           "MyService",
        AppVersion:        "1.0.0",
        JSONUseProtoNames: true,
        JSONEmitUnpopulated: true,
        JSONUseEnumNumbers: true,
        JSONAllowPartial: true,
        JSONMultiline: true,
        JSONIndent: "  ",
    }
    
    // Create server with config
    server := netgex.NewServer(
        netgex.WithConfig(cfg),
    )
    
    // Run the server
    server.Run(context.Background())
}
```

## Environment Variables

The configuration system supports the following environment variables:

| Variable                     | Description                          | Default     |
|------------------------------|--------------------------------------|-------------|
| `<PREFIX>_CLOSE_TIMEOUT`     | Timeout for graceful shutdown        | `10s`       |
| `<PREFIX>_APP_NAME`          | Application name                     | `Service`   |
| `<PREFIX>_APP_VERSION`       | Application version                  | `dev`       |
| `<PREFIX>_GRPC_ADDRESS`      | gRPC server address                  | `:9090`     |
| `<PREFIX>_HTTP_ADDRESS`      | HTTP/REST gateway address            | `:8080`     |
| `<PREFIX>_METRICS_ADDRESS`   | Metrics server address               | `:9091`     |
| `<PREFIX>_PPROF_ADDRESS`     | Pprof debug server address           | `:6060`     |
| `<PREFIX>_SWAGGER_DIR`       | Directory containing swagger files   | `./api`     |
| `<PREFIX>_SWAGGER_BASE_PATH` | Base path for swagger UI             | `/`         |
| `<PREFIX>_REFLECTION_ENABLED`| Enable gRPC reflection              | `true`      |
| `<PREFIX>_HEALTH_CHECK_ENABLED` | Enable health checks            | `true`      |
| `<PREFIX>_CORS_ENABLED`      | Enable CORS                          | `false`     |
| `<PREFIX>_CORS_ALLOWED_ORIGINS` | CORS allowed origins             | `*`         |
| `<PREFIX>_CORS_ALLOWED_METHODS` | CORS allowed methods             | `GET,POST,PUT,DELETE,OPTIONS` |
| `<PREFIX>_CORS_ALLOWED_HEADERS` | CORS allowed headers             | `Origin,Accept,Content-Type,X-Requested-With,Authorization` |
| `<PREFIX>_JSON_USE_PROTO_NAMES` | Use protobuf field names in JSON  | `true`      |
| `<PREFIX>_JSON_EMIT_UNPOPULATED` | Include unpopulated fields in JSON | `true`   |
| `<PREFIX>_JSON_USE_ENUM_NUMBERS` | Use enum numbers instead of names | `true`     |
| `<PREFIX>_JSON_ALLOW_PARTIAL` | Allow partial JSON messages         | `true`      |
| `<PREFIX>_JSON_MULTILINE`    | Format JSON with multiple lines      | `true`      |
| `<PREFIX>_JSON_INDENT`       | JSON indentation string              | `  ` (2 spaces) |

Where `<PREFIX>` is the prefix you specify when calling `config.LoadFromEnv()`.

## CORS Configuration

To enable CORS support:

```
NETGEX_CORS_ENABLED=true
NETGEX_CORS_ALLOWED_ORIGINS=https://example.com,https://api.example.com
NETGEX_CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
NETGEX_CORS_ALLOWED_HEADERS=Origin,Authorization,Content-Type
```

## JSON Configuration

To customize JSON marshaling:

```
NETGEX_JSON_USE_PROTO_NAMES=true
NETGEX_JSON_EMIT_UNPOPULATED=false
NETGEX_JSON_USE_ENUM_NUMBERS=false
NETGEX_JSON_ALLOW_PARTIAL=true
NETGEX_JSON_MULTILINE=true
NETGEX_JSON_INDENT="  "
``` 