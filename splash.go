package netgex

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
	appName         string
	appVersion      string
	environment     string
	hostname        string
	goVersion       string
	grpcAddress     string
	httpAddress     string
	metricsAddress  string
	pprofAddress    string
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
		appName:     getEnv("PROJECT_NAME", "My Service"),
		appVersion:  getEnv("VERSION", "dev"),
		environment: getEnv("ENVIRONMENT", "development"),
		hostname:    hostname,
		goVersion:   runtime.Version(),
		features:    []string{},
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithAppName sets the application name for the splash screen
func WithSplashAppName(name string) SplashOption {
	return func(s *Splash) {
		s.appName = name
	}
}

// WithAppVersion sets the application version for the splash screen
func WithSplashAppVersion(version string) SplashOption {
	return func(s *Splash) {
		s.appVersion = version
	}
}

// WithEnvironment sets the environment for the splash screen
func WithSplashEnvironment(env string) SplashOption {
	return func(s *Splash) {
		s.environment = env
	}
}

// WithGRPCAddress sets the gRPC address for the splash screen
func WithSplashGRPCAddress(address string) SplashOption {
	return func(s *Splash) {
		s.grpcAddress = address
	}
}

// WithHTTPAddress sets the HTTP address for the splash screen
func WithSplashHTTPAddress(address string) SplashOption {
	return func(s *Splash) {
		s.httpAddress = address
	}
}

// WithMetricsAddress sets the metrics address for the splash screen
func WithSplashMetricsAddress(address string) SplashOption {
	return func(s *Splash) {
		s.metricsAddress = address
	}
}

// WithPprofAddress sets the pprof address for the splash screen
func WithSplashPprofAddress(address string) SplashOption {
	return func(s *Splash) {
		s.pprofAddress = address
	}
}

// WithSwaggerBasePath sets the swagger base path for the splash screen
func WithSplashSwaggerBasePath(path string) SplashOption {
	return func(s *Splash) {
		s.swaggerBasePath = path
	}
}

// WithFeature adds a feature to the splash screen
func WithSplashFeature(feature string) SplashOption {
	return func(s *Splash) {
		s.features = append(s.features, feature)
	}
}

// String returns the splash screen as a string
//
//nolint:gocyclo // This function is complex by nature
func (s *Splash) String() string {
	// Format application name - capitalize first letter and replace hyphens with spaces
	formattedAppName := s.appName
	if formattedAppName == "" {
		formattedAppName = "Service"
	} else {
		// Replace hyphens with spaces
		formattedAppName = strings.ReplaceAll(formattedAppName, "-", " ")

		// Capitalize each word
		words := strings.Fields(formattedAppName)
		for i, word := range words {
			if word != "" {
				words[i] = strings.ToUpper(word[:1]) + word[1:]
			}
		}
		formattedAppName = strings.Join(words, " ")
	}

	// Format version - ensure it's in the format v0.0.0
	version := s.appVersion
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	// Create a clean, frameless splash screen
	splash := []string{
		"",
		fmt.Sprintf("🚀 %s %s", formattedAppName, version),
		"",
		fmt.Sprintf("🌐 Environment: %s", strings.ToUpper(s.environment)),
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
		if s.httpAddress != "" {
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
