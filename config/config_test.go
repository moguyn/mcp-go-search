package config

import (
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
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

	// Test with default values
	os.Unsetenv("BOCHA_API_KEY")
	os.Unsetenv("BOCHA_API_BASE_URL")
	os.Unsetenv("HTTP_TIMEOUT")
	os.Unsetenv("SERVER_NAME")
	os.Unsetenv("SERVER_VERSION")

	cfg := New()

	if cfg.BochaAPIKey != "" {
		t.Errorf("Expected empty API key, got %s", cfg.BochaAPIKey)
	}
	if cfg.BochaAPIBaseURL != "https://api.bochaai.com/v1/ai-search" {
		t.Errorf("Expected default API base URL, got %s", cfg.BochaAPIBaseURL)
	}
	if cfg.HTTPTimeout != 10*time.Second {
		t.Errorf("Expected default HTTP timeout, got %s", cfg.HTTPTimeout)
	}
	if cfg.ServerName != "Bocha AI Search Server" {
		t.Errorf("Expected default server name, got %s", cfg.ServerName)
	}
	if cfg.ServerVersion != "1.0.0" {
		t.Errorf("Expected default server version, got %s", cfg.ServerVersion)
	}

	// Test with custom values
	os.Setenv("BOCHA_API_KEY", "test-api-key")
	os.Setenv("BOCHA_API_BASE_URL", "https://test.api.com")
	os.Setenv("HTTP_TIMEOUT", "5s")
	os.Setenv("SERVER_NAME", "Test Server")
	os.Setenv("SERVER_VERSION", "2.0.0")

	cfg = New()

	if cfg.BochaAPIKey != "test-api-key" {
		t.Errorf("Expected custom API key, got %s", cfg.BochaAPIKey)
	}
	if cfg.BochaAPIBaseURL != "https://test.api.com" {
		t.Errorf("Expected custom API base URL, got %s", cfg.BochaAPIBaseURL)
	}
	if cfg.HTTPTimeout != 5*time.Second {
		t.Errorf("Expected custom HTTP timeout, got %s", cfg.HTTPTimeout)
	}
	if cfg.ServerName != "Test Server" {
		t.Errorf("Expected custom server name, got %s", cfg.ServerName)
	}
	if cfg.ServerVersion != "2.0.0" {
		t.Errorf("Expected custom server version, got %s", cfg.ServerVersion)
	}
}

func TestValidate(t *testing.T) {
	// Test with valid configuration
	cfg := &Config{
		BochaAPIKey:     "test-api-key",
		BochaAPIBaseURL: "https://test.api.com",
		HTTPTimeout:     10 * time.Second,
		ServerName:      "Test Server",
		ServerVersion:   "1.0.0",
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test with missing API key
	cfg.BochaAPIKey = ""
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for missing API key, got nil")
	}

	// Test with missing API base URL
	cfg.BochaAPIKey = "test-api-key"
	cfg.BochaAPIBaseURL = ""
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for missing API base URL, got nil")
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	// Save original environment variable to restore later
	origValue := os.Getenv("TEST_ENV_VAR")
	defer os.Setenv("TEST_ENV_VAR", origValue)

	// Test with unset environment variable
	os.Unsetenv("TEST_ENV_VAR")
	value := getEnvWithDefault("TEST_ENV_VAR", "default-value")
	if value != "default-value" {
		t.Errorf("Expected default value, got %s", value)
	}

	// Test with set environment variable
	os.Setenv("TEST_ENV_VAR", "custom-value")
	value = getEnvWithDefault("TEST_ENV_VAR", "default-value")
	if value != "custom-value" {
		t.Errorf("Expected custom value, got %s", value)
	}
}

func TestGetEnvDurationWithDefault(t *testing.T) {
	// Save original environment variable to restore later
	origValue := os.Getenv("TEST_DURATION")
	defer os.Setenv("TEST_DURATION", origValue)

	// Test with unset environment variable
	os.Unsetenv("TEST_DURATION")
	duration := getEnvDurationWithDefault("TEST_DURATION", 10*time.Second)
	if duration != 10*time.Second {
		t.Errorf("Expected default duration, got %s", duration)
	}

	// Test with seconds as integer
	os.Setenv("TEST_DURATION", "5")
	duration = getEnvDurationWithDefault("TEST_DURATION", 10*time.Second)
	if duration != 5*time.Second {
		t.Errorf("Expected 5s duration, got %s", duration)
	}

	// Test with duration string
	os.Setenv("TEST_DURATION", "2m30s")
	duration = getEnvDurationWithDefault("TEST_DURATION", 10*time.Second)
	if duration != 2*time.Minute+30*time.Second {
		t.Errorf("Expected 2m30s duration, got %s", duration)
	}

	// Test with invalid duration
	os.Setenv("TEST_DURATION", "invalid")
	duration = getEnvDurationWithDefault("TEST_DURATION", 10*time.Second)
	if duration != 10*time.Second {
		t.Errorf("Expected default duration for invalid input, got %s", duration)
	}
}
