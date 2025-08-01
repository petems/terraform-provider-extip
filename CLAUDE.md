# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Terraform provider that retrieves external IP addresses as a data source. The provider makes HTTP requests to IP resolution services (like AWS's checkip service) and returns the external IP address for use in Terraform configurations.

## Architecture

- **Main Entry Point**: `main.go` - Standard Terraform provider entry point using the plugin SDK v2
- **Provider Logic**: `extip/` package contains all provider functionality
  - `provider.go` - Provider schema definition and registration
  - `data_source.go` - Core data source implementation with HTTP client pooling and IP resolution
  - `data_source_test.go` - Comprehensive test suite including unit and acceptance tests

## Key Implementation Details

- **HTTP Client Pooling**: Uses a cached HTTP client pool (`httpClients` map) with mutex protection for thread safety and performance optimization
- **Configurable Timeouts**: Client timeout is configurable via `client_timeout` parameter (milliseconds)
- **IP Validation**: Optional IP validation using Go's `net.ParseIP()`
- **Error Handling**: Comprehensive error handling for HTTP requests, timeouts, and invalid responses

## Development Commands

### Building
```bash
make build           # Build and install the provider
go install           # Direct Go build
```

### Testing
```bash
make test           # Run all tests with coverage
make test-short     # Skip network-dependent tests
make test-verbose   # Verbose test output
make testacc        # Run acceptance tests (requires TF_ACC=1)
```

### Code Quality
```bash
make lint           # Run golangci-lint (requires version 2.3.0+)
make lint-fix       # Auto-fix linting issues
make fmt            # Format Go code
make fmtcheck       # Check code formatting
make vet            # Run go vet
```

### Coverage
```bash
make cover_html     # Generate HTML coverage report
```

## Testing Strategy

- **Unit Tests**: Mock HTTP servers for testing different scenarios
- **Acceptance Tests**: Real Terraform provider tests using the plugin SDK testing framework
- **Parameter Validation Tests**: Test invalid inputs and error handling
- **Network Tests**: Can be skipped with `-short` flag for faster development cycles

## Provider Configuration

The provider supports these data source attributes:
- `resolver`: HTTP/HTTPS URL for IP resolution service (default: AWS checkip)
- `client_timeout`: Request timeout in milliseconds (default: 1000ms)
- `validate_ip`: Boolean flag to validate returned IP address format

## Recent Modernization

The codebase was recently upgraded from Terraform Plugin SDK v1 to v2, with performance optimizations including HTTP client pooling and more efficient string handling.