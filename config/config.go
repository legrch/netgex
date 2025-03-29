package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the comprehensive configuration for the server.Server
type Config struct {
	// Core settings
	LogLevel     string        `envconfig:"LOG_LEVEL" default:"info"`
	CloseTimeout time.Duration `envconfig:"CLOSE_TIMEOUT" default:"10s"`

	// Server addresses
	GRPCAddress    string `envconfig:"GRPC_ADDRESS" default:":9090"`
	HTTPAddress    string `envconfig:"HTTP_ADDRESS" default:":8080"`
	MetricsAddress string `envconfig:"METRICS_ADDRESS" default:":9091"`
	PprofEnabled   bool   `envconfig:"PPROF_ENABLED" default:"true"`
	PprofAddress   string `envconfig:"PPROF_ADDRESS" default:":6060"`

	// Feature flags
	ReflectionEnabled  bool `envconfig:"REFLECTION_ENABLED" default:"true"`
	HealthCheckEnabled bool `envconfig:"HEALTH_CHECK_ENABLED" default:"true"`

	// Swagger configuration
	SwaggerEnabled  bool   `envconfig:"SWAGGER_ENABLED" default:"true"`
	SwaggerDir      string `envconfig:"SWAGGER_DIR" default:"./api"`
	SwaggerBasePath string `envconfig:"SWAGGER_BASE_PATH" default:"/"`

	// Service information for telemetry
	ServiceName    string `envconfig:"SERVICE_NAME" default:"netgex"`
	ServiceVersion string `envconfig:"SERVICE_VERSION" default:"0.0.0"`
	Environment    string `envconfig:"ENVIRONMENT" default:"development"`

	// Telemetry configuration
	Telemetry TelemetryConfig
}

// TelemetryConfig holds all observability configuration settings
type TelemetryConfig struct {
	// Tracing configuration
	Tracing TracingConfig
	// Metrics configuration
	Metrics MetricsConfig
	// Logging configuration
	Logging LoggingConfig
	// Profiling configuration
	Profiling ProfilingConfig
	// OpenTelemetry configuration (unified approach)
	OTEL OTELConfig
}

// TracingConfig configures distributed tracing
type TracingConfig struct {
	Enabled      bool          `envconfig:"TRACING_ENABLED" default:"false"`
	Backend      string        `envconfig:"TRACING_BACKEND" default:"none"` // "otlp", "jaeger", "none"
	Endpoint     string        `envconfig:"TRACING_ENDPOINT" default:"localhost:4318"`
	Insecure     bool          `envconfig:"TRACING_INSECURE" default:"true"`
	SampleRate   float64       `envconfig:"TRACING_SAMPLE_RATE" default:"1.0"`
	BatchSize    int           `envconfig:"TRACING_BATCH_SIZE" default:"100"`
	BatchTimeout time.Duration `envconfig:"TRACING_BATCH_TIMEOUT" default:"5s"`
}

// MetricsConfig configures metrics collection
type MetricsConfig struct {
	Enabled   bool   `envconfig:"METRICS_ENABLED" default:"false"`
	Backend   string `envconfig:"METRICS_BACKEND" default:"prometheus"` // "prometheus", "otlp", "none"
	Endpoint  string `envconfig:"METRICS_ENDPOINT" default:"localhost:4318"`
	Insecure  bool   `envconfig:"METRICS_INSECURE" default:"true"`
	Path      string `envconfig:"METRICS_PATH" default:"/metrics"`
	Port      int    `envconfig:"METRICS_PORT" default:"9091"`
	Namespace string `envconfig:"METRICS_NAMESPACE" default:"netgex"`
}

// LoggingConfig configures structured logging
type LoggingConfig struct {
	Enabled  bool   `envconfig:"LOGGING_ENABLED" default:"true"`
	Backend  string `envconfig:"LOGGING_BACKEND" default:"stdout"` // "stdout", "otlp", "file", "none"
	Endpoint string `envconfig:"LOGGING_ENDPOINT" default:""`
	Level    string `envconfig:"LOGGING_LEVEL" default:"info"`  // "debug", "info", "warn", "error"
	Format   string `envconfig:"LOGGING_FORMAT" default:"json"` // "json", "text", "console"
	FilePath string `envconfig:"LOGGING_FILE_PATH" default:""`
}

