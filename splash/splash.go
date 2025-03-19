package splash

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// SplashOption is a function that configures a Splash
type SplashOption func(*Splash)

// Splash represents a splash screen for the application
type Splash struct {
	hostname        string
	goVersion       string
	grpcAddress     string
	httpAddress     string
	metricsAddress  string
	pprofAddress    string
	swaggerEnabled  bool
	swaggerBasePath string
	features        []string
}

// NewSplash creates a new Splash with the given options
func NewSplash(opts ...SplashOption) *Splash {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Default values
	s := &Splash{
		hostname:  hostname,
		goVersion: runtime.Version(),
		features:  []string{},
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithGRPCAddress sets the gRPC address for the splash screen
func WithGRPCAddress(address string) SplashOption {
	return func(s *Splash) {
		s.grpcAddress = address
	}
}

// WithHTTPAddress sets the HTTP address for the splash screen
func WithHTTPAddress(address string) SplashOption {
	return func(s *Splash) {
		s.httpAddress = address
	}
}

// WithMetricsAddress sets the metrics address for the splash screen
func WithMetricsAddress(address string) SplashOption {
	return func(s *Splash) {
		s.metricsAddress = address
	}
}

// WithPprofAddress sets the pprof address for the splash screen
func WithPprofAddress(address string) SplashOption {
	return func(s *Splash) {
		s.pprofAddress = address
	}
}

// WithSwaggerBasePath sets the swagger base path for the splash screen
func WithSwaggerBasePath(path string) SplashOption {
	return func(s *Splash) {
		s.swaggerEnabled = true
		s.swaggerBasePath = path
	}
}

// WithFeature adds a feature to the splash screen
func WithFeature(feature string) SplashOption {
	return func(s *Splash) {
		s.features = append(s.features, feature)
	}
}

// String returns the splash screen as a string
//
//nolint:gocyclo // This function is complex by nature
func (s *Splash) String() string {
	// Create a clean, frameless splash screen
	splash := []string{
		"",
		fmt.Sprintf("💻 Hostname: %s", s.hostname),
		fmt.Sprintf("🔄 Go Version: %s", s.goVersion),
		"",
	}

	// Add endpoints section if any endpoint is set
	if s.grpcAddress != "" || s.httpAddress != "" || s.metricsAddress != "" || s.pprofAddress != "" {
		splash = append(splash, "📡 Endpoints:")

		if s.grpcAddress != "" {
			splash = append(splash, fmt.Sprintf("   • gRPC API: %s", s.grpcAddress))
		}

		if s.httpAddress != "" {
			splash = append(splash, fmt.Sprintf("   • HTTP API: %s", s.httpAddress))
		}

		if s.metricsAddress != "" {
			splash = append(splash, fmt.Sprintf("   • Metrics: %s", s.metricsAddress))
		}

		if s.pprofAddress != "" {
			splash = append(splash, fmt.Sprintf("   • Profiling: %s", s.pprofAddress))
		}

		// Add Swagger information if enabled
		if s.swaggerEnabled {
			// Extract port from HTTP address
			port := strings.TrimPrefix(s.httpAddress, ":")

			// Create clickable link for terminal
			swaggerURL := fmt.Sprintf("http://localhost:%s/swagger", port)
			clickableLink := fmt.Sprintf("\u001B]8;;%s\u0007%s\u001B]8;;\u0007", swaggerURL, swaggerURL)

			splash = append(splash, fmt.Sprintf("   • Swagger UI: %s", clickableLink))
		}

		splash = append(splash, "")
	}

	// Add features information if any
	if len(s.features) > 0 {
		splash = append(splash, "✨ Features:")
		for _, feature := range s.features {
			splash = append(splash, fmt.Sprintf("   • %s", feature))
		}
		splash = append(splash, "")
	}

	return strings.Join(splash, "\n")
}

// Display prints the splash screen to stdout
func (s *Splash) Display() {
	fmt.Print(s.String())
}
