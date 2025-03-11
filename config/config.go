package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	// API configuration
	BochaAPIKey     string
	BochaAPIBaseURL string
	HTTPTimeout     time.Duration

	// Server configuration
	ServerName    string
	ServerVersion string
}

// New creates a new configuration with values from environment variables
func New() *Config {
	config := &Config{
		// Default values
		BochaAPIKey:     os.Getenv("BOCHA_API_KEY"),
		BochaAPIBaseURL: getEnvWithDefault("BOCHA_API_BASE_URL", "https://api.bochaai.com/v1/ai-search"),
		HTTPTimeout:     getEnvDurationWithDefault("HTTP_TIMEOUT", 10*time.Second),
		ServerName:      getEnvWithDefault("SERVER_NAME", "Bocha AI Search Server"),
		ServerVersion:   getEnvWithDefault("SERVER_VERSION", "1.0.0"),
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

// Validate performs additional validation on the configuration
// and returns an error if the configuration is invalid
func (c *Config) Validate() error {
	if c.BochaAPIKey == "" {
		return fmt.Errorf("BOCHA_API_KEY environment variable is required")
	}

	if c.BochaAPIBaseURL == "" {
		return fmt.Errorf("BOCHA_API_BASE_URL cannot be empty")
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
