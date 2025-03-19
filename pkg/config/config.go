package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the comprehensive configuration for the netgex.Server
type Config struct {
	// Core settings
	LogLevel     string        `envconfig:"LOG_LEVEL" default:"info"`
	CloseTimeout time.Duration `envconfig:"CLOSE_TIMEOUT" default:"10s"`
	AppName      string        `envconfig:"APP_NAME" default:"Service"`
	AppVersion   string        `envconfig:"APP_VERSION" default:"dev"`

	// Server addresses
	GRPCAddress    string `envconfig:"GRPC_ADDRESS" default:":9090"`
	HTTPAddress    string `envconfig:"HTTP_ADDRESS" default:":8080"`
	MetricsAddress string `envconfig:"METRICS_ADDRESS" default:":9091"`
	PprofAddress   string `envconfig:"PPROF_ADDRESS" default:":6060"`

	// Feature flags
	ReflectionEnabled  bool `envconfig:"REFLECTION_ENABLED" default:"true"`
	HealthCheckEnabled bool `envconfig:"HEALTH_CHECK_ENABLED" default:"true"`

	// Swagger configuration
	SwaggerDir      string `envconfig:"SWAGGER_DIR" default:"./api"`
	SwaggerBasePath string `envconfig:"SWAGGER_BASE_PATH" default:"/"`
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		LogLevel:           "info",
		CloseTimeout:       10 * time.Second,
		AppName:            "Service",
		AppVersion:         "dev",
		GRPCAddress:        ":9090",
		HTTPAddress:        ":8080",
		MetricsAddress:     ":9091",
		PprofAddress:       ":6060",
		ReflectionEnabled:  true,
		HealthCheckEnabled: true,
		SwaggerDir:         "./api",
		SwaggerBasePath:    "/",
	}
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv(prefix string) (*Config, error) {
	cfg := NewConfig()
	err := envconfig.Process(prefix, cfg)
	return cfg, err
}
