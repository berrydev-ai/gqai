# Technology Stack

## Core Technologies

### Programming Language
- **Go 1.20+**: Primary implementation language
  - High performance and concurrency support
  - Strong typing and memory safety
  - Excellent HTTP server capabilities
  - Rich standard library

### Key Dependencies

#### CLI Framework
- **Cobra**: Command-line interface framework
  - Structured command definitions
  - Flag parsing and validation
  - Help generation and completion

#### Configuration Management
- **Viper**: Configuration management
  - YAML file parsing
  - Environment variable support
  - Multiple configuration sources

#### GraphQL Processing
- **gqlparser/v2**: GraphQL parsing and AST manipulation
  - Standard GraphQL syntax support
  - Query document parsing
  - Variable extraction and validation

#### HTTP Server
- **Gorilla Mux**: HTTP request routing
  - Flexible routing patterns
  - Middleware support
  - Session management for streamable HTTP

## Development Tools

### Build System
- **Go Modules**: Dependency management
- **Makefile**: Build automation
  - Standard Go build/test commands
  - Cross-platform compatibility

### Testing
- **Go Testing Framework**: Unit and integration tests
  - Table-driven tests
  - Benchmarking support
  - Coverage reporting

### Code Quality
- **go fmt**: Code formatting
- **go vet**: Static analysis
- **golint**: Linting (optional)

## Transport Protocols

### Stdio Transport
- Direct process communication
- JSON-RPC 2.0 over stdin/stdout
- Suitable for MCP client integration

### HTTP Transports
- **Streamable HTTP**: REST-like interface with sessions (default transport)
- **Server-Sent Events (SSE)**: Real-time event streaming
- Session-based request correlation

## Configuration Format

### GraphQL Config (.graphqlrc.yml)
- YAML-based configuration
- Environment variable expansion (`${VAR}`)
- Header configuration with auth support
- Document path patterns

### MCP Configuration
- JSON configuration for MCP clients
- Transport-specific settings
- Host/port configuration

## Development Environment

### Prerequisites
- Go 1.20 or later
- Git for version control
- Make for build automation

### Project Structure
- Modular package organization
- Clear separation of concerns
- Comprehensive test coverage
- Example configurations included

## Deployment Options

### CLI Tool
- Direct execution via `go run` or binary
- Configuration via command-line flags
- Suitable for development and testing

### HTTP Server
- Production-ready HTTP server
- Configurable host/port via `-H`/`-p` flags (default :8080)
- Multiple transport options (`-t` flag: stdio/http/sse)
- Session management for stateful operations
