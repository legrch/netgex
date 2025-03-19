package server

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/legrch/netgex/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockLogger is a mock implementation of the slog.Logger
type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Info(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *mockLogger) Error(msg string, args ...any) {
	m.Called(msg, args)
}

// mockProcess is a mock implementation of the Process interface
type mockProcessWithExpectations struct {
	mock.Mock
}

func (m *mockProcessWithExpectations) PreRun(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockProcessWithExpectations) Run(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockProcessWithExpectations) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestNewServer(t *testing.T) {
	// Act
	s := NewServer()

	// Assert
	assert.NotNil(t, s)
	assert.NotNil(t, s.cfg)
	assert.Empty(t, s.processes)
}

func TestNewServer_WithOptions(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := &config.Config{
		GRPCAddress: ":50051",
		HTTPAddress: ":8081",
	}
	mockProcess := &mockProcess{}

	// Act
	s := NewServer(
		WithLogger(logger),
		WithConfig(cfg),
		WithProcesses(mockProcess),
	)

	// Assert
	assert.Equal(t, logger, s.logger)
	assert.Equal(t, cfg, s.cfg)
	assert.Len(t, s.processes, 1)
	assert.Contains(t, s.processes, mockProcess)
}

func TestServer_Run_Success(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	process1 := new(mockProcessWithExpectations)
	process1.On("PreRun", mock.Anything).Return(nil)
	process1.On("Run", mock.Anything).Return(nil)
	process1.On("Shutdown", mock.Anything).Return(nil)

	process2 := new(mockProcessWithExpectations)
	process2.On("PreRun", mock.Anything).Return(nil)
	process2.On("Run", mock.Anything).Return(nil)
	process2.On("Shutdown", mock.Anything).Return(nil)

	// Create server with mocked processes
	s := NewServer(
		WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))),
		WithProcesses(process1, process2),
	)

	// Cancel context after a delay to simulate graceful shutdown
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	// Act
	err := s.Run(ctx)

	// Assert
	assert.NoError(t, err)
	process1.AssertExpectations(t)
	process2.AssertExpectations(t)
}

func TestServer_Run_PreRunError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("prerun error")

	process1 := new(mockProcessWithExpectations)
	process1.On("PreRun", mock.Anything).Return(expectedErr)

	// Create server with mocked processes
	s := NewServer(
		WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))),
		WithProcesses(process1),
	)

	// Act
	err := s.Run(ctx)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pre-run error")
	process1.AssertExpectations(t)
}

func TestServer_Run_ProcessError(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	expectedErr := errors.New("process error")

	process1 := new(mockProcessWithExpectations)
	process1.On("PreRun", mock.Anything).Return(nil)
	process1.On("Run", mock.Anything).Return(expectedErr)
	process1.On("Shutdown", mock.Anything).Return(nil)

	process2 := new(mockProcessWithExpectations)
	process2.On("PreRun", mock.Anything).Return(nil)
	// Run will be called but we don't specify expectations because the goroutine might not start
	process2.On("Shutdown", mock.Anything).Return(nil)

	// Create server with mocked processes
	s := NewServer(
		WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))),
		WithProcesses(process1, process2),
	)

	// Act
	err := s.Run(ctx)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "process error")
	process1.AssertExpectations(t)
}

func TestServer_Run_ShutdownError(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	shutdownErr := errors.New("shutdown error")

	process1 := new(mockProcessWithExpectations)
	process1.On("PreRun", mock.Anything).Return(nil)
	process1.On("Run", mock.Anything).Return(nil)
	process1.On("Shutdown", mock.Anything).Return(shutdownErr)

	process2 := new(mockProcessWithExpectations)
	process2.On("PreRun", mock.Anything).Return(nil)
	process2.On("Run", mock.Anything).Return(nil)
	process2.On("Shutdown", mock.Anything).Return(nil)

	// Create server with mocked processes
	s := NewServer(
		WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))),
		WithProcesses(process1, process2),
	)

	// Cancel context after a delay to simulate graceful shutdown
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	// Act
	err := s.Run(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, shutdownErr, err)
	process1.AssertExpectations(t)
	process2.AssertExpectations(t)
}

func TestServer_DisplaySplash(t *testing.T) {
	// Arrange
	s := NewServer(
		WithGRPCAddress(":50051"),
		WithHTTPAddress(":8081"),
		WithMetricsAddress(":9091"),
		WithPprofAddress(":6060"),
		WithReflection(true),
		WithHealthCheck(true),
		WithSwaggerDir("./api"),
		WithSwaggerBasePath("/api/v1"),
	)

	// Act & Assert - Just make sure it doesn't panic
	require.NotPanics(t, func() {
		s.displaySplash()
	})
}
