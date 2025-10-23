# System Architecture

## Overview
gqai is a Go-based MCP server that bridges GraphQL APIs to AI models through the Model Context Protocol. It consists of four main packages: `cmd`, `graphql`, `mcp`, and `tool`.

## Package Structure

### cmd/
**Purpose**: CLI interface and command handling
- `root.go`: Main CLI entry point with cobra commands
- `serve.go`: HTTP server commands for different transports

**Key Components**:
- Command definitions for `run`, `serve`, `tools/call`, `tools/list`, `describe`
- Transport selection (stdio, SSE, HTTP) via `-t` flag
- Host/port configuration via `-H`/`-p` flags (default :8080)
- Configuration loading and initialization

### graphql/
**Purpose**: GraphQL configuration parsing and execution
- `config.go`: Loads and parses `.graphqlrc.yml` configuration files
- `operations.go`: Discovers and loads GraphQL operations from `.graphql` files
- `executor.go`: Executes GraphQL requests against configured endpoints
- `model.go`: Data structures for GraphQL requests/responses

**Key Features**:
- Environment variable expansion in config values
- Support for multiple schema endpoints with headers
- Automatic operation discovery from file system
- HTTP client for GraphQL execution

### mcp/
**Purpose**: Model Context Protocol implementation
- `model.go`: JSON-RPC and MCP data structures
- `router.go`: Request routing and method handling
- `initialize.go`: MCP initialization handshake
- `tools_list.go`: Tool discovery endpoint
- `tools_call.go`: Tool execution endpoint
- `stdio.go`: Stdio transport implementation
- `sse.go`: Server-Sent Events transport
- `streamable_http.go`: HTTP transport with session management

**Key Components**:
- JSON-RPC 2.0 protocol handling
- Multiple transport layers (stdio, SSE, HTTP)
- Session management for HTTP transports
- Error handling and response formatting
- SSE server implementation for real-time communication
- Streamable HTTP server with session-based request correlation

### tool/
**Purpose**: Tool definition and schema generation
- `mcptool.go`: MCP tool data structure
- `tools.go`: Tool creation from GraphQL operations
- `schema.go`: Input schema extraction from GraphQL AST

**Key Features**:
- Automatic tool generation from GraphQL operations
- Input schema inference from query variables
- Tool annotations for hints (read-only, destructive, etc.)

## Data Flow

1. **Configuration Loading**: CLI loads `.graphqlrc.yml` and discovers `.graphql` files
2. **Operation Discovery**: GraphQL parser extracts operations and variables
3. **Tool Generation**: Operations are converted to MCP tools with schemas
4. **Transport Setup**: Server starts on selected transport (stdio/SSE/HTTP)
5. **Request Handling**: MCP requests are routed to appropriate handlers
6. **Execution**: Tools execute GraphQL queries against configured endpoints
7. **Response**: Results are formatted as MCP responses

## Transport Options

### Stdio Transport
- Used for direct process communication
- Suitable for MCP clients that spawn the server as a subprocess
- Default transport for `run` command

### SSE Transport
- Server-Sent Events for real-time communication
- Suitable for web-based MCP clients
- Uses `serve --transport sse` command

### HTTP Transport
- Streamable HTTP with session management
- REST-like interface with session headers
- Uses `serve --transport http` command (default)

## Configuration Structure

```yaml
schema: https://api.example.com/graphql
documents: operations/
headers:
  Authorization: Bearer ${TOKEN}
```

## Key Design Decisions

1. **Configuration-Driven**: All behavior defined by `.graphqlrc.yml` files
2. **Transport Agnostic**: Core logic separated from transport implementations
3. **GraphQL Native**: Uses standard GraphQL tooling and patterns
4. **Environment Aware**: Supports environment variable expansion
5. **Session-Based HTTP**: HTTP transport uses sessions for request correlation
