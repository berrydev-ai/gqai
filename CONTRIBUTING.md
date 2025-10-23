# Contributing to gqai

Thank you for your interest in contributing to gqai! This document provides guidelines and instructions for contributors.

## Table of Contents

- [Development Setup](#development-setup)
- [Running Tests](#running-tests)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Testing Requirements](#testing-requirements)
- [Project Structure](#project-structure)

## Development Setup

### Prerequisites

- Go 1.20 or later
- Git

### Getting Started

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/gqai.git
   cd gqai
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Verify the setup:
   ```bash
   go build ./...
   ```

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Run Tests for Specific Package
```bash
go test ./graphql
go test ./mcp
go test ./tool
go test ./cmd
```

### Run Tests with Verbose Output
```bash
go test -v ./...
```

### Run Integration Tests
```bash
go test -tags=integration ./...
```

## Code Style

### Go Code

- Follow standard Go formatting: `go fmt`
- Use `go vet` to check for common errors
- Use `golint` or `golangci-lint` for additional linting
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in imperative mood (e.g., "Add", "Fix", "Update")
- Keep the first line under 50 characters
- Add more detailed explanation in the body if needed

Example:
```
Add error handling for GraphQL execution failures

- Handle network timeouts gracefully
- Return appropriate error messages to clients
- Add tests for error scenarios
```

## Submitting Changes

### Pull Request Process

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and ensure tests pass:
   ```bash
   go test ./...
   ```

3. Update documentation if needed

4. Commit your changes:
   ```bash
   git commit -m "Add your descriptive commit message"
   ```

5. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

6. Create a Pull Request on GitHub

### Pull Request Requirements

- All tests must pass
- Code must be reviewed by at least one maintainer
- Include a clear description of the changes
- Reference any related issues

## Testing Requirements

### Unit Tests

- All new code must include unit tests
- Tests should cover both happy path and error cases
- Use table-driven tests for multiple test cases
- Mock external dependencies when appropriate

### Integration Tests

- Integration tests are located in `integration_test.go`
- These tests verify end-to-end functionality
- Run with: `go test -tags=integration`

### Test Coverage

- Aim for high test coverage
- Use `go test -cover` to check coverage
- Focus on testing public APIs and error conditions

## Project Structure

```
gqai/
├── cmd/                    # CLI commands
│   ├── root.go            # Main CLI entry point
│   └── serve.go           # HTTP server commands
├── graphql/               # GraphQL configuration and execution
│   ├── config.go          # GraphQL config parsing
│   ├── executor.go        # GraphQL query execution
│   ├── operations.go      # GraphQL operation loading
│   └── model.go           # GraphQL data models
├── mcp/                   # MCP (Model Context Protocol) implementation
│   ├── error.go           # Error handling
│   ├── initialize.go      # MCP initialization
│   ├── json.go            # JSON utilities
│   ├── model.go           # MCP data models
│   ├── router.go          # MCP request routing
│   ├── sse.go             # SSE transport
│   ├── stdio.go           # Stdio transport
│   ├── streamable_http.go # HTTP transport
│   ├── tools_call.go      # Tool calling
│   └── tools_list.go      # Tool listing
├── tool/                  # Tool definitions and schemas
│   ├── mcptool.go         # MCP tool structure
│   ├── schema.go          # Schema generation
│   └── tools.go           # Tool management
├── examples/              # Example configurations
├── main.go                # Application entry point
├── go.mod                 # Go module file
└── README.md              # Project documentation
```

## Development Workflow

1. **Choose an Issue**: Look for open issues labeled "good first issue" or "help wanted"
2. **Create a Branch**: Use descriptive branch names
3. **Write Tests First**: Follow TDD principles when possible
4. **Implement Changes**: Keep changes focused and minimal
5. **Test Thoroughly**: Run all tests and verify functionality
6. **Update Documentation**: Keep README and other docs current
7. **Submit PR**: Follow the PR guidelines above

## Getting Help

- Check existing issues and documentation first
- Create an issue for bugs or feature requests
- Join discussions in existing issues

## License

By contributing to gqai, you agree that your contributions will be licensed under the same license as the project (see LICENSE file).

Thank you for contributing to gqai! 🎉
