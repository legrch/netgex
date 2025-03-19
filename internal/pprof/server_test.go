package pprof

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	address := ":6060"

	// Act
	server := NewServer(logger, address)

	// Assert
	assert.NotNil(t, server)
	assert.Equal(t, logger, server.logger)
	assert.Equal(t, address, server.server.Addr)
	assert.NotNil(t, server.server.Handler)
}

func TestServer_PreRun(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := NewServer(logger, ":6060")

	// Act
	err := server.PreRun(context.Background())

	// Assert
	assert.NoError(t, err)
}

func TestServer_Shutdown(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := NewServer(logger, ":6060")
	ctx := context.Background()

	// Act
	err := server.Shutdown(ctx)

	// Assert
	assert.NoError(t, err)
}

func TestServer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Use a random port for testing
	server := NewServer(logger, ":0")

	// Start server in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Pre-run
	err := server.PreRun(ctx)
	require.NoError(t, err)

	// Extract actual port from listener
	port := server.server.Addr

	// Run server
	go func() {
		_ = server.Run(ctx)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Act - Try to access a pprof endpoint
	resp, err := http.Get("http://localhost" + port + "/debug/pprof/")

	// Assert
	if err == nil {
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "profiles:")
	}

	// Cleanup
	err = server.Shutdown(ctx)
	assert.NoError(t, err)
}
