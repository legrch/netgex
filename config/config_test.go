package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	// Act
	cfg := NewConfig()

	// Assert
	assert.Equal(t, "info", cfg.LogLevel, "default log level should be 'info'")
	assert.Equal(t, 10*time.Second, cfg.CloseTimeout, "default close timeout should be 10s")
	assert.Equal(t, ":9090", cfg.GRPCAddress, "default gRPC address should be ':9090'")
	assert.Equal(t, ":8080", cfg.HTTPAddress, "default HTTP address should be ':8080'")
	assert.Equal(t, ":9091", cfg.MetricsAddress, "default metrics address should be ':9091'")
	assert.Equal(t, ":6060", cfg.PprofAddress, "default pprof address should be ':6060'")
	assert.True(t, cfg.ReflectionEnabled, "reflection should be enabled by default")
	assert.True(t, cfg.HealthCheckEnabled, "health check should be enabled by default")
	assert.True(t, cfg.SwaggerEnabled, "swagger should be enabled by default")
	assert.Equal(t, "./api", cfg.SwaggerDir, "default swagger dir should be './api'")
	assert.Equal(t, "/", cfg.SwaggerBasePath, "default swagger base path should be '/'")
}

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		validate func(*testing.T, *Config)
	}{
		{
			name:    "default values when env vars not set",
			envVars: map[string]string{},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "info", cfg.LogLevel)
				assert.Equal(t, 10*time.Second, cfg.CloseTimeout)
				assert.Equal(t, ":9090", cfg.GRPCAddress)
			},
		},
		{
			name: "custom values from env vars",
			envVars: map[string]string{
				"TEST_LOG_LEVEL":       "debug",
				"TEST_CLOSE_TIMEOUT":   "5s",
				"TEST_GRPC_ADDRESS":    ":50051",
				"TEST_HTTP_ADDRESS":    ":8081",
				"TEST_SWAGGER_DIR":     "/custom/swagger",
				"TEST_SWAGGER_ENABLED": "false",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "debug", cfg.LogLevel)
				assert.Equal(t, 5*time.Second, cfg.CloseTimeout)
				assert.Equal(t, ":50051", cfg.GRPCAddress)
				assert.Equal(t, ":8081", cfg.HTTPAddress)
				assert.Equal(t, "/custom/swagger", cfg.SwaggerDir)
				assert.False(t, cfg.SwaggerEnabled)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			// Act
			cfg, err := LoadFromEnv("TEST")

			// Assert
			require.NoError(t, err)
			tt.validate(t, cfg)
		})
	}
}

func TestLoadFromEnv_InvalidFormat(t *testing.T) {
	// Setup environment with invalid format
	os.Setenv("TEST_CLOSE_TIMEOUT", "invalid")

	// Act
	_, err := LoadFromEnv("TEST")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CLOSE_TIMEOUT")
}
