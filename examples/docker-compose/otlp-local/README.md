# NetGeX OTLP Local Telemetry Example

This example demonstrates how to use the NetGeX telemetry system with modern OpenTelemetry backends:
- **OpenTelemetry Collector** for data collection and routing
- **Tempo** for distributed tracing
- **Prometheus** for metrics
- **Pyroscope** for continuous profiling
- **Grafana** for visualization

## Running the Example

### 1. Start the Infrastructure

First, start the required infrastructure using Docker Compose:

```bash
docker compose up -d
```

This will start:
- OpenTelemetry Collector: localhost:4318 (HTTP)
- Tempo: http://localhost:3200
- Prometheus: http://localhost:9090
- Pyroscope: http://localhost:4040
- Grafana: http://localhost:3000

### 2. Run the Example App

Run the example app directly:

```bash
go run app.go
```

The app will:
- Generate traces and send them to Tempo via OTLP
- Expose metrics for Prometheus to scrape
- Send continuous profiling data to Pyroscope

## Endpoints

- HTTP API: http://localhost:8080
- gRPC: localhost:9090
- Metrics: http://localhost:9091/metrics

## Exploring the Data

### Traces in Grafana Tempo
1. Open Grafana at http://localhost:3000
2. Go to Explore and select the Tempo data source
3. Click "Search" to find traces

### Metrics in Prometheus
1. Open Prometheus at http://localhost:9090
2. Go to Graph and query metrics like `go_goroutines`

### Metrics in Grafana
1. Open Grafana at http://localhost:3000
2. Go to Explore and select the Prometheus data source
3. Query metrics like `go_goroutines`

### Profiles in Pyroscope
1. Open Pyroscope at http://localhost:4040
2. The service should appear in the application dropdown 
3. Select different profile types (CPU, Heap, etc.)

## OpenTelemetry Collector

The OpenTelemetry Collector is configured to:
1. Receive data via OTLP (HTTP)
2. Export traces to Tempo
3. Export metrics to Prometheus

Check the `configs/otel-collector-config.yaml` file for details.

## Configuration

The example app uses the NetGeX configuration system. You can override settings using environment variables:

```bash
TRACING_ENDPOINT=localhost:4318 PROFILING_ENDPOINT=http://localhost:4040 go run app.go
```

See `app.go` for the default configuration and available settings. 