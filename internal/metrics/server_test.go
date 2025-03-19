package metrics

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	address := ":9091"
	closeTimeout := 5 * time.Second

	// Act
	server := NewServer(logger, address, closeTimeout)

	// Assert
	assert.NotNil(t, server)
	assert.Equal(t, logger, server.logger)
	assert.Equal(t, closeTimeout, server.closeTimeout)
	assert.Equal(t, address, server.server.Addr)
	assert.NotNil(t, server.server.Handler)
}

func TestServer_PreRun(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := NewServer(logger, ":9091", 5*time.Second)

	// Unregister metrics to avoid test pollution
	prometheus.Unregister(AppVersion)

	// Act
	err := server.PreRun(context.Background())

	// Assert
	assert.NoError(t, err)

	// Clean up
	prometheus.Unregister(AppVersion)
}

func TestServer_Shutdown(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := NewServer(logger, ":9091", 5*time.Second)
	ctx := context.Background()

	// Act
	err := server.Shutdown(ctx)

	// Assert
	assert.NoError(t, err)
}

func TestSetAppVersion(t *testing.T) {
	// Arrange - unregister to avoid test pollution
	prometheus.Unregister(AppVersion)
	prometheus.MustRegister(AppVersion)

	// Act
	SetAppVersion("1.0.0")

	// Assert
	count, err := testutil.GatherAndCount(prometheus.DefaultGatherer, "app_version")
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "app_version metric should be registered and have a value")

	// Clean up
	prometheus.Unregister(AppVersion)
}

func TestRegisterAndUnregisterAppMetrics(t *testing.T) {
	// Ensure metric is unregistered at the start
	prometheus.Unregister(AppVersion)

	// Act - Register
	RegisterAppMetrics()

	// Assert - Check if metric exists after registration
	_, err := testutil.GatherAndCount(prometheus.DefaultGatherer, "app_version")
	assert.NoError(t, err, "Metric should be registered")

	// Act - Unregister
	UnregisterAppMetrics()

	// Assert - Check metric no longer exists
	count, _ := testutil.GatherAndCount(prometheus.DefaultGatherer, "app_version")
	assert.Equal(t, 0, count, "Metric should not be registered")
}
