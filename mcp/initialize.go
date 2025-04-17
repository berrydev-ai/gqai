package mcp

// Handles the MCP [initialization request](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/45466e9cccdb49e68746f6c3a005d67a8aa0a16d/schema/2025-03-26/schema.ts#L164)

const initializeMethod = "initialize"

var protocolVersions = []string{
	"2024-11-05",
	"2025-03-26",
}

func mcpInitialize(request JSONRPCRequest) JSONRPCResponse {
	// Check if the request is valid
	if request.Params == nil {
		return errorResponse(request, InvalidParams, "Invalid parameters")
	}
	if _, ok := request.Params.(map[string]interface{}); !ok {
		return errorResponse(request, InvalidParams, "Invalid parameters format")
	}
	if _, ok := request.Params.(map[string]interface{})["protocolVersion"]; !ok {
		return errorResponse(request, InvalidParams, "Missing protocolVersion")
	}
	var clientProtocolVersion = request.Params.(map[string]interface{})["protocolVersion"].(string)

	// # Version Negotiation
	// https://modelcontextprotocol.io/specification/2025-03-26/basic/lifecycle#version-negotiation
	//
	// In the initialize request, the client MUST send a protocol version it supports.
	// This SHOULD be the latest version supported by the client.
	//
	// If the server supports the requested protocol version, it MUST respond with the same version. Otherwise, the
	// server MUST respond with another protocol version it supports.
	// This SHOULD be the latest version supported by the server.
	var protocolVersion string
	for _, version := range protocolVersions {
		if version == clientProtocolVersion {
			protocolVersion = clientProtocolVersion
			break
		}
	}
	if protocolVersion == "" {
		protocolVersion = protocolVersions[len(protocolVersions)-1]
	}

	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: InitializeResult{
			ProtocolVersion: protocolVersion,
			ServerInfo: ServerInfo{
				Name:    "gqai",
				Version: "0.0.0",
			},
			Capabilities: Capabilities{
				Tools: map[string]interface{}{},
			},
		},
	}
}
