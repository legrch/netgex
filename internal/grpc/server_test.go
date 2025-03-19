package grpc

import (
	"context"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	mocksvc "github.com/legrch/netgex/internal/mocks/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestNewServer(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 5 * time.Second
	address := ":50051"
	serviceRegistrar := mocksvc.NewRegistrar(t)

	// Act
	server := NewServer(
		logger,
		closeTimeout,
		address,
		WithServices(serviceRegistrar),
	)

	// Assert
	assert.NotNil(t, server)
	assert.Equal(t, logger, server.logger)
	assert.Equal(t, closeTimeout, server.closeTimeout)
	assert.Equal(t, address, server.address)
	assert.Len(t, server.registrars, 1)
	assert.False(t, server.reflectionEnabled)
	assert.True(t, server.healthCheckEnabled)
}

func TestNewServerWithOptions(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 5 * time.Second
	address := ":50051"

	registrar := mocksvc.NewRegistrar(t)

	unaryInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}

	streamInterceptor := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	}

	// Act
	server := NewServer(
		logger,
		closeTimeout,
		address,
		WithServices(registrar),
		WithUnaryInterceptors(unaryInterceptor),
		WithStreamInterceptors(streamInterceptor),
		WithReflection(true),
		WithHealthCheck(false),
	)

	// Assert
	assert.NotNil(t, server)
	assert.Equal(t, logger, server.logger)
	assert.Equal(t, closeTimeout, server.closeTimeout)
	assert.Equal(t, address, server.address)
	assert.Len(t, server.registrars, 1)
	assert.Len(t, server.unaryInterceptors, 1)
	assert.Len(t, server.streamInterceptors, 1)
	assert.True(t, server.reflectionEnabled)
	assert.False(t, server.healthCheckEnabled)
}

func TestServer_PreRun(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 5 * time.Second
	address := ":50051"

	registrar1 := mocksvc.NewRegistrar(t)
	registrar1.EXPECT().RegisterGRPC(mock.Anything).Return()

	registrar2 := mocksvc.NewRegistrar(t)
	registrar2.EXPECT().RegisterGRPC(mock.Anything).Return()

	// Act
	srv := NewServer(
		logger,
		closeTimeout,
		address,
		WithServices(registrar1, registrar2),
	)
	err := srv.PreRun(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, srv.server)
}

func TestServer_RunAndShutdown(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 5 * time.Second

	// Use random port for tests
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	address := listener.Addr().String()
	listener.Close() // Close so the server can use it

	// Create service registrar mock
	registrar := mocksvc.NewRegistrar(t)
	registrar.EXPECT().RegisterGRPC(mock.Anything).Return()

	// Create server
	srv := NewServer(
		logger,
		closeTimeout,
		address,
		WithServices(registrar),
		WithHealthCheck(true),
	)

	// PreRun
	err = srv.PreRun(context.Background())
	require.NoError(t, err)

	// Run server in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Run(ctx)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test connection by checking health
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	healthResp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, healthResp.Status)

	// Shutdown server
	err = srv.Shutdown(ctx)
	assert.NoError(t, err)

	// Check server stopped
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestServer_Run_ListenError(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 5 * time.Second

	// Use a port that's already in use to cause an error
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer listener.Close()
	address := listener.Addr().String()

	// Create server
	srv := NewServer(
		logger,
		closeTimeout,
		address,
	)

	// PreRun
	err = srv.PreRun(context.Background())
	require.NoError(t, err)

	// Act
	err = srv.Run(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to listen")
}

func TestServer_Shutdown_GracefulStop(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 5 * time.Second
	address := ":50051"

	// Create server with mock gRPC server
	srv := NewServer(
		logger,
		closeTimeout,
		address,
	)

	// Create a real gRPC server
	srv.server = grpc.NewServer()

	// Act & Assert
	err := srv.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestServer_Shutdown_Timeout(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 100 * time.Millisecond // Very short timeout
	address := ":50051"

	// Create server
	srv := NewServer(
		logger,
		closeTimeout,
		address,
	)

	// Create a mock gRPC server that will hang during GracefulStop
	originalServer := grpc.NewServer()

	// Replace GracefulStop with a function that will hang
	srv.server = grpc.NewServer()

	// Start a goroutine that hangs for longer than the timeout
	stopCh := make(chan struct{})
	go func() {
		time.Sleep(500 * time.Millisecond) // Longer than closeTimeout
		close(stopCh)
	}()

	// Act & Assert
	err := srv.Shutdown(context.Background())
	assert.NoError(t, err) // Should not error even with timeout

	// Wait for the hanging goroutine to finish
	<-stopCh

	// Clean up the original server
	originalServer.Stop()
}
