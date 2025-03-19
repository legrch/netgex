package netgex

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/legrch/netgex/internal/gateway"
	"github.com/legrch/netgex/internal/metrics"
	"github.com/legrch/netgex/internal/pprof"
	"github.com/legrch/netgex/pkg/config"
	"github.com/legrch/netgex/pkg/service"
	"github.com/legrch/netgex/pkg/splash"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	grpcserver "github.com/legrch/netgex/internal/grpc"
)

// Constants
const (
	// StartupDelay is the time to wait for processes to start before displaying the splash screen
	StartupDelay = 100 * time.Millisecond
)

// Process is an interface for components that can be started and stopped
type Process interface {
	PreRun(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// Server represents the main entry point for the application
type Server struct {
	cfg                          *config.Config
	processes                    []Process
	logger                       *slog.Logger
	services                     []service.Registrar
	grpcServerOptions            []grpc.ServerOption
	grpcUnaryServerInterceptors  []grpc.UnaryServerInterceptor
	grpcStreamServerInterceptors []grpc.StreamServerInterceptor
	gwServerMuxOptions           []runtime.ServeMuxOption
	gwCORSEnabled                bool
	gwCORSOptions                cors.Options
}

// NewServer creates a new Server with the given options
func NewServer(opts ...Option) *Server {
	s := &Server{
		cfg: config.NewConfig(),
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Run starts the Server and all its processes
func (s *Server) Run(ctx context.Context) error {

	if s.logger == nil {
		s.logger = slog.Default()
	}

	s.logger.Info("starting application")

	// Create gRPC server
	grpcServer := grpcserver.NewServer(
		s.logger,
		s.cfg.CloseTimeout,
		s.cfg.GRPCAddress,
		grpcserver.WithServices(s.services...),
		grpcserver.WithUnaryInterceptors(s.grpcUnaryServerInterceptors...),
		grpcserver.WithStreamInterceptors(s.grpcStreamServerInterceptors...),
		grpcserver.WithReflection(s.cfg.ReflectionEnabled),
		grpcserver.WithHealthCheck(s.cfg.HealthCheckEnabled),
		grpcserver.WithOptions(s.grpcServerOptions...),
	)

	// Create gateway server
	gatewayOpts := []gateway.Option{
		gateway.WithServices(s.services...),
	}

	// Add mux options if provided
	if len(s.gwServerMuxOptions) > 0 {
		gatewayOpts = append(gatewayOpts, gateway.WithMuxOptions(s.gwServerMuxOptions...))
	}

	// Add CORS if enabled
	if s.gwCORSEnabled {
		gatewayOpts = append(gatewayOpts, gateway.WithCORS(&s.gwCORSOptions))
	}

	// Add swagger if configured
	if s.cfg.SwaggerEnabled {
		gatewayOpts = append(gatewayOpts, gateway.WithSwagger(s.cfg.SwaggerDir, s.cfg.SwaggerBasePath))
	}

	gatewayServer := gateway.NewServer(
		s.logger,
		s.cfg.CloseTimeout,
		s.cfg.GRPCAddress,
		s.cfg.HTTPAddress,
		gatewayOpts...,
	)

	metricsServer := metrics.NewServer(s.logger, s.cfg.MetricsAddress, s.cfg.CloseTimeout)
	pprofServer := pprof.NewServer(s.logger, s.cfg.PprofAddress)

	// Add servers to processes
	processes := []Process{grpcServer, gatewayServer, metricsServer, pprofServer}
	processes = append(processes, s.processes...)

	// Run PreRun for all processes
	for _, p := range processes {
		if err := p.PreRun(ctx); err != nil {
			return fmt.Errorf("pre-run error: %w", err)
		}
	}

	// Create error channel
	errCh := make(chan error, len(processes))

	// Start all processes
	for i, p := range processes {
		process := p
		index := i

		go func() {
			s.logger.Info("starting process", "index", index)
			if err := process.Run(ctx); err != nil {
				errCh <- fmt.Errorf("process %d error: %w", index, err)
			}
		}()
	}

	// Give processes a moment to start
	time.Sleep(StartupDelay)

	// Display splash screen after processes have started
	s.displaySplash()

	// Wait for context cancellation or error
	var err error
	select {
	case <-ctx.Done():
		s.logger.Info("context canceled, shutting down")
	case err = <-errCh:
		s.logger.Error("process error", "error", err)
	}

	// Create shutdown context
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.CloseTimeout)
	defer cancel()

	// Shutdown all processes in reverse order
	for i := len(processes) - 1; i >= 0; i-- {
		p := processes[i]
		if shutdownErr := p.Shutdown(shutdownCtx); shutdownErr != nil {
			s.logger.Error("shutdown error", "error", shutdownErr)
			if err == nil {
				err = shutdownErr
			}
		}
	}

	s.logger.Info("application stopped")
	return err
}

// displaySplash initializes and displays the splash screen
func (s *Server) displaySplash() {
	splashOpts := []splash.SplashOption{
		splash.WithSplashGRPCAddress(s.cfg.GRPCAddress),
		splash.WithSplashHTTPAddress(s.cfg.HTTPAddress),
		splash.WithSplashMetricsAddress(s.cfg.MetricsAddress),
		splash.WithSplashPprofAddress(s.cfg.PprofAddress),
	}

	// Add features
	if s.cfg.ReflectionEnabled {
		splashOpts = append(splashOpts, splash.WithSplashFeature("gRPC Reflection"))
	}
	if s.cfg.HealthCheckEnabled {
		splashOpts = append(splashOpts, splash.WithSplashFeature("Health Checks"))
	}
	if s.gwCORSEnabled {
		splashOpts = append(splashOpts, splash.WithSplashFeature("CORS"))
	}

	// Add swagger if enabled
	if s.cfg.SwaggerDir != "" {
		splashOpts = append(splashOpts, splash.WithSplashSwaggerBasePath(s.cfg.SwaggerBasePath))
	}

	// Create and display splash
	splash := splash.NewSplash(splashOpts...)
	splash.Display()
}
