package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"com.moguyn/mcp-go-search/config"
	"com.moguyn/mcp-go-search/mcp"
	"com.moguyn/mcp-go-search/search"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Printf("Configuration error: %v", err)
		log.Println("Please set the BOCHA_API_KEY environment variable.")
		log.Println("Example: export BOCHA_API_KEY=\"your-api-key-here\"")
		os.Exit(1)
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		cfg.ServerName,
		cfg.ServerVersion,
		server.WithLogging(),
	)

	// Create the search service
	searchService := search.NewBochaServiceWithConfig(cfg)

	// Create the search tool
	searchTool := mcp.NewSearchTool(searchService)

	// Add the search tool to the server
	s.AddTool(searchTool.Definition(), searchTool.Handler())

	// Start the server
	log.Printf("Starting %s v%s...", cfg.ServerName, cfg.ServerVersion)
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
