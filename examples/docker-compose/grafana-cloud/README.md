# NetGeX Grafana Cloud Telemetry Example

This example demonstrates how to use the NetGeX telemetry system with Grafana Cloud:
- **Grafana Tempo** for distributed tracing
- **Grafana Prometheus** for metrics
- **Grafana Pyroscope** for continuous profiling
- **Grafana Dashboard** for visualization

## Prerequisites

You need a Grafana Cloud account and API key to run this example. Get one at [grafana.com](https://grafana.com/products/cloud/).

## Running the Example

### 1. Set Your Grafana Cloud API Key

Set the API key as an environment variable:

```bash
export GRAFANA_CLOUD_API_KEY="your-api-key"
```

### 2. Start Local Grafana (Optional)

If you want to view your Grafana Cloud data locally:

```bash
docker compose up -d
```

This will start:
- Grafana: http://localhost:3000

### 3. Run the Example App

Run the example app directly:

```bash
go run app.go
```

The app will:
- Generate traces and send them to Grafana Tempo via OTLP
- Send metrics to Grafana Prometheus via OTLP
- Send continuous profiling data to Grafana Pyroscope

## Endpoints

- HTTP API: http://localhost:8080
- gRPC: localhost:9090
- Metrics: http://localhost:9091/metrics

## Exploring the Data

### Grafana Cloud
1. Log in to your Grafana Cloud instance
2. Navigate to Explore to view:
   - Traces in Tempo
   - Metrics in Prometheus
   - Profiles in Pyroscope

### Local Grafana (if using)
1. Open Grafana at http://localhost:3000
2. The Grafana Cloud data sources should be pre-configured
3. Navigate to Explore to view your data

## Configuration

The example app uses the NetGeX configuration system. You can override settings using environment variables:

```bash
SERVICE_NAME=my-cloud-app ENVIRONMENT=staging go run app.go
```

### Important Environment Variables

- `GRAFANA_CLOUD_API_KEY`: Your Grafana Cloud API key (required)
- `TRACING_ENDPOINT`: Grafana Cloud Tempo endpoint 
- `METRICS_ENDPOINT`: Grafana Cloud Prometheus endpoint
- `PROFILING_ENDPOINT`: Grafana Cloud Pyroscope endpoint

See `app.go` for the default configuration and available settings.

## Grafana Cloud Regions

The default endpoints in this example are for the US Central region. If your Grafana Cloud instance is in a different region, update the endpoints accordingly:

- US: `tempo-us-central1.grafana.net:4318`
- EU: `tempo-eu-west-0.grafana.net:4318`
- APAC: `tempo-ap-southeast-0.grafana.net:4318`

Similar patterns apply for Prometheus and Pyroscope endpoints. 