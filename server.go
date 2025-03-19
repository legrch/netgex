package netgex

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/legrch/netgex/gateway"
	"github.com/legrch/netgex/metrics"
	"github.com/legrch/netgex/pprof"
	"github.com/legrch/netgex/service"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	grpcserver "github.com/legrch/netgex/grpc"
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
	logger             *slog.Logger
	closeTimeout       time.Duration
	grpcAddress        string
	httpAddress        string
	metricsAddress     string
	pprofAddress       string
	swaggerDir         string
	swaggerBasePath    string
	reflection         bool
	healthCheck        bool
	registrars         []service.Registrar
	processes          []Process
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	corsOptions        cors.Options
	corsEnabled        bool
	jsonConfig         *gateway.JSONConfig
	appName            string
	appVersion         string
}

// NewServer creates a new Server with the given options
func NewServer(opts ...Option) *Server {
	// Default values
	s := &Server{
		logger:          slog.Default(),
		closeTimeout:    10 * time.Second,
		grpcAddress:     getEnv("GRPC_ADDRESS", ":9090"),
		httpAddress:     getEnv("HTTP_ADDRESS", ":8080"),
		metricsAddress:  getEnv("METRICS_ADDRESS", ":9091"),
		pprofAddress:    getEnv("PPROF_ADDRESS", ":6060"),
		swaggerDir:      getEnv("SWAGGER_DIR", "./api"),
		swaggerBasePath: getEnv("SWAGGER_BASE_PATH", "/"),
		reflection:      getEnvBool("REFLECTION_ENABLED"),
		healthCheck:     true,
		jsonConfig:      gateway.DefaultJSONConfig(),
		appName:         getEnv("PROJECT_NAME", "Service"),
		appVersion:      getEnv("VERSION", "dev"),
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Run starts the Server and all its processes
func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("starting application")

	// Create gRPC server
	grpcServer := grpcserver.NewServer(
		s.logger,
		s.closeTimeout,
		s.grpcAddress,
		grpcserver.WithRegistrars(s.registrars...),
		grpcserver.WithUnaryInterceptors(s.unaryInterceptors...),
		grpcserver.WithStreamInterceptors(s.streamInterceptors...),
		grpcserver.WithReflection(s.reflection),
		grpcserver.WithHealthCheck(s.healthCheck),
	)

	// Create gateway server
	gatewayOpts := []gateway.Option{
		gateway.WithRegistrars(s.registrars...),
		gateway.WithJSONConfig(s.jsonConfig),
	}

	if s.corsEnabled {
		gatewayOpts = append(gatewayOpts, gateway.WithCORS(&s.corsOptions))
	}

	if s.swaggerDir != "" {
		gatewayOpts = append(gatewayOpts, gateway.WithSwagger(s.swaggerDir, s.swaggerBasePath))
	}

	gatewayServer := gateway.NewServer(
		s.logger,
		s.closeTimeout,
		s.grpcAddress,
		s.httpAddress,
		gatewayOpts...,
	)

	metricsServer := metrics.NewServer(s.logger, s.metricsAddress, s.closeTimeout)
	pprofServer := pprof.NewServer(s.logger, s.pprofAddress)

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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.closeTimeout)
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
	splashOpts := []SplashOption{
		WithSplashAppName(s.appName),
		WithSplashAppVersion(s.appVersion),
		WithSplashGRPCAddress(s.grpcAddress),
		WithSplashHTTPAddress(s.httpAddress),
		WithSplashMetricsAddress(s.metricsAddress),
		WithSplashPprofAddress(s.pprofAddress),
	}

	// Add features
	if s.reflection {
		splashOpts = append(splashOpts, WithSplashFeature("gRPC Reflection"))
	}
	if s.healthCheck {
		splashOpts = append(splashOpts, WithSplashFeature("Health Checks"))
	}
	if s.corsEnabled {
		splashOpts = append(splashOpts, WithSplashFeature("CORS"))
	}

	// Add swagger if enabled
	if s.swaggerDir != "" {
		splashOpts = append(splashOpts, WithSplashSwaggerBasePath(s.swaggerBasePath))
	}

	// Create and display splash
	splash := NewSplash(splashOpts...)
	splash.Display()
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvBool(key string) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "1" || value == "yes"
	}
	return true
}
