package gateway

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	mocksvc "github.com/legrch/netgex/internal/mocks/service"
	"github.com/legrch/netgex/service"
	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// mockServiceRegistrar implements service.Registrar for testing
type mockServiceRegistrar struct {
	mock.Mock
}

func (m *mockServiceRegistrar) RegisterGRPC(srv *grpc.Server) {
	m.Called(srv)
}

func (m *mockServiceRegistrar) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	args := m.Called(ctx, mux, endpoint, opts)
	return args.Error(0)
}

func TestNewServer(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 5 * time.Second
	grpcAddress := ":50051"
	httpAddress := ":8080"

	// Act
	server := NewServer(logger, closeTimeout, grpcAddress, httpAddress)

	// Assert
	assert.NotNil(t, server)
	assert.Equal(t, logger, server.logger)
	assert.Equal(t, closeTimeout, server.closeTimeout)
	assert.Equal(t, grpcAddress, server.grpcAddress)
	assert.Equal(t, httpAddress, server.httpAddress)
	assert.Equal(t, httpAddress, server.server.Addr)
	assert.NotNil(t, server.jsonConfig)
}

func TestWithServices(t *testing.T) {
	// Arrange
	server := &Server{
		registrars: []service.Registrar{},
	}

	// Create mocks
	svc1 := mocksvc.NewRegistrar(t)
	svc2 := mocksvc.NewRegistrar(t)

	// Act
	WithServices(svc1, svc2)(server)

	// Assert
	assert.Len(t, server.registrars, 2)
	assert.Contains(t, server.registrars, svc1)
	assert.Contains(t, server.registrars, svc2)
}

func TestWithMuxOptions(t *testing.T) {
	// Arrange
	srv := &Server{}
	opt1 := runtime.WithMarshalerOption("application/json", &runtime.JSONPb{})
	opt2 := runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	})

	// Act
	opt := WithMuxOptions(opt1, opt2)
	opt(srv)

	// Assert
	assert.Len(t, srv.muxOptions, 2)
}

func TestWithIncomingHeaderMatcher(t *testing.T) {
	// Arrange
	srv := &Server{}
	matcher := func(s string) (string, bool) {
		return s, true
	}

	// Act
	opt := WithIncomingHeaderMatcher(matcher)
	opt(srv)

	// Assert
	assert.NotNil(t, srv.incomingHeaderMatcher)
	result, ok := srv.incomingHeaderMatcher("test")
	assert.Equal(t, "test", result)
	assert.True(t, ok)
}

func TestWithOutgoingHeaderMatcher(t *testing.T) {
	// Arrange
	srv := &Server{}
	matcher := func(s string) (string, bool) {
		return s, true
	}

	// Act
	opt := WithOutgoingHeaderMatcher(matcher)
	opt(srv)

	// Assert
	assert.NotNil(t, srv.outgoingHeaderMatcher)
	result, ok := srv.outgoingHeaderMatcher("test")
	assert.Equal(t, "test", result)
	assert.True(t, ok)
}

func TestWithCORS(t *testing.T) {
	// Arrange
	srv := &Server{}
	corsOptions := &cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	}

	// Act
	opt := WithCORS(corsOptions)
	opt(srv)

	// Assert
	assert.True(t, srv.corsEnabled)
	assert.Equal(t, *corsOptions, srv.corsOptions)
}

func TestWithPprof(t *testing.T) {
	// Arrange
	srv := &Server{}

	// Act
	opt := WithPprof(true)
	opt(srv)

	// Assert
	assert.True(t, srv.pprofEnabled)
}

func TestWithSwagger(t *testing.T) {
	// Arrange
	srv := &Server{}
	dir := "/api/swagger"
	basePath := "/api/v1"

	// Act
	opt := WithSwagger(dir, basePath)
	opt(srv)

	// Assert
	assert.True(t, srv.swaggerEnabled)
	assert.Equal(t, dir, srv.swaggerDir)
	assert.Equal(t, basePath, srv.swaggerBasePath)
}

func TestWithJSONConfig(t *testing.T) {
	// Arrange
	srv := &Server{}
	config := &JSONConfig{
		UseProtoNames:   false,
		EmitUnpopulated: false,
		UseEnumNumbers:  false,
		AllowPartial:    false,
		Multiline:       false,
		Indent:          "    ",
	}

	// Act
	opt := WithJSONConfig(config)
	opt(srv)

	// Assert
	assert.Equal(t, config, srv.jsonConfig)
}

func TestServer_PreRun(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 5 * time.Second
	grpcAddress := ":50051"
	httpAddress := ":8081"

	// Create server
	srv := NewServer(logger, closeTimeout, grpcAddress, httpAddress)

	// Act
	err := srv.PreRun(context.Background())

	// Assert
	assert.NoError(t, err)
}

func TestServer_Run_RegisterHTTPError(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 5 * time.Second
	grpcAddress := ":50051"
	httpAddress := ":8081"

	// Create mock registrar that returns an error
	registrar := new(mockServiceRegistrar)
	registrar.On("RegisterHTTP", mock.Anything, mock.Anything, grpcAddress, mock.Anything).Return(assert.AnError)

	// Create server
	srv := NewServer(
		logger,
		closeTimeout,
		grpcAddress,
		httpAddress,
		WithServices(registrar),
	)

	// Act
	err := srv.Run(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to register gateway")
	registrar.AssertExpectations(t)
}

func TestServer_Shutdown(t *testing.T) {
	// Arrange
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	closeTimeout := 5 * time.Second
	grpcAddress := ":50051"
	httpAddress := ":8081"

	// Create a test server that we can control
	testServer := &http.Server{}

	// Create gateway server with a mock HTTP server
	srv := NewServer(logger, closeTimeout, grpcAddress, httpAddress)
	srv.server = testServer

	// Act
	err := srv.Shutdown(context.Background())

	// Assert
	assert.NoError(t, err)
}
