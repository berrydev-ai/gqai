package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/fotoetienne/gqai/graphql"
)

// StreamableHTTPServer represents a streamable HTTP server instance
type StreamableHTTPServer struct {
	config     *graphql.GraphQLConfig
	sessions   map[string]*StreamableHTTPSession
	sessionsMux sync.RWMutex
}

// StreamableHTTPSession represents a session for streamable HTTP
type StreamableHTTPSession struct {
	sessionID string
	responses chan JSONRPCResponse
	done      chan struct{}
}

// NewStreamableHTTPServer creates a new streamable HTTP server
func NewStreamableHTTPServer(config *graphql.GraphQLConfig) *StreamableHTTPServer {
	return &StreamableHTTPServer{
		config:   config,
		sessions: make(map[string]*StreamableHTTPSession),
	}
}

// HandleStreamableHTTP handles the streamable HTTP endpoint
func (s *StreamableHTTPServer) HandleStreamableHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleStreamableHTTPGet(w, r)
	case http.MethodPost:
		s.handleStreamableHTTPPost(w, r)
	case http.MethodDelete:
		s.handleStreamableHTTPDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleStreamableHTTPGet handles GET requests for establishing streamable HTTP connection
func (s *StreamableHTTPServer) handleStreamableHTTPGet(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers for streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Generate session ID
	sessionID := fmt.Sprintf("http_%d", len(s.sessions)+1)

	// Create session
	session := &StreamableHTTPSession{
		sessionID: sessionID,
		responses: make(chan JSONRPCResponse, 10),
		done:      make(chan struct{}),
	}

	// Register session
	s.sessionsMux.Lock()
	s.sessions[sessionID] = session
	s.sessionsMux.Unlock()

	// Send session ID in initial event
	fmt.Fprintf(w, "event: session\ndata: %s\n\n", sessionID)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Clean up on disconnect
	defer func() {
		s.sessionsMux.Lock()
		delete(s.sessions, sessionID)
		s.sessionsMux.Unlock()
		close(session.done)
		close(session.responses)
	}()

	// Stream responses
	for {
		select {
		case response := <-session.responses:
			if response != (JSONRPCResponse{}) {
				responseData, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshaling response: %v", err)
					continue
				}
				fmt.Fprintf(w, "data: %s\n\n", responseData)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		case <-session.done:
			return
		case <-r.Context().Done():
			return
		}
	}
}

// handleStreamableHTTPPost handles POST requests for sending messages
func (s *StreamableHTTPServer) handleStreamableHTTPPost(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("Mcp-Session-Id")
	if sessionID == "" {
		http.Error(w, "Missing Mcp-Session-Id header", http.StatusBadRequest)
		return
	}

	s.sessionsMux.RLock()
	session, exists := s.sessions[sessionID]
	s.sessionsMux.RUnlock()

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

	// Send response to session channel
	select {
	case session.responses <- response:
	default:
		log.Printf("Response channel full for session %s", sessionID)
	}

	w.WriteHeader(http.StatusAccepted)
}

// handleStreamableHTTPDelete handles DELETE requests for ending sessions
func (s *StreamableHTTPServer) handleStreamableHTTPDelete(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("Mcp-Session-Id")
	if sessionID == "" {
		http.Error(w, "Missing Mcp-Session-Id header", http.StatusBadRequest)
		return
	}

	s.sessionsMux.Lock()
	if session, exists := s.sessions[sessionID]; exists {
		close(session.done)
		delete(s.sessions, sessionID)
	}
	s.sessionsMux.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// RunMCPStreamableHTTP starts the MCP server with streamable HTTP transport
func (s *StreamableHTTPServer) RunMCPStreamableHTTP(addr string) {
	http.HandleFunc("/mcp", s.HandleStreamableHTTP)

	log.Printf("Starting MCP Streamable HTTP server on %s", addr)
	log.Printf("Streamable HTTP endpoint: http://%s/mcp", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
