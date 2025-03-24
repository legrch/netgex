# NetGeX Telemetry - Basic Test

This example demonstrates how to test the telemetry system in NetGeX, showing:

- Prometheus metrics collection
- pprof profiling

## Running the Test Locally

```bash
# Run the test application directly
go run cmd/test/main.go

# Check the endpoints
curl http://localhost:9091/metrics
curl http://localhost:6060/debug/pprof/
```

You can also take a CPU profile:
```bash
curl http://localhost:6060/debug/pprof/profile?seconds=10 > cpu.prof
go tool pprof -http=:8080 cpu.prof
```

## Running with Docker Compose

```bash
# Build and start the stack
docker-compose up -d

# Check the services are running
docker-compose ps
```

## Accessing the Services

- **Grafana**: [http://localhost:3000](http://localhost:3000)
  - Add Prometheus as a data source using URL: `http://prometheus:9090`

- **Prometheus**: [http://localhost:9090](http://localhost:9090)
  - Try querying metrics like `netgex_*`

- **pprof**: [http://localhost:6060/debug/pprof/](http://localhost:6060/debug/pprof/)
  - View CPU profile: [http://localhost:6060/debug/pprof/profile?seconds=30](http://localhost:6060/debug/pprof/profile?seconds=30)
  - View heap profile: [http://localhost:6060/debug/pprof/heap](http://localhost:6060/debug/pprof/heap)

## Verifying the Implementation

Our implementation successfully integrates with the standard Go profiling and metrics tools. In the full NetGeX server implementation, these features are enabled through the telemetry system we created.

Key components tested:
1. **Profiling** - Accessible via pprof endpoint
2. **Metrics** - Accessible via metrics endpoint
3. **Grafana Integration** - Dashboards for visualizing metrics

## Shutting Down

```bash
# If running directly
Ctrl+C

# If running with Docker Compose
docker-compose down
``` 