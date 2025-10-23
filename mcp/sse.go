package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/fotoetienne/gqai/graphql"
)

// SSEServer represents an SSE server instance
type SSEServer struct {
	config     *graphql.GraphQLConfig
	clients    map[string]*SSEClient
	clientsMux sync.RWMutex
}

// SSEClient represents a connected SSE client
type SSEClient struct {
	sessionID string
	writer    http.ResponseWriter
	done      chan struct{}
}

// NewSSEServer creates a new SSE server
func NewSSEServer(config *graphql.GraphQLConfig) *SSEServer {
	return &SSEServer{
		config:  config,
		clients: make(map[string]*SSEClient),
	}
}

// HandleSSE handles the SSE endpoint
func (s *SSEServer) HandleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Generate session ID
	sessionID := fmt.Sprintf("sse_%d", len(s.clients)+1)

	// Create client
	client := &SSEClient{
		sessionID: sessionID,
		writer:    w,
		done:      make(chan struct{}),
	}

	// Register client
	s.clientsMux.Lock()
	s.clients[sessionID] = client
	s.clientsMux.Unlock()

	// Send endpoint event
	fmt.Fprintf(w, "event: endpoint\ndata: /message?sessionId=%s\n\n", sessionID)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Clean up on disconnect
	defer func() {
		s.clientsMux.Lock()
		delete(s.clients, sessionID)
		s.clientsMux.Unlock()
		close(client.done)
	}()

	// Keep connection alive
	for {
		select {
		case <-client.done:
			return
		case <-r.Context().Done():
			return
		}
	}
}

// HandleMessage handles incoming messages for SSE clients
func (s *SSEServer) HandleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		http.Error(w, "Missing sessionId parameter", http.StatusBadRequest)
		return
	}

	s.clientsMux.RLock()
	client, exists := s.clients[sessionID]
	s.clientsMux.RUnlock()

	if !exists {
		http.Error(w, "Invalid session", http.StatusNotFound)
		return
	}

	// Parse the incoming JSON-RPC request
	var request JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Route the request
	response := RouteMCPRequest(request, s.config)

	// Send response back via SSE if it's not empty
	if response != (JSONRPCResponse{}) {
		responseData, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			return
		}

		fmt.Fprintf(client.writer, "event: message\ndata: %s\n\n", responseData)
		if f, ok := client.writer.(http.Flusher); ok {
			f.Flush()
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

// RunMCPSSE starts the MCP server with SSE transport
func RunMCPSSE(config *graphql.GraphQLConfig, addr string) {
	server := NewSSEServer(config)

	http.HandleFunc("/sse", server.HandleSSE)
	http.HandleFunc("/message", server.HandleMessage)

	log.Printf("Starting MCP SSE server on %s", addr)
	log.Printf("SSE endpoint: http://%s/sse", addr)
	log.Printf("Message endpoint: http://%s/message", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
