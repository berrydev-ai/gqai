package main

import "github.com/fotoetienne/gqai/cmd"

func main() {
	cmd.Execute()
}

//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"os"
//	"path/filepath"
//
//	"gopkg.in/yaml.v3"
//)
//
//type GraphQLRC struct {
//	Schema     []string `yaml:"schema"`
//	Documents  string   `yaml:"documents"`
//	Extensions struct {
//		Endpoints map[string]struct {
//			URL string `yaml:"url"`
//		} `yaml:"endpoints"`
//	} `yaml:"extensions"`
//}
//
//type GraphQLRequest struct {
//	Query     string                 `json:"query"`
//	Variables map[string]interface{} `json:"variables,omitempty"`
//}
//
//type Tool struct {
//	Name  string
//	Query string
//	Type  string // query or mutation
//}
//
//var tools = map[string]Tool{}
//
//func loadGraphQLRC(path string) (*GraphQLRC, error) {
//	data, err := ioutil.ReadFile(path)
//	if err != nil {
//		return nil, err
//	}
//	var config GraphQLRC
//	err = yaml.Unmarshal(data, &config)
//	if err != nil {
//		return nil, err
//	}
//	return &config, nil
//}
//
//func loadGraphQLOperations(path string) error {
//	files, err := filepath.Glob(path)
//	if err != nil {
//		return err
//	}
//	for _, file := range files {
//		content, err := ioutil.ReadFile(file)
//		if err != nil {
//			return err
//		}
//		name := filepath.Base(file)
//		base := name[:len(name)-len(filepath.Ext(name))]
//		typeStr := "query"
//		if string(content[:8]) == "mutation" {
//			typeStr = "mutation"
//		}
//		tools[base] = Tool{
//			Name:  base,
//			Query: string(content),
//			Type:  typeStr,
//		}
//	}
//	return nil
//}
//
//func callTool(endpoint string, tool Tool, vars map[string]interface{}) ([]byte, error) {
//	reqBody := GraphQLRequest{
//		Query:     tool.Query,
//		Variables: vars,
//	}
//	jsonBody, _ := json.Marshal(reqBody)
//	resp, err := http.Post(endpoint, "application/json", ioutil.NopCloser(
//		bytes.NewReader(jsonBody)))
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	return ioutil.ReadAll(resp.Body)
//}
//
//func main() {
//	if len(os.Args) < 3 {
//		fmt.Println("Usage: mcp-proxy <tool-name> <vars-json>")
//		os.Exit(1)
//	}
//	toolName := os.Args[1]
//	varsJSON := os.Args[2]
//
//	config, err := loadGraphQLRC(".graphqlrc.yml")
//	if err != nil {
//		log.Fatalf("Failed to load .graphqlrc: %v", err)
//	}
//
//	docsPath := config.Documents
//	if err := loadGraphQLOperations(docsPath); err != nil {
//		log.Fatalf("Failed to load GraphQL operations: %v", err)
//	}
//
//	tool, ok := tools[toolName]
//	if !ok {
//		log.Fatalf("Tool %s not found", toolName)
//	}
//
//	var vars map[string]interface{}
//	err = json.Unmarshal([]byte(varsJSON), &vars)
//	if err != nil {
//		log.Fatalf("Failed to parse vars JSON: %v", err)
//	}
//
//	endpoint := config.Extensions.Endpoints["default"].URL
//	resp, err := callTool(endpoint, tool, vars)
//	if err != nil {
//		log.Fatalf("GraphQL call failed: %v", err)
//	}
//
//	fmt.Println(string(resp))
//}
