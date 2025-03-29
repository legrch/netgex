package telemetry

import (
	"context"
	"fmt"
	"strings"

	"github.com/legrch/netgex/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// setupOTEL configures the unified OpenTelemetry provider
func (s *Service) setupOTEL(ctx context.Context) error {
	cfg := s.config.Telemetry.OTEL

	if !cfg.Enabled {
		s.logger.Info("OpenTelemetry is disabled")
		return nil
	}

	s.logger.Info("initializing OpenTelemetry provider",
		"endpoint", cfg.Endpoint,
		"protocol", cfg.Protocol,
		"traces_enabled", cfg.TracesEnabled,
		"metrics_enabled", cfg.MetricsEnabled,
		"logs_enabled", cfg.LogsEnabled)

	// Create common resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(s.config.ServiceName),
			semconv.ServiceVersion(s.config.ServiceVersion),
			attribute.String("environment", s.config.Environment),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Parse headers for authentication/metadata if provided
	headers := parseHeaders(cfg.Headers)

	// Set up tracing if enabled
	if cfg.TracesEnabled {
		traceProvider, err := s.setupOTELTracing(ctx, cfg, res, headers)
		if err != nil {
			return fmt.Errorf("failed to set up OTEL tracing: %w", err)
		}
		s.tracer = traceProvider
	}

	// Set up metrics if enabled
	if cfg.MetricsEnabled {
		meterProvider, err := s.setupOTELMetrics(ctx, cfg, res, headers)
		if err != nil {
			return fmt.Errorf("failed to set up OTEL metrics: %w", err)
		}
		s.meter = meterProvider
	}

	// TODO: Set up logging if enabled when OTLP logging is fully supported

	s.logger.Info("OpenTelemetry initialized successfully")
	return nil
}

// setupOTELTracing configures OpenTelemetry tracing
func (s *Service) setupOTELTracing(
	ctx context.Context,
	cfg config.OTELConfig,
	res *resource.Resource,
	headers map[string]string,
) (*sdktrace.TracerProvider, error) {
	// Create HTTP exporter as the default
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.Endpoint),
	}

	if cfg.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	// Add headers if provided
	if len(headers) > 0 {
		opts = append(opts, otlptracehttp.WithHeaders(headers))
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP trace exporter: %w", err)
	}

	// Create TracerProvider with the exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithMaxExportBatchSize(cfg.BatchSize),
			sdktrace.WithBatchTimeout(cfg.BatchTimeout),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRate)),
	)

	// Set global TracerProvider and propagators
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	s.logger.Info("OTLP tracing initialized",
		"endpoint", cfg.Endpoint,
		"sample_rate", cfg.SampleRate)

	return tp, nil
}

// setupOTELMetrics configures OpenTelemetry metrics
func (s *Service) setupOTELMetrics(
	ctx context.Context,
	cfg config.OTELConfig,
	res *resource.Resource,
	headers map[string]string,
) (*metric.MeterProvider, error) {
	// Create HTTP exporter as the default
	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(cfg.Endpoint),
	}

	if cfg.Insecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}

	// Add headers if provided
	if len(headers) > 0 {
		opts = append(opts, otlpmetrichttp.WithHeaders(headers))
	}

	exp, err := otlpmetrichttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP metric exporter: %w", err)
	}

	reader := metric.NewPeriodicReader(exp)

	// Create MeterProvider
	mp := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(res),
	)

	// Set global MeterProvider
	otel.SetMeterProvider(mp)

	s.logger.Info("OTLP metrics initialized",
		"endpoint", cfg.Endpoint)

	return mp, nil
}

// parseHeaders parses a comma-separated list of key=value pairs into a map
func parseHeaders(headerStr string) map[string]string {
	headers := make(map[string]string)
	if headerStr == "" {
		return headers
	}

	pairs := strings.Split(headerStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	return headers
}
