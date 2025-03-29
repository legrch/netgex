# NetGeX Observability

NetGeX supports comprehensive observability via OpenTelemetry and other backends. This document explains how to configure and use the observability features.

## Overview

The observability stack includes:

- **Tracing**: Distributed tracing for request flows
- **Metrics**: Measurements of system behavior and performance
- **Logging**: Structured application logs
- **Profiling**: Continuous profiling for CPU and memory

## Configuration

You can configure observability through environment variables or by using the server options:

### Using Environment Variables

```bash
# Service Info
export SERVICE_NAME=netgex
export SERVICE_VERSION=1.0.0
export ENVIRONMENT=production

# Tracing
export TRACING_ENABLED=true
export TRACING_BACKEND=otlp  # otlp, jaeger, none
export TRACING_ENDPOINT=otel-collector:4318
export TRACING_SAMPLE_RATE=0.1  # 10% sampling

# Metrics
export METRICS_ENABLED=true
export METRICS_BACKEND=prometheus  # prometheus, otlp, none
export METRICS_PATH=/metrics

# Logging
export LOGGING_LEVEL=info
export LOGGING_FORMAT=json  # json, text

# Profiling
export PROFILING_ENABLED=true
export PROFILING_BACKEND=pyroscope  # pyroscope, none
export PROFILING_ENDPOINT=http://pyroscope:4040
```

### Using Server Options

```go
srv := server.New(
    // Core server options
    server.WithConfig(cfg),
    server.WithLogger(logger),
    
    // Basic telemetry activation
    server.WithTelemetry(),
    
    // Specific backends
    server.WithTracingBackend("otlp", "otel-collector:4318"),
    server.WithMetricsBackend("prometheus", ""),
    server.WithProfilingBackend("pyroscope", "http://pyroscope:4040"),
)
```

## Backends

### Tracing Backends

- **OTLP**: OpenTelemetry Protocol - modern and vendor-neutral
  - Compatible with Tempo, Jaeger, or any OTLP-compatible backend
  - Example: `WithTracingBackend("otlp", "otel-collector:4318")`

- **Jaeger**: Legacy tracing (direct)
  - Example: `WithTracingBackend("jaeger", "http://jaeger:14268/api/traces")`

### Metrics Backends

- **Prometheus**: Exposes metrics at /metrics endpoint
  - Compatible with Prometheus, Mimir, or VictoriaMetrics
  - Example: `WithMetricsBackend("prometheus", "")`

- **OTLP**: OpenTelemetry Protocol for metrics
  - Example: `WithMetricsBackend("otlp", "otel-collector:4318")`

### Profiling Backends

- **Pyroscope**: Continuous profiling
  - Compatible with Grafana Phlare and Pyroscope
  - Example: `WithProfilingBackend("pyroscope", "http://pyroscope:4040")`

## Docker Compose Setup

For local development, you can use this Docker Compose setup:

```yaml
version: "3.8"

services:
  app:
    # Your app configuration
    environment:
      - SERVICE_NAME=netgex
      - TRACING_ENABLED=true
      - TRACING_BACKEND=otlp
      - TRACING_ENDPOINT=otel-collector:4318
      - METRICS_ENABLED=true
      - METRICS_BACKEND=prometheus
      - PROFILING_ENABLED=true
      - PROFILING_BACKEND=pyroscope
      - PROFILING_ENDPOINT=http://pyroscope:4040

  otel-collector:
    image: otel/opentelemetry-collector:0.91.0
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    command: ["--config=/etc/otel-collector-config.yaml"]
    ports:
      - "4317:4317"
      - "4318:4318"

  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  tempo:
    image: grafana/tempo:2.3.1
    command: ["-config.file=/etc/tempo.yaml"]
    ports:
      - "3200:3200"

  pyroscope:
    image: pyroscope/pyroscope:latest
    ports:
      - "4040:4040"

  grafana:
    image: grafana/grafana:10.2.3
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

## Advanced Configuration

For advanced scenarios, you can fully customize the telemetry configuration:

```go
cfg := config.NewConfig()

// Tracing configuration
cfg.Telemetry.Tracing.Enabled = true
cfg.Telemetry.Tracing.Backend = "otlp"
cfg.Telemetry.Tracing.Endpoint = "tempo-us-central1.grafana.net:4318"
cfg.Telemetry.Tracing.SampleRate = 0.1
cfg.Telemetry.Tracing.Insecure = false
cfg.Telemetry.Tracing.BatchSize = 100
cfg.Telemetry.Tracing.BatchTimeout = 5 * time.Second

// Metrics configuration
cfg.Telemetry.Metrics.Enabled = true
cfg.Telemetry.Metrics.Backend = "prometheus"
cfg.Telemetry.Metrics.Path = "/metrics"
cfg.Telemetry.Metrics.Namespace = "netgex"

// Profiling configuration
cfg.Telemetry.Profiling.Enabled = true
cfg.Telemetry.Profiling.Backend = "pyroscope"
cfg.Telemetry.Profiling.Endpoint = "http://pyroscope:4040"
cfg.Telemetry.Profiling.SampleRate = 1.0
cfg.Telemetry.Profiling.Types = "cpu,heap,goroutine"

srv := server.New(
    server.WithConfig(cfg),
    server.WithTelemetry(),
)
```

## Grafana Cloud Setup

For production use with Grafana Cloud:

```go
srv := server.New(
    server.WithTracingBackend("otlp", "tempo-us-central1.grafana.net:4318"),
    server.WithMetricsBackend("otlp", "prometheus-us-central1.grafana.net:4318"),
    server.WithProfilingBackend("pyroscope", "profiles-prod-006.grafana.net:4040"),
)
```

## Self-hosted Production Setup

For self-hosted production:

```go
srv := server.New(
    server.WithTracingBackend("otlp", "tempo.monitoring:4317"),
    server.WithMetricsBackend("prometheus", ""),
    server.WithProfilingBackend("pyroscope", "http://pyroscope.monitoring:4040"),
)
```

---

For more information on observability best practices, see the documentation in `.private/Netgex/Observability.md`. 