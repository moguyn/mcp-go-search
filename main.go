package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/server"

	"com.moguyn/mcp-go-search/config"
	"com.moguyn/mcp-go-search/mcp"
	"com.moguyn/mcp-go-search/search"
)

// Logger provides a simple structured logging interface
type Logger struct {
	prefix string
}

// NewLogger creates a new logger with the given prefix
func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

// Info logs an informational message with structured data
func (l *Logger) Info(msg string, data map[string]interface{}) {
	l.log("INFO", msg, data)
}

// Error logs an error message with structured data
func (l *Logger) Error(msg string, err error, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	if err != nil {
		data["error"] = err.Error()
	}
	l.log("ERROR", msg, data)
}

// log formats and prints a log message
func (l *Logger) log(level, msg string, data map[string]interface{}) {
	timestamp := time.Now().Format(time.RFC3339)

	// Format the data as key=value pairs
	dataStr := ""
	for k, v := range data {
		dataStr += fmt.Sprintf(" %s=%v", k, v)
	}

	log.Printf("%s [%s] %s: %s%s", timestamp, level, l.prefix, msg, dataStr)
}

func main() {
	logger := NewLogger("main")

	// Log startup
	logger.Info("Starting server", map[string]interface{}{
		"time": time.Now().Format(time.RFC3339),
	})

	// Load configuration
	cfg := config.New()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Error("Configuration error", err, map[string]interface{}{
			"suggestion": "Please set the BOCHA_API_KEY environment variable.",
		})
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
	logger.Info("Server ready", map[string]interface{}{
		"name":    cfg.ServerName,
		"version": cfg.ServerVersion,
	})

	if err := server.ServeStdio(s); err != nil {
		logger.Error("Server error", err, nil)
		os.Exit(1)
	}
}
