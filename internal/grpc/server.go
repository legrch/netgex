package grpc

import (
	"context"
	"fmt"
	"github.com/legrch/netgex/service"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthGrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Option is a function that configures a Server
type Option func(*Server)

// Server represents a gRPC server
type Server struct {
	logger             *slog.Logger
	server             *grpc.Server
	closeTimeout       time.Duration
	address            string
	registrars         []service.Registrar
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	serverOptions      []grpc.ServerOption
	reflectionEnabled  bool
	healthCheckEnabled bool
}

// NewServer creates a new gRPC server
func NewServer(
	logger *slog.Logger,
	closeTimeout time.Duration,
	address string,
	opts ...Option,
) *Server {
	s := &Server{
		logger:             logger,
		closeTimeout:       closeTimeout,
		address:            address,
		reflectionEnabled:  false,
		healthCheckEnabled: true, // Enable health checks by default
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithServices sets the service registrars for the gRPC server
func WithServices(registrars ...service.Registrar) Option {
	return func(s *Server) {
		s.registrars = append(s.registrars, registrars...)
	}
}

// WithUnaryInterceptors sets the unary interceptors for the gRPC server
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(s *Server) {
		s.unaryInterceptors = append(s.unaryInterceptors, interceptors...)
	}
}

// WithStreamInterceptors sets the stream interceptors for the gRPC server
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(s *Server) {
		s.streamInterceptors = append(s.streamInterceptors, interceptors...)
	}
}

// WithOptions sets additional gRPC server options
func WithOptions(options ...grpc.ServerOption) Option {
	return func(s *Server) {
		s.serverOptions = append(s.serverOptions, options...)
	}
}

// WithReflection enables or disables gRPC reflection
func WithReflection(enabled bool) Option {
	return func(s *Server) {
		s.reflectionEnabled = enabled
	}
}

// WithHealthCheck enables or disables gRPC health checks
func WithHealthCheck(enabled bool) Option {
	return func(s *Server) {
		s.healthCheckEnabled = enabled
	}
}

// PreRun prepares the gRPC server
func (s *Server) PreRun(_ context.Context) error {
	// Prepare server options

	opts := make([]grpc.ServerOption, 0, len(s.serverOptions)+len(s.unaryInterceptors)+len(s.streamInterceptors))
	opts = append(opts, s.serverOptions...)
	opts = append(opts, grpc.ChainUnaryInterceptor(s.unaryInterceptors...), grpc.ChainStreamInterceptor(s.streamInterceptors...))

	// Create gRPC server
	srv := grpc.NewServer(opts...)

	// Register health check service if enabled
	if s.healthCheckEnabled {
		healthServer := health.NewServer()
		healthGrpc.RegisterHealthServer(srv, healthServer)
	}

	// Register all service implementations
	for _, registrar := range s.registrars {
		registrar.RegisterGRPC(srv)
	}

	// Enable reflection if requested
	if s.reflectionEnabled {
		reflection.Register(srv)
	}

	// Store the server
	s.server = srv

	return nil
}

// Run starts the gRPC server
func (s *Server) Run(_ context.Context) error {
	// Create listener
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Start server
	s.logger.Info("starting gRPC server", "address", s.address)
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Shutdown gracefully stops the gRPC server
func (s *Server) Shutdown(_ context.Context) error {
	s.logger.Info("shutting down gRPC server")

	// Create a channel to signal completion
	stopped := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	// Wait for shutdown or timeout
	select {
	case <-stopped:
		s.logger.Info("gRPC server stopped gracefully")
	case <-time.After(s.closeTimeout):
		s.logger.Warn("gRPC server shutdown timed out, forcing stop")
		s.server.Stop()
	}

	return nil
}
