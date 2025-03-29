# NetGeX Legacy Telemetry Example

This example demonstrates how to use the NetGeX telemetry system with legacy backends:
- **Jaeger** for distributed tracing
- **Prometheus** for metrics
- **pprof** for profiling

## Running the Example

### 1. Start the Infrastructure

First, start the required infrastructure using Docker Compose:

```bash
docker compose up -d
```

This will start:
- Jaeger UI: http://localhost:16686
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000

### 2. Run the Example App

Run the example app directly:

```bash
go run app.go
```

The app will:
- Generate traces and send them to Jaeger
- Expose metrics for Prometheus to scrape
- Expose pprof endpoints for profiling

## Endpoints

The app uses non-standard ports to avoid conflicts:
- HTTP API: http://localhost:8082
- gRPC: localhost:50052
- Metrics: http://localhost:9093/metrics
- pprof: http://localhost:6062/debug/pprof/

## Exploring the Data

### Traces in Jaeger
1. Open Jaeger UI at http://localhost:16686
2. Select the "netgex-legacy-example" service
3. Click "Find Traces" to view captured traces

### Metrics in Prometheus
1. Open Prometheus at http://localhost:9090
2. Go to Graph and query metrics like `go_goroutines`

### Metrics in Grafana
1. Open Grafana at http://localhost:3000
2. The Prometheus data source should be pre-configured
3. Create dashboards to visualize your metrics

### Profiling with pprof
1. Capture a 30-second CPU profile:
```bash
curl http://localhost:6062/debug/pprof/profile?seconds=30 > cpu.prof
```

2. Analyze the profile:
```bash
go tool pprof -http=:8080 cpu.prof
```

## Configuration

The example app uses the NetGeX configuration system. You can override settings using environment variables:

```bash
TRACING_ENDPOINT=http://localhost:14268/api/traces go run app.go
```

See `app.go` for the default configuration and available settings. 