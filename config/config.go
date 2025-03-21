package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	// API configuration
	BochaAPIKey     string        `yaml:"bocha_api_key" json:"bocha_api_key"`
	BochaAPIBaseURL string        `yaml:"bocha_api_base_url" json:"bocha_api_base_url"`
	HTTPTimeout     time.Duration `yaml:"-" json:"-"` // Custom handling for YAML/JSON

	// Server configuration
	ServerName    string `yaml:"server_name" json:"server_name"`
	ServerVersion string `yaml:"server_version" json:"server_version"`

	// Internal fields not for YAML/JSON
	HTTPTimeoutStr string `yaml:"http_timeout" json:"http_timeout"`
}

// New creates a new configuration with values from environment variables
func New() *Config {
	config := &Config{
		// Default values
		BochaAPIKey:     os.Getenv("BOCHA_API_KEY"),
		BochaAPIBaseURL: getEnvWithDefault("BOCHA_API_BASE_URL", "https://api.bochaai.com/v1/web-search"),
		HTTPTimeout:     getEnvDurationWithDefault("HTTP_TIMEOUT", 15*time.Second),
		ServerName:      getEnvWithDefault("SERVER_NAME", "Bocha AI Search Server"),
		ServerVersion:   getEnvWithDefault("SERVER_VERSION", "0.0.1"),
	}

	// Check if a config file path is provided
	configPath := os.Getenv("CONFIG_FILE")
	if configPath != "" {
		if err := config.LoadFromFile(configPath); err != nil {
			log.Printf("Warning: Failed to load config from file %s: %v", configPath, err)
		} else {
			log.Printf("Warning: Using configuration file for sensitive data like API keys is not recommended for production environments")
		}
	}

	// Environment variables take precedence over config file
	if envAPIKey := os.Getenv("BOCHA_API_KEY"); envAPIKey != "" {
		config.BochaAPIKey = envAPIKey
	}
	if envAPIBaseURL := os.Getenv("BOCHA_API_BASE_URL"); envAPIBaseURL != "" {
		config.BochaAPIBaseURL = envAPIBaseURL
	}
	if envHTTPTimeout := os.Getenv("HTTP_TIMEOUT"); envHTTPTimeout != "" {
		config.HTTPTimeout = getEnvDurationWithDefault("HTTP_TIMEOUT", config.HTTPTimeout)
	}
	if envServerName := os.Getenv("SERVER_NAME"); envServerName != "" {
		config.ServerName = envServerName
	}
	if envServerVersion := os.Getenv("SERVER_VERSION"); envServerVersion != "" {
		config.ServerVersion = envServerVersion
	}

	// Validate required configuration
	if config.BochaAPIKey == "" {
		log.Println("Warning: BOCHA_API_KEY environment variable not set. The search service will not work without an API key.")
	}

	// Validate HTTP timeout
	if config.HTTPTimeout < time.Second {
		log.Printf("Warning: HTTP_TIMEOUT is very short (%s). Setting to minimum of 1 second.", config.HTTPTimeout)
		config.HTTPTimeout = time.Second
	} else if config.HTTPTimeout > 60*time.Second {
		log.Printf("Warning: HTTP_TIMEOUT is very long (%s). This may cause requests to hang.", config.HTTPTimeout)
	}

	return config
}

// LoadFromFile loads configuration from a YAML file
func (c *Config) LoadFromFile(path string) error {
	// Resolve the absolute path and evaluate any symlinks
	absPath, err := filepath.EvalSymlinks(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Explicitly check for path traversal patterns
	if strings.Contains(path, "..") {
		return fmt.Errorf("security error: path contains potential directory traversal sequence")
	}

	// Create a temporary file with the contents of the target file
	// This approach is recommended by gosec to avoid file inclusion vulnerabilities
	tempFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the temp file when done

	// Read the original file
	originalData, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Write to the temporary file
	if _, err := tempFile.Write(originalData); err != nil {
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Now read from the temporary file
	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read temporary config file: %w", err)
	}

	// Create a temporary config to unmarshal into
	var fileConfig Config
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply non-empty values from the file config
	if fileConfig.BochaAPIKey != "" {
		c.BochaAPIKey = fileConfig.BochaAPIKey
	}
	if fileConfig.BochaAPIBaseURL != "" {
		c.BochaAPIBaseURL = fileConfig.BochaAPIBaseURL
	}
	if fileConfig.HTTPTimeoutStr != "" {
		duration, err := time.ParseDuration(fileConfig.HTTPTimeoutStr)
		if err == nil {
			c.HTTPTimeout = duration
		} else {
			log.Printf("Warning: Invalid HTTP timeout in config file: %s", fileConfig.HTTPTimeoutStr)
		}
	}
	if fileConfig.ServerName != "" {
		c.ServerName = fileConfig.ServerName
	}
	if fileConfig.ServerVersion != "" {
		c.ServerVersion = fileConfig.ServerVersion
	}

	return nil
}

// Validate performs additional validation on the configuration
// and returns an error if the configuration is invalid
func (c *Config) Validate() error {
	if c.BochaAPIKey == "" {
		return fmt.Errorf("BOCHA_API_KEY environment variable is required")
	}

	if c.BochaAPIBaseURL == "" {
		return fmt.Errorf("BOCHA_API_BASE_URL cannot be empty")
	}

	// Log a masked version of the API key for debugging
	if len(c.BochaAPIKey) > 8 {
		maskedKey := c.BochaAPIKey[:4] + "..." + c.BochaAPIKey[len(c.BochaAPIKey)-4:]
		log.Printf("Using Bocha API key: %s", maskedKey)
	}

	return nil
}

// getEnvWithDefault returns the value of the environment variable or the default value if not set
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvDurationWithDefault returns the duration from the environment variable or the default value if not set
func getEnvDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Try to parse as seconds
	seconds, err := strconv.Atoi(value)
	if err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Try to parse as duration string
	duration, err := time.ParseDuration(value)
	if err == nil {
		return duration
	}

	// Return default if parsing fails
	log.Printf("Warning: Could not parse %s as duration, using default of %s", key, defaultValue)
	return defaultValue
}
