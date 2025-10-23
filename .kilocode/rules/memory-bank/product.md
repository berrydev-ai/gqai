# Product Description

## What gqai Does
gqai is a lightweight proxy that transforms GraphQL operations into Model Context Protocol (MCP) tools, enabling AI models like Claude, Cursor, and ChatGPT to directly interact with GraphQL backends through natural language commands.

## Core Value Proposition
- **Zero Code Required**: Define tools using standard GraphQL queries/mutations
- **AI-Native Integration**: Exposes GraphQL operations as MCP tools that AI can discover and call
- **Multiple Transports**: Supports stdio, SSE, and HTTP transports for different integration scenarios
- **Configuration-Driven**: Uses familiar `.graphqlrc.yml` + `.graphql` files for tool definition

## User Experience Goals
- **Simple Setup**: Configure once with GraphQL config, get MCP server automatically
- **Transparent Operation**: AI models interact with GraphQL backend as if calling native tools
- **Flexible Deployment**: Run as CLI tool, HTTP server, or embedded in other systems
- **Developer-Friendly**: Uses familiar GraphQL tooling and patterns

## Problems Solved
1. **AI-GraphQL Integration Gap**: Traditional GraphQL clients require code; gqai enables natural language interaction
2. **Tool Discovery**: AI models can automatically discover available GraphQL operations as tools
3. **Authentication**: Supports header-based auth through config files with environment variable expansion
4. **Multiple Endpoints**: Can work with different GraphQL backends through configuration

## Target Users
- **AI Application Developers**: Want to give AI models access to their GraphQL APIs
- **GraphQL API Owners**: Want to expose their APIs to AI assistants
- **MCP Server Operators**: Need to create MCP servers from existing GraphQL backends
- **DevOps Teams**: Need to deploy GraphQL APIs as AI-accessible services