// ProfilingConfig configures continuous profiling
type ProfilingConfig struct {
	Enabled    bool    `envconfig:"PROFILING_ENABLED" default:"false"`
	Backend    string  `envconfig:"PROFILING_BACKEND" default:"none"` // "pyroscope", "otlp", "none"
	Endpoint   string  `envconfig:"PROFILING_ENDPOINT" default:"http://localhost:4040"`
	SampleRate float64 `envconfig:"PROFILING_SAMPLE_RATE" default:"1.0"`
	Types      string  `envconfig:"PROFILING_TYPES" default:"cpu,heap"` // Comma-separated: "cpu,heap,goroutine,mutex,block"
}

// OTELConfig configures OpenTelemetry as a unified observability provider
type OTELConfig struct {
	Enabled  bool   `envconfig:"OTEL_ENABLED" default:"false"`
	Endpoint string `envconfig:"OTEL_ENDPOINT" default:"localhost:4318"`
	Insecure bool   `envconfig:"OTEL_INSECURE" default:"true"`
	Headers  string `envconfig:"OTEL_HEADERS" default:""`      // Format: "key1=value1,key2=value2"
	Protocol string `envconfig:"OTEL_PROTOCOL" default:"http"` // "http" or "grpc"

	// Signal-specific configuration
	TracesEnabled  bool          `envconfig:"OTEL_TRACES_ENABLED" default:"true"`
	MetricsEnabled bool          `envconfig:"OTEL_METRICS_ENABLED" default:"true"`
	LogsEnabled    bool          `envconfig:"OTEL_LOGS_ENABLED" default:"false"`
	SampleRate     float64       `envconfig:"OTEL_SAMPLE_RATE" default:"1.0"`
	BatchSize      int           `envconfig:"OTEL_BATCH_SIZE" default:"100"`
	BatchTimeout   time.Duration `envconfig:"OTEL_BATCH_TIMEOUT" default:"5s"`
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		LogLevel:           "info",
		CloseTimeout:       10 * time.Second,
		GRPCAddress:        ":9090",
		HTTPAddress:        ":8080",
		MetricsAddress:     ":9091",
		PprofEnabled:       true,
		PprofAddress:       ":6060",
		ReflectionEnabled:  true,
		HealthCheckEnabled: true,
		SwaggerEnabled:     true,
		SwaggerDir:         "./api",
		SwaggerBasePath:    "/",
		ServiceName:        "netgex",
		ServiceVersion:     "0.0.0",
		Environment:        "development",
		Telemetry: TelemetryConfig{
			Tracing: TracingConfig{
				Enabled:      false,
				Backend:      "none",
				Endpoint:     "localhost:4318",
				Insecure:     true,
				SampleRate:   1.0,
				BatchSize:    100,
				BatchTimeout: 5 * time.Second,
			},
			Metrics: MetricsConfig{
				Enabled:   false,
				Backend:   "prometheus",
				Endpoint:  "localhost:4318",
				Insecure:  true,
				Path:      "/metrics",
				Port:      9091,
				Namespace: "netgex",
			},
			Logging: LoggingConfig{
				Enabled:  true,
				Backend:  "stdout",
				Level:    "info",
				Format:   "json",
				FilePath: "",
			},
			Profiling: ProfilingConfig{
				Enabled:    false,
				Backend:    "none",
				Endpoint:   "http://localhost:4040",
				SampleRate: 1.0,
				Types:      "cpu,heap",
			},
			OTEL: OTELConfig{
				Enabled:        false,
				Endpoint:       "localhost:4318",
				Insecure:       true,
				Headers:        "",
				Protocol:       "http",
				TracesEnabled:  true,
				MetricsEnabled: true,
				LogsEnabled:    false,
				SampleRate:     1.0,
				BatchSize:      100,
				BatchTimeout:   5 * time.Second,
			},
		},
	}
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv(prefix string) (*Config, error) {
	cfg := NewConfig()
	err := envconfig.Process(prefix, cfg)
	return cfg, err
}
