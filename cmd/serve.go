package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fotoetienne/gqai/tool"
	"github.com/gorilla/mux"
)

func listToolsHandler(w http.ResponseWriter, r *http.Request) {
	tools, err := tool.ToolsFromConfig(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading tools: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tools); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode tools: %v", err), http.StatusInternalServerError)
	}
}

func callToolHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ToolName string         `json:"toolName"`
		Input    map[string]any `json:"input"`
	}

	// Parse input
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Execute the tool
	result, err := tool.Execute(config, payload.ToolName, payload.Input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		return
	}

	// Return result wrapped in MCP response format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"output": result,
	})
}

func serveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	toolName := vars["toolName"]

	var payload struct {
		Input map[string]any `json:"input"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	result, err := tool.Execute(config, toolName, payload.Input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"output": result,
	})
}
