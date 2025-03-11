package main

import (
	"os"
	"testing"
)

// TestMainConfigValidation tests the configuration validation in main
func TestMainConfigValidation(t *testing.T) {
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

	// We can't easily test the main function directly since it calls os.Exit
	// Instead, we'll test the config validation logic separately

	// Set up a mock os.Exit function
	origOsExit := osExit
	defer func() { osExit = origOsExit }()

	var exitCode int
	osExit = func(code int) {
		exitCode = code
		panic("os.Exit called")
	}

	// Recover from the panic caused by our mock os.Exit
	defer func() {
		if r := recover(); r != nil {
			if r != "os.Exit called" {
				t.Errorf("Unexpected panic: %v", r)
			}
		}

		// Check that the exit code is 1 (error)
		if exitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", exitCode)
		}
	}()

	// Call the function that would trigger the validation error
	// This should call our mock os.Exit and be caught by the recover above
	main()
}

// Override os.Exit for testing
var osExit = os.Exit
