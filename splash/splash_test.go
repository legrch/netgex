package splash

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSplash(t *testing.T) {
	// Act
	splash := NewSplash()

	// Assert
	hostname, _ := os.Hostname()
	assert.Equal(t, hostname, splash.hostname)
	assert.NotEmpty(t, splash.goVersion)
	assert.Empty(t, splash.features)
}

func TestWithOptions(t *testing.T) {
	// Arrange & Act
	splash := NewSplash(
		WithGRPCAddress(":50051"),
		WithHTTPAddress(":8081"),
		WithMetricsAddress(":9091"),
		WithPprofAddress(":6060"),
		WithSwaggerBasePath("/api/v1"),
		WithFeature("Test Feature 1"),
		WithFeature("Test Feature 2"),
	)

	// Assert
	assert.Equal(t, ":50051", splash.grpcAddress)
	assert.Equal(t, ":8081", splash.httpAddress)
	assert.Equal(t, ":9091", splash.metricsAddress)
	assert.Equal(t, ":6060", splash.pprofAddress)
	assert.True(t, splash.swaggerEnabled)
	assert.Equal(t, "/api/v1", splash.swaggerBasePath)
	assert.Contains(t, splash.features, "Test Feature 1")
	assert.Contains(t, splash.features, "Test Feature 2")
}

func TestSplash_String(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		splash   *Splash
		contains []string
		excludes []string
	}{
		{
			name:   "basic splash with hostname and go version",
			splash: NewSplash(),
			contains: []string{
				"Hostname",
				"Go Version",
			},
			excludes: []string{
				"Endpoints",
				"Features",
			},
		},
		{
			name: "splash with endpoints",
			splash: NewSplash(
				WithGRPCAddress(":50051"),
				WithHTTPAddress(":8081"),
			),
			contains: []string{
				"Endpoints",
				"gRPC API: :50051",
				"HTTP API: :8081",
			},
		},
		{
			name: "splash with metrics and pprof",
			splash: NewSplash(
				WithMetricsAddress(":9091"),
				WithPprofAddress(":6060"),
			),
			contains: []string{
				"Endpoints",
				"Metrics: :9091",
				"Profiling: :6060",
			},
		},
		{
			name: "splash with swagger",
			splash: NewSplash(
				WithHTTPAddress(":8081"),
				WithSwaggerBasePath("/api/v1"),
			),
			contains: []string{
				"Swagger UI",
				"http://localhost:8081/swagger",
			},
		},
		{
			name: "splash with features",
			splash: NewSplash(
				WithFeature("Feature 1"),
				WithFeature("Feature 2"),
			),
			contains: []string{
				"Features",
				"Feature 1",
				"Feature 2",
			},
		},
		{
			name: "complete splash",
			splash: NewSplash(
				WithGRPCAddress(":50051"),
				WithHTTPAddress(":8081"),
				WithMetricsAddress(":9091"),
				WithPprofAddress(":6060"),
				WithSwaggerBasePath("/api/v1"),
				WithFeature("Feature 1"),
				WithFeature("Feature 2"),
			),
			contains: []string{
				"Hostname",
				"Go Version",
				"Endpoints",
				"gRPC API: :50051",
				"HTTP API: :8081",
				"Metrics: :9091",
				"Profiling: :6060",
				"Swagger UI",
				"Features",
				"Feature 1",
				"Feature 2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			output := tt.splash.String()

			// Assert
			for _, s := range tt.contains {
				assert.Contains(t, output, s)
			}
			for _, s := range tt.excludes {
				assert.NotContains(t, output, s)
			}
		})
	}
}

func TestSplash_Display(t *testing.T) {
	// We can't easily test the output to stdout, but we can at least make sure
	// the function doesn't panic
	s := NewSplash()
	assert.NotPanics(t, func() {
		s.Display()
	})
}

func TestWithGRPCAddress(t *testing.T) {
	// Arrange
	s := NewSplash()
	opt := WithGRPCAddress(":50051")

	// Act
	opt(s)

	// Assert
	assert.Equal(t, ":50051", s.grpcAddress)
}

func TestWithHTTPAddress(t *testing.T) {
	// Arrange
	s := NewSplash()
	opt := WithHTTPAddress(":8081")

	// Act
	opt(s)

	// Assert
	assert.Equal(t, ":8081", s.httpAddress)
}

func TestWithMetricsAddress(t *testing.T) {
	// Arrange
	s := NewSplash()
	opt := WithMetricsAddress(":9091")

	// Act
	opt(s)

	// Assert
	assert.Equal(t, ":9091", s.metricsAddress)
}

func TestWithPprofAddress(t *testing.T) {
	// Arrange
	s := NewSplash()
	opt := WithPprofAddress(":6060")

	// Act
	opt(s)

	// Assert
	assert.Equal(t, ":6060", s.pprofAddress)
}

func TestWithSwaggerBasePath(t *testing.T) {
	// Arrange
	s := NewSplash()
	opt := WithSwaggerBasePath("/api/v1")

	// Act
	opt(s)

	// Assert
	assert.True(t, s.swaggerEnabled)
	assert.Equal(t, "/api/v1", s.swaggerBasePath)
}

func TestWithFeature(t *testing.T) {
	// Arrange
	s := NewSplash()

	// Act
	opt1 := WithFeature("Feature 1")
	opt1(s)
	opt2 := WithFeature("Feature 2")
	opt2(s)

	// Assert
	assert.Equal(t, []string{"Feature 1", "Feature 2"}, s.features)
}
