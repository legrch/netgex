package telemetry

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GetUnaryInterceptors returns the unary interceptors for telemetry
func (s *Service) GetUnaryInterceptors() []grpc.UnaryServerInterceptor {
	var interceptors []grpc.UnaryServerInterceptor

	// Add tracing interceptor if enabled
	if s.config.Telemetry.Tracing.Enabled {
		interceptors = append(interceptors, s.TracingUnaryInterceptor())
	}

	// Add metrics interceptor if enabled
	if s.config.Telemetry.Metrics.Enabled && s.config.Telemetry.Metrics.Backend == "prometheus" {
		interceptors = append(interceptors, s.MetricsUnaryInterceptor())
	}

	return interceptors
}

// GetStreamInterceptors returns the stream interceptors for telemetry
func (s *Service) GetStreamInterceptors() []grpc.StreamServerInterceptor {
	var interceptors []grpc.StreamServerInterceptor

	// Add tracing interceptor if enabled
	if s.config.Telemetry.Tracing.Enabled {
		interceptors = append(interceptors, s.TracingStreamInterceptor())
	}

	// Add metrics interceptor if enabled
	if s.config.Telemetry.Metrics.Enabled && s.config.Telemetry.Metrics.Backend == "prometheus" {
		interceptors = append(interceptors, s.MetricsStreamInterceptor())
	}

	return interceptors
}

// TracingUnaryInterceptor creates a gRPC unary interceptor for tracing
func (s *Service) TracingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get tracer
		tracer := otel.Tracer("grpc.server")

		// Extract method name
		methodName := info.FullMethod

		// Start span
		ctx, span := tracer.Start(ctx, methodName,
			trace.WithAttributes(
				attribute.String("rpc.service", s.config.ServiceName),
				attribute.String("rpc.method", methodName),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Handle request
		resp, err := handler(ctx, req)

		// Set status based on error
		if err != nil {
			st, _ := status.FromError(err)
			span.SetStatus(codes.Error, st.Message())
			span.SetAttributes(attribute.String("error.code", st.Code().String()))
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return resp, err
	}
}

// TracingStreamInterceptor creates a gRPC stream interceptor for tracing
func (s *Service) TracingStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Get tracer
		tracer := otel.Tracer("grpc.server")

		// Extract method name
		methodName := info.FullMethod

		// Start span
		ctx, span := tracer.Start(ss.Context(), methodName,
			trace.WithAttributes(
				attribute.String("rpc.service", s.config.ServiceName),
				attribute.String("rpc.method", methodName),
				attribute.Bool("rpc.stream", true),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Wrap server stream to propagate the context
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// Handle stream
		err := handler(srv, wrappedStream)

		// Set status based on error
		if err != nil {
			st, _ := status.FromError(err)
			span.SetStatus(codes.Error, st.Message())
			span.SetAttributes(attribute.String("error.code", st.Code().String()))
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return err
	}
}

// MetricsUnaryInterceptor creates a gRPC unary interceptor for Prometheus metrics
func (s *Service) MetricsUnaryInterceptor() grpc.UnaryServerInterceptor {
	// Initialize metrics
	grpcRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: s.config.Telemetry.Metrics.Namespace,
			Name:      "grpc_requests_total",
			Help:      "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	grpcRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: s.config.Telemetry.Metrics.Namespace,
			Name:      "grpc_request_duration_seconds",
			Help:      "Duration of gRPC requests in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"method"},
	)

	// Register metrics
	prometheus.MustRegister(
		grpcRequestsTotal,
		grpcRequestDuration,
	)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(startTime).Seconds()

		// Record metrics
		statusCode := "success"
		if err != nil {
			statusCode = status.Code(err).String()
		}

		grpcRequestsTotal.WithLabelValues(info.FullMethod, statusCode).Inc()
		grpcRequestDuration.WithLabelValues(info.FullMethod).Observe(duration)

		return resp, err
	}
}

// MetricsStreamInterceptor creates a gRPC stream interceptor for Prometheus metrics
func (s *Service) MetricsStreamInterceptor() grpc.StreamServerInterceptor {
	// Initialize metrics
	grpcStreamRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: s.config.Telemetry.Metrics.Namespace,
			Name:      "grpc_stream_requests_total",
			Help:      "Total number of gRPC stream requests",
		},
		[]string{"method", "status"},
	)

	grpcStreamDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: s.config.Telemetry.Metrics.Namespace,
			Name:      "grpc_stream_duration_seconds",
			Help:      "Duration of gRPC streams in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10, 30, 60, 120, 300, 600},
		},
		[]string{"method"},
	)

	// Register metrics
	prometheus.MustRegister(
		grpcStreamRequestsTotal,
		grpcStreamDuration,
	)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		err := handler(srv, ss)
		duration := time.Since(startTime).Seconds()

		// Record metrics
		statusCode := "success"
		if err != nil {
			statusCode = status.Code(err).String()
		}

		grpcStreamRequestsTotal.WithLabelValues(info.FullMethod, statusCode).Inc()
		grpcStreamDuration.WithLabelValues(info.FullMethod).Observe(duration)

		return err
	}
}

// wrappedServerStream wraps grpc.ServerStream to modify the context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the wrapped context
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
