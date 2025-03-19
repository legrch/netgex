package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/legrch/netgex/config"
	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// mockRegistrar implements service.Registrar
type mockRegistrar struct{}

func (m *mockRegistrar) RegisterGRPC(*grpc.Server) {}
func (m *mockRegistrar) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return nil
}

// mockProcess implements Process
type mockProcess struct{}

func (m *mockProcess) PreRun(ctx context.Context) error   { return nil }
func (m *mockProcess) Run(ctx context.Context) error      { return nil }
func (m *mockProcess) Shutdown(ctx context.Context) error { return nil }

func TestWithLogger(t *testing.T) {
	// Arrange
	s := &Server{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Act
	opt := WithLogger(logger)
	opt(s)

	// Assert
	assert.Equal(t, logger, s.logger)
}

func TestWithConfig(t *testing.T) {
	// Arrange
	s := &Server{}
	cfg := &config.Config{
		GRPCAddress: ":50051",
		HTTPAddress: ":8081",
	}

	// Act
	opt := WithConfig(cfg)
	opt(s)

	// Assert
	assert.Equal(t, cfg, s.cfg)
}

func TestWithServices(t *testing.T) {
	// Arrange
	s := &Server{}
	svc1 := &mockRegistrar{}
	svc2 := &mockRegistrar{}

	// Act
	opt := WithServices(svc1, svc2)
	opt(s)

	// Assert
	assert.Len(t, s.services, 2)
	assert.Contains(t, s.services, svc1)
	assert.Contains(t, s.services, svc2)
}

func TestWithProcesses(t *testing.T) {
	// Arrange
	s := &Server{}
	p1 := &mockProcess{}
	p2 := &mockProcess{}

	// Act
	opt := WithProcesses(p1, p2)
	opt(s)

	// Assert
	assert.Len(t, s.processes, 2)
	assert.Contains(t, s.processes, p1)
	assert.Contains(t, s.processes, p2)
}

func TestWithGRPCServerOptions(t *testing.T) {
	// Arrange
	s := &Server{}
	opt1 := grpc.MaxRecvMsgSize(4096)
	opt2 := grpc.MaxSendMsgSize(4096)

	// Act
	o := WithGRPCServerOptions(opt1, opt2)
	o(s)

	// Assert
	assert.Len(t, s.grpcServerOptions, 2)
}

func TestWithGRPCUnaryInterceptors(t *testing.T) {
	// Arrange
	s := &Server{}
	interceptor1 := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
	interceptor2 := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}

	// Act
	o := WithGRPCUnaryInterceptors(interceptor1, interceptor2)
	o(s)

	// Assert
	assert.Len(t, s.grpcUnaryServerInterceptors, 2)
}

func TestWithGRPCStreamInterceptors(t *testing.T) {
	// Arrange
	s := &Server{}
	interceptor1 := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	}
	interceptor2 := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	}

	// Act
	o := WithGRPCStreamInterceptors(interceptor1, interceptor2)
	o(s)

	// Assert
	assert.Len(t, s.grpcStreamServerInterceptors, 2)
}

func TestWithGatewayMuxOptions(t *testing.T) {
	// Arrange
	s := &Server{}
	opt1 := runtime.WithMarshalerOption("application/json", &runtime.JSONPb{})
	opt2 := runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	})

	// Act
	o := WithGatewayMuxOptions(opt1, opt2)
	o(s)

	// Assert
	assert.Len(t, s.gwServerMuxOptions, 2)
}

func TestWithGatewayCORS(t *testing.T) {
	// Arrange
	s := &Server{}
	corsOpts := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	}

	// Act
	o := WithGatewayCORS(corsOpts)
	o(s)

	// Assert
	assert.True(t, s.gwCORSEnabled)
	assert.Equal(t, corsOpts, s.gwCORSOptions)
}

func TestConfigShortcuts(t *testing.T) {
	tests := []struct {
		name     string
		option   Option
		validate func(*testing.T, *Server)
	}{
		{
			name:   "WithGRPCAddress",
			option: WithGRPCAddress(":50051"),
			validate: func(t *testing.T, s *Server) {
				assert.Equal(t, ":50051", s.cfg.GRPCAddress)
			},
		},
		{
			name:   "WithHTTPAddress",
			option: WithHTTPAddress(":8081"),
			validate: func(t *testing.T, s *Server) {
				assert.Equal(t, ":8081", s.cfg.HTTPAddress)
			},
		},
		{
			name:   "WithMetricsAddress",
			option: WithMetricsAddress(":9092"),
			validate: func(t *testing.T, s *Server) {
				assert.Equal(t, ":9092", s.cfg.MetricsAddress)
			},
		},
		{
			name:   "WithPprofAddress",
			option: WithPprofAddress(":6061"),
			validate: func(t *testing.T, s *Server) {
				assert.Equal(t, ":6061", s.cfg.PprofAddress)
			},
		},
		{
			name:   "WithCloseTimeout",
			option: WithCloseTimeout(5 * time.Second),
			validate: func(t *testing.T, s *Server) {
				assert.Equal(t, 5*time.Second, s.cfg.CloseTimeout)
			},
		},
		{
			name:   "WithReflection",
			option: WithReflection(false),
			validate: func(t *testing.T, s *Server) {
				assert.False(t, s.cfg.ReflectionEnabled)
			},
		},
		{
			name:   "WithHealthCheck",
			option: WithHealthCheck(false),
			validate: func(t *testing.T, s *Server) {
				assert.False(t, s.cfg.HealthCheckEnabled)
			},
		},
		{
			name:   "WithSwaggerDir",
			option: WithSwaggerDir("/custom/swagger"),
			validate: func(t *testing.T, s *Server) {
				assert.True(t, s.cfg.SwaggerEnabled)
				assert.Equal(t, "/custom/swagger", s.cfg.SwaggerDir)
			},
		},
		{
			name:   "WithSwaggerBasePath",
			option: WithSwaggerBasePath("/api/v1"),
			validate: func(t *testing.T, s *Server) {
				assert.True(t, s.cfg.SwaggerEnabled)
				assert.Equal(t, "/api/v1", s.cfg.SwaggerBasePath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			s := &Server{cfg: config.NewConfig()}

			// Act
			tt.option(s)

			// Assert
			tt.validate(t, s)
		})
	}
}
