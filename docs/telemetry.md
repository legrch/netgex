# NetGeX Telemetry Guide

This guide explains how to configure observability in NetGeX using the telemetry system.

## Overview

NetGeX includes a comprehensive observability system supporting:

- **Tracing**: Distributed request tracing with OpenTelemetry or Jaeger
- **Metrics**: Prometheus or OTLP metrics
- **Logging**: Structured logging with slog
- **Profiling**: Continuous profiling with Pyroscope or pprof

## Configuration

You can configure telemetry through:

1. Environment variables
2. Go code with functional options
3. Docker Compose files

### Using Environment Variables

```bash
# Service Info
export SERVICE_NAME=netgex
export SERVICE_VERSION=1.0.0
export ENVIRONMENT=production

# Tracing Configuration
export TRACING_ENABLED=true
export TRACING_BACKEND=otlp  # Options: otlp, jaeger, none
export TRACING_ENDPOINT=otel-collector:4318
export TRACING_INSECURE=true
export TRACING_SAMPLE_RATE=0.1  # 10% sampling in production

# Metrics Configuration
export METRICS_ENABLED=true
export METRICS_BACKEND=prometheus  # Options: prometheus, otlp, none
export METRICS_PATH=/metrics
export METRICS_PORT=9091

# Logging Configuration
export LOGGING_ENABLED=true
export LOGGING_BACKEND=stdout  # Options: stdout, file, otlp, global
export LOGGING_LEVEL=info      # Options: debug, info, warn, error
export LOGGING_FORMAT=json     # Options: json, text

# Profiling Configuration
export PROFILING_ENABLED=true
export PROFILING_BACKEND=pyroscope  # Options: pyroscope, pprof, none
export PROFILING_ENDPOINT=http://pyroscope:4040
export PROFILING_TYPES=cpu,heap  # Comma-separated list
```

### Using Server Options

```go
// Basic usage - enable all telemetry with default settings
srv := server.New(
    server.WithTelemetry(),
)

// Advanced usage - configure specific backends
srv := server.New(
    server.WithLogger(logger),  // Your custom logger
    server.WithConfig(cfg),     // Your custom config
    
    // Enable telemetry
    server.WithTelemetry(),
    
    // Configure backends
    server.WithTracingBackend("otlp", "otel-collector:4318"),
    server.WithMetricsBackend("prometheus", ""),
    server.WithProfilingBackend("pyroscope", "http://pyroscope:4040"),
)

// Full programmatic configuration
cfg := config.NewConfig()
cfg.Telemetry.Tracing.Enabled = true
cfg.Telemetry.Tracing.Backend = "otlp"
cfg.Telemetry.Tracing.Endpoint = "otel-collector:4318"
cfg.Telemetry.Tracing.SampleRate = 1.0

cfg.Telemetry.Metrics.Enabled = true
cfg.Telemetry.Metrics.Backend = "prometheus"
cfg.Telemetry.Metrics.Path = "/metrics"

cfg.Telemetry.Profiling.Enabled = true
cfg.Telemetry.Profiling.Backend = "pprof"

srv := server.New(
    server.WithConfig(cfg),
    server.WithTelemetry(),
)
```

## Deployment Scenarios

### Local Development with OpenTelemetry

For local development with a full observability stack:

```bash
cd examples/docker-compose/otlp-local
docker-compose up
```

This starts:
- Your application
- OpenTelemetry Collector
- Prometheus
- Tempo (for traces)
- Pyroscope (for profiles)
- Grafana (UI for all data)

### Legacy Systems (Prometheus + Jaeger + pprof)

For integrating with existing Prometheus + Jaeger setups:

```bash
cd examples/docker-compose/legacy
docker-compose up
```

This starts:
- Your application
- Prometheus
- Jaeger
- Grafana
- pprof endpoint (for on-demand profiling)

### Grafana Cloud Integration

For production use with Grafana Cloud:

```bash
# Set your API key
export GRAFANA_CLOUD_API_KEY=your-api-key

cd examples/docker-compose/grafana-cloud
docker-compose up
```

Configuration:
```go
srv := server.New(
    // Cloud endpoints
    server.WithTracingBackend("otlp", "tempo-us-central1.grafana.net:4318"),
    server.WithMetricsBackend("otlp", "prometheus-us-central1.grafana.net:4318"),
    server.WithProfilingBackend("pyroscope", "profiles-prod-006.grafana.net:4040"),
)
```

## Relationship with Logging

The telemetry system integrates with logging in several ways:

1. **Using WithLogger**: If you provide a logger with `server.WithLogger(logger)`, the telemetry system will respect it and add service attributes.

2. **Configuring via Environment**: You can configure logging format, level, and backend via environment variables.

3. **OTLP Logging**: Support for OTLP logging is in progress (currently experimental in the OTel SDK).

4. **Structured Logs**: All logs are structured and include service, version, and environment context.

Example:
```go
// Your custom logger
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// Telemetry will respect and enhance your logger
srv := server.New(
    server.WithLogger(logger),
    server.WithTelemetry(),
)
```

## FAQ

### Do I need to use both `WithLogger` and configure `Telemetry.Logging`?

No. If you use `WithLogger`, that logger will be respected and enhanced with service attributes. The `Telemetry.Logging` configuration controls the format and level when a logger isn't provided, or when sending logs to backends like OTLP.

### How do I use pprof instead of Pyroscope?

Set `PROFILING_BACKEND=pprof` or use `server.WithProfilingBackend("pprof", "")`. The pprof endpoint will be available at the address configured with `PPROF_ADDRESS` (default `:6060`).

### How can I see what telemetry is enabled?

The server splash screen shows active telemetry backends:

```
 ╭──────────────────────────────────────────╮
 │   NetGeX API Server                      │
 │   Version: 1.0.0                         │
 ├──────────────────────────────────────────┤
 │   gRPC:       :9090                      │
 │   HTTP:       :8080                      │
 │   Metrics:    :9091                      │
 │   pprof:      :6060                      │
 ├──────────────────────────────────────────┤
 │   ✓ gRPC Reflection                      │
 │   ✓ Health Checks                        │
 │   ✓ Tracing (otlp)                       │
 │   ✓ Metrics (prometheus)                 │
 │   ✓ Profiling (pprof)                    │
 ╰──────────────────────────────────────────╯
```

## Recommended Configuration

### Development

```go
srv := server.New(
    server.WithTelemetry(),
    server.WithTracingBackend("otlp", "otel-collector:4318"),
    server.WithMetricsBackend("prometheus", ""),
    server.WithProfilingBackend("pprof", ""),
)
```

### Production (Grafana Cloud)

```go
srv := server.New(
    server.WithTelemetry(),
    server.WithTracingBackend("otlp", "tempo-us-central1.grafana.net:4318"),
    server.WithMetricsBackend("otlp", "prometheus-us-central1.grafana.net:4318"),
    server.WithProfilingBackend("pyroscope", "profiles-prod-006.grafana.net:4040"),
)
```

### Production (Self-hosted)

```go
srv := server.New(
    server.WithTelemetry(),
    server.WithTracingBackend("otlp", "tempo:4317"),
    server.WithMetricsBackend("prometheus", ""),
    server.WithProfilingBackend("pyroscope", "http://pyroscope:4040"),
)
``` 