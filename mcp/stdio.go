package mcp

import (
	"encoding/json"
	"github.com/fotoetienne/gqai/graphql"
	"log"
	"os"
)

// RunMCPStdIO starts the MCP server in stdin/stdout mode.
func RunMCPStdIO(config *graphql.Config) {
	// Read from stdin and write to stdout
	// This is a blocking call, so it will not return until the program exits

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	// Set up logging to stderr
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Printf("Starting MCP server ...")

	for {
		var request JSONRPCRequest
		if err := decoder.Decode(&request); err != nil {
			log.Printf("Error decoding request: %v", err)
			var response = errorResponse(request, ParseError, "Failed to parse JSON")
			sendResponse(encoder, response)
			break
		}

		log.Printf("Received request: %v", PrettyJSON(request))

		if request.JSONRPC != "2.0" {
			var response = errorResponse(request, InvalidRequest, "Only JSON-RPC 2.0 is supported")
			sendResponse(encoder, response)
			continue
		}

		var response = RouteMCPRequest(request, config)

		// Send response if it's not empty
		if !(response == JSONRPCResponse{}) {
			log.Printf("Sending response: %v", PrettyJSON(response))
			sendResponse(encoder, response)
		}
	}

}

func handleInput(request JSONRPCRequest) {
	// Handle input from the request
	// This is a placeholder function
	log.Printf("Handling input: %v", request)
}

func sendResponse(encoder *json.Encoder, response interface{}) {
	if err := encoder.Encode(response); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}
