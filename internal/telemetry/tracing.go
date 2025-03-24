package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// setupTracing configures distributed tracing based on the provided configuration
func (s *Service) setupTracing(ctx context.Context) error {
	cfg := s.config.Telemetry.Tracing

	if !cfg.Enabled {
		s.logger.Info("tracing is disabled")
		return nil
	}

	// Create resource with service information
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

	var exporter sdktrace.SpanExporter

	// Configure based on backend type
	switch cfg.Backend {
	case "otlp":
		// Create OTLP exporter
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(cfg.Endpoint),
		}

		if cfg.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}

		exp, err := otlptracehttp.New(ctx, opts...)
		if err != nil {
			return fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}
		exporter = exp
		s.logger.Info("initialized OTLP trace exporter", "endpoint", cfg.Endpoint)

	case "jaeger":
		// Create Jaeger exporter
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Endpoint)))
		if err != nil {
			return fmt.Errorf("failed to create Jaeger exporter: %w", err)
		}
		exporter = exp
		s.logger.Info("initialized Jaeger trace exporter", "endpoint", cfg.Endpoint)

	default:
		return fmt.Errorf("unsupported tracing backend: %s", cfg.Backend)
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

	// Set global TracerProvider
	otel.SetTracerProvider(tp)
	s.tracerProvider = tp

	s.logger.Info("tracing initialized successfully",
		"backend", cfg.Backend,
		"sample_rate", cfg.SampleRate)

	return nil
}
