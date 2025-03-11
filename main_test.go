package main

import (
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

// TestConfigValidation tests the configuration validation
func TestConfigValidation(t *testing.T) {
	// Save original environment variables to restore later
	origAPIKey := os.Getenv("BOCHA_API_KEY")
	origAPIBaseURL := os.Getenv("BOCHA_API_BASE_URL")
	origHTTPTimeout := os.Getenv("HTTP_TIMEOUT")
	origServerName := os.Getenv("SERVER_NAME")
	origServerVersion := os.Getenv("SERVER_VERSION")

	// Restore environment variables after the test
	defer func() {
		os.Setenv("BOCHA_API_KEY", origAPIKey)
		os.Setenv("BOCHA_API_BASE_URL", origAPIBaseURL)
		os.Setenv("HTTP_TIMEOUT", origHTTPTimeout)
		os.Setenv("SERVER_NAME", origServerName)
		os.Setenv("SERVER_VERSION", origServerVersion)
	}()

	// Test with missing API key
	os.Unsetenv("BOCHA_API_KEY")
	os.Setenv("BOCHA_API_BASE_URL", "https://test.api.com")
	os.Setenv("HTTP_TIMEOUT", "5s")
	os.Setenv("SERVER_NAME", "Test Server")
	os.Setenv("SERVER_VERSION", "1.0.0")

	// Call runServer - it should return an error
	err := runServer()
	if err == nil {
		t.Error("Expected error when API key is not set, but got nil")
	}
}

// TestConfigSuccess tests the successful configuration case
func TestConfigSuccess(t *testing.T) {
	// Save original environment variables to restore later
	origAPIKey := os.Getenv("BOCHA_API_KEY")
	origAPIBaseURL := os.Getenv("BOCHA_API_BASE_URL")
	origHTTPTimeout := os.Getenv("HTTP_TIMEOUT")
	origServerName := os.Getenv("SERVER_NAME")
	origServerVersion := os.Getenv("SERVER_VERSION")

	// Save original serveStdio function
	origServeStdio := serveStdio
	defer func() { serveStdio = origServeStdio }()

	// Mock the serveStdio function
	serveStdio = func(_ *server.MCPServer) error {
		return nil
	}

	// Restore environment variables after the test
	defer func() {
		os.Setenv("BOCHA_API_KEY", origAPIKey)
		os.Setenv("BOCHA_API_BASE_URL", origAPIBaseURL)
		os.Setenv("HTTP_TIMEOUT", origHTTPTimeout)
		os.Setenv("SERVER_NAME", origServerName)
		os.Setenv("SERVER_VERSION", origServerVersion)
	}()

	// Set valid configuration
	os.Setenv("BOCHA_API_KEY", "test-api-key-for-testing")
	os.Setenv("BOCHA_API_BASE_URL", "https://test.api.com")
	os.Setenv("HTTP_TIMEOUT", "5s")
	os.Setenv("SERVER_NAME", "Test Server")
	os.Setenv("SERVER_VERSION", "1.0.0")

	// Call runServer - it should not return an error
	err := runServer()
	if err != nil {
		t.Errorf("Expected no error with valid configuration, but got: %v", err)
	}
}
