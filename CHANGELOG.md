# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2025-10-23 -->

### Added
- HTTP transport support via `-t http` flag for web-based MCP integrations
- SSE (Server-Sent Events) transport support via `-t sse` flag for event-driven clients
- Transport selection flag `-t` with support for stdio (default), http, and sse transports
- Address flag `-a` for configuring listening address for http/sse transports (default: :8080)
- `--help` and `-h` flags to display comprehensive usage information
- `--version` flag to display version information
- Makefile targets for running different transports: `run-http`, `run-http-custom`, `run-sse`, `run-sse-custom`

## Changed

## [1.0.0] - 2025-04-14

### Added
- Initial MCP implementation by @fotoetienne (github.com/fotoetienne/gqai)
