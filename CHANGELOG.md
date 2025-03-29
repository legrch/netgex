# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- None

## [1.1.0] - 2025-03-29

### Added
- Unified OpenTelemetry configuration (tracing, metrics, logging).
- Support for Grafana and Pyroscope integrations.
- HTTP endpoint for explicit trace generation.
- Configurable log level and `pprof` profiling via config.
- Dedicated telemetry server with new initialization logic.

### Changed
- Streamlined telemetry and metrics initialization in server run process.
- Renamed telemetry variables for consistency (`meterProvider` → `meter`, etc.).
- Updated Docker Compose configurations for OTLP and Tempo compatibility.
- Improved documentation: telemetry guide and examples.

### Removed
- Obsolete `test.go` file and outdated JSONConfig explanation from guide.

## [1.0.0] - 2025-03-19

### Added
- Introduced a **dedicated configuration package** for `netgex` server. 
- Added **unit tests** for configuration and gateway components.

### Changed
- Replaced `WithRegistrars` with `WithServices` for naming consistency.
- Updated **import paths** to use `pkg/service` for better structure. 
- Streamlined **server configuration** and improved option functions.
- Updated **gateway import paths** and added a new server implementation. 
- Enhanced **server initialization** and improved logging setup. 
- Removed `server.go` and restructured server initialization logic. 
- Updated `go.mod` dependencies and improved `Taskfile` commands.
- Adjusted **mock generation commands** in `Taskfile`. 
- Updated Swagger configuration to use a boolean flag. 
- Simplified splash screen **configuration** and removed unused app details.
- Enhanced **server tests** with error-handling scenarios
- Improved **Taskfile linting** and tidy-up imports. 

### Removed
- Removed `server.go` and restructured **server initialization logic**.)
- Removed **unused app name and version** from the splash screen. 

## [0.1.0] - 2024-03-19

### Added
- Initial release
- gRPC server implementation with reflection and health checks
- HTTP/REST gateway for gRPC services
- Metrics server for Prometheus
- Profiling server with pprof
- Configuration utilities
- Terminal splash screen
- Examples demonstrating various features