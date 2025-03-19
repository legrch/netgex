package server

import (
	"github.com/legrch/netgex/config"
	"github.com/legrch/netgex/service"
	"log/slog"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

// Option is a function that configures a Server
type Option func(*Server)

// WithLogger sets the logger for the Server
func WithLogger(logger *slog.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithConfig sets the configuration for the Server
func WithConfig(config *config.Config) Option {
	return func(s *Server) {
		s.cfg = config
	}
}

// WithServices sets the service implementations
func WithServices(services ...service.Registrar) Option {
	return func(s *Server) {
		s.services = services
	}
}

// WithProcesses adds additional processes to the Server
func WithProcesses(processes ...Process) Option {
	return func(s *Server) {
		s.processes = append(s.processes, processes...)
	}
}

// WithGRPCServerOptions sets additional options for the gRPC server
func WithGRPCServerOptions(options ...grpc.ServerOption) Option {
	return func(s *Server) {
		s.grpcServerOptions = options
	}
}

// WithGRPCUnaryInterceptors sets the unary interceptors for the gRPC server
func WithGRPCUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(s *Server) {
		s.grpcUnaryServerInterceptors = interceptors
	}
}

// WithGRPCStreamInterceptors sets the stream interceptors for the gRPC server
func WithGRPCStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(s *Server) {
		s.grpcStreamServerInterceptors = interceptors
	}
}

// WithGatewayMuxOptions sets the ServeMux options for the gateway server
func WithGatewayMuxOptions(options ...runtime.ServeMuxOption) Option {
	return func(s *Server) {
		s.gwServerMuxOptions = options
	}
}

// WithGatewayCORS enables CORS with the specified options for the gateway
func WithGatewayCORS(options cors.Options) Option {
	return func(s *Server) {
		s.gwCORSEnabled = true
		s.gwCORSOptions = options
	}
}

// Configuration shortcuts for common config fields

// WithGRPCAddress sets the gRPC server address
func WithGRPCAddress(address string) Option {
	return func(s *Server) {
		s.cfg.GRPCAddress = address
	}
}

// WithHTTPAddress sets the HTTP server address
func WithHTTPAddress(address string) Option {
	return func(s *Server) {
		s.cfg.HTTPAddress = address
	}
}

// WithMetricsAddress sets the metrics server address
func WithMetricsAddress(address string) Option {
	return func(s *Server) {
		s.cfg.MetricsAddress = address
	}
}

// WithPprofAddress sets the pprof server address
func WithPprofAddress(address string) Option {
	return func(s *Server) {
		s.cfg.PprofAddress = address
	}
}

// WithCloseTimeout sets the timeout for graceful shutdown
func WithCloseTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.cfg.CloseTimeout = timeout
	}
}

// WithReflection enables or disables gRPC reflection
func WithReflection(enabled bool) Option {
	return func(s *Server) {
		s.cfg.ReflectionEnabled = enabled
	}
}

// WithHealthCheck enables or disables health checks
func WithHealthCheck(enabled bool) Option {
	return func(s *Server) {
		s.cfg.HealthCheckEnabled = enabled
	}
}

// WithSwaggerDir sets the directory containing swagger files
func WithSwaggerDir(dir string) Option {
	return func(s *Server) {
		s.cfg.SwaggerEnabled = true
		s.cfg.SwaggerDir = dir
	}
}

// WithSwaggerBasePath sets the base path for swagger UI
func WithSwaggerBasePath(path string) Option {
	return func(s *Server) {
		s.cfg.SwaggerEnabled = true
		s.cfg.SwaggerBasePath = path
	}
}
