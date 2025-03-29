package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/legrch/netgex/config"
	"github.com/legrch/netgex/internal/telemetry"
	"github.com/legrch/netgex/service"
	"github.com/legrch/netgex/splash"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/legrch/netgex/internal/gateway"
	"github.com/legrch/netgex/internal/metrics"
	"github.com/legrch/netgex/internal/pprof"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	grpcserver "github.com/legrch/netgex/internal/grpc"
)

// Constants
const (
	// StartupDelay is the time to wait for processes to start before displaying the splash screen
	StartupDelay = 100 * time.Millisecond
)

// parseLogLevel converts a string log level to slog.Level
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

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
	telemetryEnabled             bool
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
		// Set LogLevel from config
		slog.SetLogLoggerLevel(parseLogLevel(s.cfg.LogLevel))
	}

	s.logger.Info("starting application")

	// Initialize telemetry if enabled
	var telemetryService *telemetry.Service
	if s.telemetryEnabled {
		telemetryService = telemetry.NewService(s.logger, s.cfg)
		s.addProcesses(telemetryService)
		s.addGRPCUnaryInterceptors(telemetryService.GetUnaryInterceptors()...)
		s.addGRPCStreamInterceptors(telemetryService.GetStreamInterceptors()...)
	}

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
		gateway.WithMuxOptions(s.gwServerMuxOptions...),
		gateway.WithCORS(&s.gwCORSOptions),
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

	// Initialize metrics server
	metricsServer := metrics.NewServer(s.logger, s.cfg.MetricsAddress, s.cfg.CloseTimeout)

	// Initialize pprof server
	pprofServer := pprof.NewServer(s.logger, s.cfg.PprofAddress)

	// Create system processes
	systemProcesses := []Process{grpcServer, gatewayServer, metricsServer, pprofServer}

	s.addProcesses(systemProcesses...)
	// Run PreRun for all processes
	for _, p := range s.processes {
		if err := p.PreRun(ctx); err != nil {
			return fmt.Errorf("pre-run error: %w", err)
		}
	}

	// Create error channel
	errCh := make(chan error, len(s.processes))

	// Start all processes
	for i, p := range s.processes {
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
	for i := len(s.processes) - 1; i >= 0; i-- {
		p := s.processes[i]
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

func (s *Server) addProcesses(processes ...Process) {
	s.processes = append(s.processes, processes...)
}

func (s *Server) addGRPCUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	s.grpcUnaryServerInterceptors = append(s.grpcUnaryServerInterceptors, interceptors...)
}

func (s *Server) addGRPCStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	s.grpcStreamServerInterceptors = append(s.grpcStreamServerInterceptors, interceptors...)
}

// displaySplash initializes and displays the splash screen
func (s *Server) displaySplash() {
	splashOpts := []splash.SplashOption{
		splash.WithGRPCAddress(s.cfg.GRPCAddress),
		splash.WithHTTPAddress(s.cfg.HTTPAddress),
		splash.WithMetricsAddress(s.cfg.MetricsAddress),
		splash.WithPprofAddress(s.cfg.PprofAddress),
	}

	// Add features
	if s.cfg.ReflectionEnabled {
		splashOpts = append(splashOpts, splash.WithFeature("gRPC Reflection"))
	}
	if s.cfg.HealthCheckEnabled {
		splashOpts = append(splashOpts, splash.WithFeature("Health Checks"))
	}
	if s.gwCORSEnabled {
		splashOpts = append(splashOpts, splash.WithFeature("CORS"))
	}

	// Add swagger if enabled
	if s.cfg.SwaggerEnabled {
		splashOpts = append(splashOpts, splash.WithSwaggerBasePath(s.cfg.SwaggerBasePath))
	}

	// Add telemetry features if enabled
	if s.telemetryEnabled {
		if s.cfg.Telemetry.OTEL.Enabled {
			splashOpts = append(splashOpts, splash.WithFeature(
				"OpenTelemetry",
			))
			if s.cfg.Telemetry.OTEL.TracesEnabled {
				splashOpts = append(splashOpts, splash.WithFeature(
					"  ↳ Traces",
				))
			}
			if s.cfg.Telemetry.OTEL.MetricsEnabled {
				splashOpts = append(splashOpts, splash.WithFeature(
					"  ↳ Metrics",
				))
			}
			if s.cfg.Telemetry.OTEL.LogsEnabled {
				splashOpts = append(splashOpts, splash.WithFeature(
					"  ↳ Logs",
				))
			}
		} else {
			// Only show individual backend information if OTEL is not enabled
			if s.cfg.Telemetry.Tracing.Enabled {
				splashOpts = append(splashOpts, splash.WithFeature(
					fmt.Sprintf("Tracing (%s)", s.cfg.Telemetry.Tracing.Backend),
				))
			}

			if s.cfg.Telemetry.Metrics.Enabled {
				splashOpts = append(splashOpts, splash.WithFeature(
					fmt.Sprintf("Metrics (%s)", s.cfg.Telemetry.Metrics.Backend),
				))
			}
		}

		if s.cfg.Telemetry.Profiling.Enabled {
			splashOpts = append(splashOpts, splash.WithFeature(
				fmt.Sprintf("Profiling (%s)", s.cfg.Telemetry.Profiling.Backend),
			))
		}
	}

	// Create and display splash
	splash := splash.NewSplash(splashOpts...)
	splash.Display()
}
