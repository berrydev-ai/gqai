package mcp

func errorResponse(request JSONRPCRequest, code int, message string) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
		},
	}
}
