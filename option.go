package netgex

import (
	"log/slog"
	"time"

	"github.com/rs/cors"
	"google.golang.org/grpc"

	"github.com/legrch/netgex/internal/gateway"
	"github.com/legrch/netgex/service"
)

// Option is a function that configures a Server
type Option func(*Server)

// WithLogger sets the logger for the Server
func WithLogger(logger *slog.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithCloseTimeout sets the timeout for graceful shutdown
func WithCloseTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.closeTimeout = timeout
	}
}

// WithGRPCAddress sets the gRPC server address
func WithGRPCAddress(address string) Option {
	return func(s *Server) {
		s.grpcAddress = address
	}
}

// WithHTTPAddress sets the HTTP server address
func WithHTTPAddress(address string) Option {
	return func(s *Server) {
		s.httpAddress = address
	}
}

// WithMetricsAddress sets the metrics server address
func WithMetricsAddress(address string) Option {
	return func(s *Server) {
		s.metricsAddress = address
	}
}

// WithPprofAddress sets the pprof server address
func WithPprofAddress(address string) Option {
	return func(s *Server) {
		s.pprofAddress = address
	}
}

// WithSwaggerDir sets the directory containing swagger files
func WithSwaggerDir(dir string) Option {
	return func(s *Server) {
		s.swaggerDir = dir
	}
}

// WithSwaggerBasePath sets the base path for swagger UI
func WithSwaggerBasePath(path string) Option {
	return func(s *Server) {
		s.swaggerBasePath = path
	}
}

// WithReflection enables or disables gRPC reflection
func WithReflection(enabled bool) Option {
	return func(s *Server) {
		s.reflection = enabled
	}
}

// WithHealthCheck enables or disables health checks
func WithHealthCheck(enabled bool) Option {
	return func(s *Server) {
		s.healthCheck = enabled
	}
}

// WithRegistrars sets the service registrars
func WithRegistrars(registrars ...service.Registrar) Option {
	return func(s *Server) {
		s.registrars = registrars
	}
}

// WithProcesses adds additional processes to the Server
func WithProcesses(processes ...Process) Option {
	return func(s *Server) {
		s.processes = append(s.processes, processes...)
	}
}

// WithUnaryInterceptors sets the unary interceptors for the gRPC server
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(s *Server) {
		s.unaryInterceptors = interceptors
	}
}

// WithStreamInterceptors sets the stream interceptors for the gRPC server
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(s *Server) {
		s.streamInterceptors = interceptors
	}
}

// WithCORS enables CORS with the specified options
func WithCORS(options *cors.Options) Option {
	return func(s *Server) {
		s.corsEnabled = true
		s.corsOptions = *options
	}
}

// WithJSONConfig sets the JSON configuration for the gateway
func WithJSONConfig(config *gateway.JSONConfig) Option {
	return func(s *Server) {
		s.jsonConfig = config
	}
}

// WithAppName sets the application name for the splash screen
func WithAppName(name string) Option {
	return func(s *Server) {
		s.appName = name
	}
}

// WithAppVersion sets the application version for the splash screen
func WithAppVersion(version string) Option {
	return func(s *Server) {
		s.appVersion = version
	}
}

// JSONConfigFromEnv creates a JSONConfig from environment variables
func JSONConfigFromEnv() *gateway.JSONConfig {
	return &gateway.JSONConfig{
		UseProtoNames:   getEnvBool("JSON_USE_PROTO_NAMES"),
		EmitUnpopulated: getEnvBool("JSON_EMIT_UNPOPULATED"),
		UseEnumNumbers:  getEnvBool("JSON_USE_ENUM_NUMBERS"),
		AllowPartial:    getEnvBool("JSON_ALLOW_PARTIAL"),
		Multiline:       getEnvBool("JSON_MULTILINE"),
		Indent:          getEnv("JSON_INDENT", "  "),
	}
}

// WithJSONConfigFromEnv sets the JSON configuration from environment variables
func WithJSONConfigFromEnv() Option {
	return func(s *Server) {
		s.jsonConfig = JSONConfigFromEnv()
	}
}

// WithJSONUseProtoNames sets whether to use proto names in JSON output
func WithJSONUseProtoNames(useProtoNames bool) Option {
	return func(s *Server) {
		s.jsonConfig.UseProtoNames = useProtoNames
	}
}

// WithJSONEmitUnpopulated sets whether to emit unpopulated fields in JSON output
func WithJSONEmitUnpopulated(emitUnpopulated bool) Option {
	return func(s *Server) {
		s.jsonConfig.EmitUnpopulated = emitUnpopulated
	}
}

// WithJSONUseEnumNumbers sets whether to use enum numbers in JSON output
func WithJSONUseEnumNumbers(useEnumNumbers bool) Option {
	return func(s *Server) {
		s.jsonConfig.UseEnumNumbers = useEnumNumbers
	}
}

// WithJSONAllowPartial sets whether to allow partial messages in JSON output
func WithJSONAllowPartial(allowPartial bool) Option {
	return func(s *Server) {
		s.jsonConfig.AllowPartial = allowPartial
	}
}

// WithJSONMultiline sets whether to use multiline formatting in JSON output
func WithJSONMultiline(multiline bool) Option {
	return func(s *Server) {
		s.jsonConfig.Multiline = multiline
	}
}

// WithJSONIndent sets the indentation to use in multiline JSON output
func WithJSONIndent(indent string) Option {
	return func(s *Server) {
		s.jsonConfig.Indent = indent
	}
}
