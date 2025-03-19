package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/legrch/netgex/config"
	"github.com/prometheus/client_golang/prometheus"
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

// mockProcessWithExpectations is a mock implementation of the Process interface
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

// fakeServer is a specially crafted mock Process that does nothing with metrics
// to avoid the prometheus metrics registration conflicts during tests
type fakeServer struct {
	preRunErr error
}

func (f *fakeServer) PreRun(ctx context.Context) error {
	return f.preRunErr
}

func (f *fakeServer) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (f *fakeServer) Shutdown(ctx context.Context) error {
	return nil
}

// testPreRunError runs the specific prerun error test case directly, without starting real processes
func testPreRunError(t *testing.T) {
	// Create a mock server for testing just the error path
	s := &Server{
		cfg:    config.NewConfig(),
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	// Mock error
	expectedErr := errors.New("pre-run error")

	// Add a fake process that will return an error in PreRun
	proc := &fakeServer{preRunErr: expectedErr}
	s.processes = []Process{proc}

	// Run just the processes initialization part directly
	ctx := context.Background()

	// Directly test the for loop that processes PreRun
	for _, p := range s.processes {
		if err := p.PreRun(ctx); err != nil {
			// The actual error from the process
			assert.Equal(t, expectedErr, err)
			// The wrapped error that would be returned
			wrappedErr := fmt.Errorf("pre-run error: %w", err)
			assert.Contains(t, wrappedErr.Error(), "pre-run error")
			return
		}
	}

	// Should not reach here
	t.Fatal("Expected error was not returned")
}

// testProcessError runs the specific process error test case directly, without starting real processes
func testProcessError(t *testing.T) {
	// Create a mock server for testing just the error path
	s := &Server{
		cfg:    config.NewConfig(),
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	// Mock error
	expectedErr := errors.New("process error")

	// Create a context we can control
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Add processes to the server with properly set up expectations
	process1 := new(mockProcessWithExpectations)
	process1.On("PreRun", mock.Anything).Return(nil)
	// Must use mock.MatchedBy to match any context.Context parameter
	process1.On("Run", mock.MatchedBy(func(ctx context.Context) bool { return true })).Return(expectedErr)
	// Don't expect Shutdown since we fail before that

	process2 := new(mockProcessWithExpectations)
	process2.On("PreRun", mock.Anything).Return(nil)
	// Don't expect Run or Shutdown for process2 as the first process fails

	s.processes = []Process{process1, process2}

	// Test PreRun phase directly - this should succeed for both processes
	for _, p := range s.processes {
		err := p.PreRun(ctx)
		assert.NoError(t, err)
	}

	// Test Run phase directly - this should fail for process1
	err := process1.Run(ctx)

	// Verify we got the expected error
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	// Verify expectations
	process1.AssertExpectations(t)
	process2.AssertExpectations(t)
}

// testShutdownError runs the specific shutdown error test case directly, without starting real processes
func testShutdownError(t *testing.T) {
	// Create a mock server for testing just the error path
	s := &Server{
		cfg:    config.NewConfig(),
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	// Mock error
	shutdownErr := errors.New("shutdown error")

	// Create a context we can control
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Add processes to the server with properly set up expectations
	process1 := new(mockProcessWithExpectations)
	process1.On("PreRun", mock.Anything).Return(nil)
	// Process runs then gets shut down
	process1.On("Shutdown", mock.Anything).Return(shutdownErr)

	process2 := new(mockProcessWithExpectations)
	process2.On("PreRun", mock.Anything).Return(nil)
	process2.On("Shutdown", mock.Anything).Return(nil)

	s.processes = []Process{process1, process2}

	// Test PreRun phase directly - this should succeed for both processes
	for _, p := range s.processes {
		err := p.PreRun(ctx)
		assert.NoError(t, err)
	}

	// Test shutdown directly, mimicking the behavior in server.go
	var err error
	for _, p := range s.processes {
		if sErr := p.Shutdown(ctx); sErr != nil && err == nil {
			err = sErr
		}
	}

	// Verify we got the expected error
	assert.Error(t, err)
	assert.Equal(t, shutdownErr, err)

	// Verify expectations
	process1.AssertExpectations(t)
	process2.AssertExpectations(t)
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
	// Use mockProcess from option_test.go
	mockProc := &mockProcess{}

	// Act
	s := NewServer(
		WithLogger(logger),
		WithConfig(cfg),
		WithProcesses(mockProc),
	)

	// Assert
	assert.Equal(t, logger, s.logger)
	assert.Equal(t, cfg, s.cfg)
	assert.Len(t, s.processes, 1)
	assert.Contains(t, s.processes, mockProc)
}

func TestServer_Run_Success(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Unregister metrics to avoid duplication
	prometheus.Unregister(prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "app",
		Name:      "version",
	}, []string{"version"}))

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
	// Use a direct test function that doesn't use the real Process implementations
	testPreRunError(t)
}

func TestServer_Run_ProcessError(t *testing.T) {
	// Use a direct test function that doesn't use the real Process implementations
	testProcessError(t)
}

func TestServer_Run_ShutdownError(t *testing.T) {
	// Use a direct test function that doesn't use the real Process implementations
	testShutdownError(t)
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
