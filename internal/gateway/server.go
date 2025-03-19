package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/legrch/netgex/pkg/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

// HeaderMatcherFunc is a function for matching headers in gRPC gateway
type HeaderMatcherFunc = func(string) (string, bool)

// Option is a function that configures a Server
type Option func(*Server)

// Server represents a gRPC-Gateway server
type Server struct {
	logger                *slog.Logger
	server                *http.Server
	closeTimeout          time.Duration
	grpcAddress           string
	httpAddress           string
	registrars            []service.Registrar
	muxOptions            []runtime.ServeMuxOption
	incomingHeaderMatcher HeaderMatcherFunc
	outgoingHeaderMatcher HeaderMatcherFunc
	corsEnabled           bool
	corsOptions           cors.Options
	pprofEnabled          bool
	swaggerDir            string
	swaggerBasePath       string
	jsonConfig            *JSONConfig
}

// NewServer creates a new gRPC-Gateway server
func NewServer(
	logger *slog.Logger,
	closeTimeout time.Duration,
	grpcAddress string,
	httpAddress string,
	opts ...Option,
) *Server {
	s := &Server{
		logger:       logger,
		closeTimeout: closeTimeout,
		grpcAddress:  grpcAddress,
		httpAddress:  httpAddress,
		server: &http.Server{
			Addr:              httpAddress,
			ReadHeaderTimeout: 5 * time.Second, // Prevent Slowloris attacks
		},
		jsonConfig: DefaultJSONConfig(),
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithServices sets the service registrars for the gateway
func WithServices(registrars ...service.Registrar) Option {
	return func(s *Server) {
		s.registrars = append(s.registrars, registrars...)
	}
}

// WithMuxOptions sets the gRPC-Gateway mux options
func WithMuxOptions(options ...runtime.ServeMuxOption) Option {
	return func(s *Server) {
		s.muxOptions = append(s.muxOptions, options...)
	}
}

// WithIncomingHeaderMatcher sets the incoming header matcher function
func WithIncomingHeaderMatcher(matcher HeaderMatcherFunc) Option {
	return func(s *Server) {
		s.incomingHeaderMatcher = matcher
	}
}

// WithOutgoingHeaderMatcher sets the outgoing header matcher function
func WithOutgoingHeaderMatcher(matcher HeaderMatcherFunc) Option {
	return func(s *Server) {
		s.outgoingHeaderMatcher = matcher
	}
}

// WithCORS enables CORS support
func WithCORS(options *cors.Options) Option {
	return func(s *Server) {
		s.corsEnabled = true
		s.corsOptions = *options
	}
}

// WithPprof enables the pprof profiler
func WithPprof(enabled bool) Option {
	return func(s *Server) {
		s.pprofEnabled = enabled
	}
}

// WithSwagger enables Swagger UI
func WithSwagger(dir, basePath string) Option {
	return func(s *Server) {
		s.swaggerDir = dir
		s.swaggerBasePath = basePath
	}
}

// WithJSONConfig sets the JSON configuration for the gateway
func WithJSONConfig(config *JSONConfig) Option {
	return func(s *Server) {
		s.jsonConfig = config
	}
}

// PreRun prepares the gateway server
func (*Server) PreRun(_ context.Context) error {
	return nil
}

// Run starts the gRPC-Gateway server
func (s *Server) Run(ctx context.Context) error {
	// Create JSON marshaling options
	jsonOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   s.jsonConfig.UseProtoNames,
			EmitUnpopulated: s.jsonConfig.EmitUnpopulated,
			UseEnumNumbers:  s.jsonConfig.UseEnumNumbers,
			AllowPartial:    s.jsonConfig.AllowPartial,
			Multiline:       s.jsonConfig.Multiline,
			Indent:          s.jsonConfig.Indent,
		},
	})

	// Add JSON options to mux options
	muxOptions := append([]runtime.ServeMuxOption{jsonOpts}, s.muxOptions...)

	// Create gRPC-Gateway mux
	gwmux := runtime.NewServeMux(muxOptions...)

	// Set up gRPC connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Register all service handlers
	for _, registrar := range s.registrars {
		if err := registrar.RegisterHTTP(ctx, gwmux, s.grpcAddress, opts); err != nil {
			return fmt.Errorf("failed to register gateway: %w", err)
		}
	}

	// Create root HTTP mux
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	// Add health check endpoints
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Add Swagger UI if configured
	if s.swaggerDir != "" {
		s.registerSwaggerHandler(mux)
	}

	// Apply CORS if enabled
	var handler http.Handler = mux
	if s.corsEnabled {
		handler = cors.New(s.corsOptions).Handler(mux)
	}

	// Set the handler
	s.server.Handler = handler

	// Start the HTTP server
	s.logger.Info("starting gRPC-Gateway server", "address", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("gateway server error: %w", err)
	}

	return nil
}

// Shutdown gracefully stops the gRPC-Gateway server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down gRPC-Gateway server")

	shutdownCtx, cancel := context.WithTimeout(ctx, s.closeTimeout)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("gateway server shutdown error: %w", err)
	}

	return nil
}

// registerSwaggerHandler registers the Swagger UI handler
func (s *Server) registerSwaggerHandler(mux *http.ServeMux) {
	// Check if swagger directory exists
	if _, err := os.Stat(s.swaggerDir); os.IsNotExist(err) {
		s.logger.Warn("swagger directory does not exist", "dir", s.swaggerDir)
		return
	}

	// Find first swagger file
	entries, err := os.ReadDir(s.swaggerDir)
	if err != nil {
		s.logger.Warn("failed to read swagger directory", "error", err)
		return
	}

	// Find first swagger file and serve it as doc.json
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".swagger.json") {
			swaggerFile := filepath.Join(s.swaggerDir, entry.Name())
			mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, swaggerFile)
			})
			break
		}
	}

	// Configure swagger options
	swaggerOptions := []func(config *httpSwagger.Config){
		httpSwagger.URL("doc.json"),
	}

	// Add base path configuration if provided
	if s.swaggerBasePath != "" {
		swaggerOptions = append(swaggerOptions,
			httpSwagger.BeforeScript(`const UrlMutatorPlugin = (system) => ({
			  rootInjects: {
				setScheme: (scheme) => {
				  const jsonSpec = system.getState().toJSON().spec.json;
				  const schemes = Array.isArray(scheme) ? scheme : [scheme];
				  const newJsonSpec = Object.assign({}, jsonSpec, { schemes });
			
				  return system.specActions.updateJsonSpec(newJsonSpec);
				},
				setHost: (host) => {
				  const jsonSpec = system.getState().toJSON().spec.json;
				  const newJsonSpec = Object.assign({}, jsonSpec, { host });
			
				  return system.specActions.updateJsonSpec(newJsonSpec);
				},
				setBasePath: (basePath) => {
				  const jsonSpec = system.getState().toJSON().spec.json;
				  const newJsonSpec = Object.assign({}, jsonSpec, { basePath });
			
				  return system.specActions.updateJsonSpec(newJsonSpec);
				}
			  }
			});`),
			httpSwagger.Plugins([]string{"UrlMutatorPlugin"}),
			httpSwagger.UIConfig(map[string]string{
				"onComplete": fmt.Sprintf(`() => { window.ui.setBasePath('%s') }`, s.swaggerBasePath),
			}),
		)
	}

	// Register swagger handler
	mux.Handle("/swagger/", httpSwagger.Handler(swaggerOptions...))
	s.logger.Info("swagger UI enabled", "basePath", s.swaggerBasePath)
}
