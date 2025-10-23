package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/berrydev-ai/gqai/graphql"
	"github.com/berrydev-ai/gqai/server"
	"github.com/berrydev-ai/gqai/tool"
	"github.com/spf13/cobra"
)

var config *graphql.GraphQLConfig
var configPath string
var host string
var port int

var rootCmd = &cobra.Command{
	Use:   "gqai",
	Short: "gqai - expose GraphQL operations as AI tools",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run gqai as an MCP server in stdin/stdout mode",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := server.SetupMCPServer(config)
		if err != nil {
			log.Fatalf("Failed to setup MCP server: %v", err)
		}

		if err := server.StartServer(s, "stdio", ""); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	},
}


var toolsCallCmd = &cobra.Command{
	Use:   "tools/call [toolName] [jsonInput]",
	Short: "Call a GraphQL operation as a tool",
	Args:  cobra.MinimumNArgs(1), // allow just the tool name
	Run: func(cmd *cobra.Command, args []string) {
		toolName := args[0]
		var input map[string]any
		if len(args) > 1 {
			if err := json.Unmarshal([]byte(args[1]), &input); err != nil {
				fmt.Println("Invalid JSON input:", err)
				os.Exit(1)
			}
		} else {
			input = map[string]any{} // default to empty input
		}

		// Load tool directly for CLI usage
		tool, err := tool.LoadTool(config, toolName)
		if err != nil {
			fmt.Printf("Error loading tool: %v\n", err)
			os.Exit(1)
		}

		// Execute the tool
		result, err := tool.Execute(input)
		if err != nil {
			fmt.Printf("Error executing tool: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	},
}

var toolsListCmd = &cobra.Command{
	Use:   "tools/list",
	Short: "List available tools",
	Run: func(cmd *cobra.Command, args []string) {
		tools, err := tool.ToolsFromConfig(config)
		if err != nil {
			fmt.Println("Error loading tools:", err)
			os.Exit(1)
		}

		out, err := json.MarshalIndent(tools, "", "  ")
		if err != nil {
			fmt.Println("Error serializing tools:", err)
			os.Exit(1)
		}
		fmt.Println(string(out))
	},
}

var describeCmd = &cobra.Command{
	Use:   "describe [toolName]",
	Short: "Describe a tool and show its full schema",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		toolName := args[0]
		tools, err := tool.ToolsFromConfig(config)
		if err != nil {
			fmt.Println("Error loading tools:", err)
			os.Exit(1)
		}

		for _, t := range tools {
			if t.Name == toolName {
				out, _ := json.MarshalIndent(t, "", "  ")
				fmt.Println(string(out))
				return
			}
		}

		fmt.Printf("Tool %s not found\n", toolName)
	},
}

var transport string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve MCP server over HTTP with configurable transport",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := server.SetupMCPServer(config)
		if err != nil {
			log.Fatalf("Failed to setup MCP server: %v", err)
		}

		addr := fmt.Sprintf("%s:%d", host, port)

		if err := server.StartServer(s, transport, addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", ".graphqlrc.yml", "Path to .graphqlrc.yml")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "localhost", "Host to bind to")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Port to bind to")

	serveCmd.Flags().StringVarP(&transport, "transport", "t", "http", "Transport type: 'sse' or 'http'")

	cobra.OnInitialize(func() {
		var err error
		config, err = graphql.LoadGraphQLConfig(configPath)
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
	})

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(toolsCallCmd)
	rootCmd.AddCommand(toolsListCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.Execute()
}
