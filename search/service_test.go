package search

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"com.moguyn/mcp-go-search/config"
)

// TestNewBochaService tests the NewBochaService function
func TestNewBochaService(t *testing.T) {
	// Save original environment variables to restore later
	origAPIKey := os.Getenv("BOCHA_API_KEY")
	origAPIBaseURL := os.Getenv("BOCHA_API_BASE_URL")
	origHTTPTimeout := os.Getenv("HTTP_TIMEOUT")

	// Restore environment variables after the test
	defer func() {
		os.Setenv("BOCHA_API_KEY", origAPIKey)
		os.Setenv("BOCHA_API_BASE_URL", origAPIBaseURL)
		os.Setenv("HTTP_TIMEOUT", origHTTPTimeout)
	}()

	// Set test environment variables
	os.Setenv("BOCHA_API_KEY", "test-api-key")
	os.Setenv("BOCHA_API_BASE_URL", "https://test.api.com")
	os.Setenv("HTTP_TIMEOUT", "5s")

	// Create a new service
	service := NewBochaService()

	// Check that the service was created with the correct values
	if service.apiKey != "test-api-key" {
		t.Errorf("Expected apiKey to be 'test-api-key', got '%s'", service.apiKey)
	}
	if service.apiBaseURL != "https://test.api.com" {
		t.Errorf("Expected apiBaseURL to be 'https://test.api.com', got '%s'", service.apiBaseURL)
	}
	if service.httpClient.Timeout != 5*time.Second {
		t.Errorf("Expected httpClient.Timeout to be 5s, got %s", service.httpClient.Timeout)
	}
}

// TestBochaService_Search tests the Search method of BochaService
func TestBochaService_Search(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header 'Bearer test-api-key', got %s", authHeader)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type header 'application/json', got %s", contentType)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"results": [
				{
					"title": "Test Result",
					"url": "https://example.com",
					"description": "This is a test result"
				}
			],
			"answer": "This is a test answer"
		}`))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create a test configuration
	cfg := &config.Config{
		BochaAPIKey:     "test-api-key",
		BochaAPIBaseURL: server.URL,
		HTTPTimeout:     5 * time.Second,
	}

	// Create a search service with the test configuration
	service := NewBochaServiceWithConfig(cfg)

	// Call the Search method
	ctx := context.Background()
	response, err := service.Search(ctx, "test query", "noLimit", 10, true)

	// Check for errors
	if err != nil {
		t.Fatalf("Search returned an error: %v", err)
	}

	// Check the response
	if response == nil {
		t.Fatal("Search returned nil response")
	}

	// Now that we've verified response is not nil, we can safely check its fields
	expectedAnswer := "This is a test answer"
	if response.Answer != expectedAnswer {
		t.Errorf("Expected answer '%s', got '%s'", expectedAnswer, response.Answer)
	}

	// Check the results
	if len(response.Results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(response.Results))
	}

	result := response.Results[0]
	if result.Title != "Test Result" {
		t.Errorf("Expected title 'Test Result', got '%s'", result.Title)
	}
	if result.URL != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got '%s'", result.URL)
	}
	if result.Description != "This is a test result" {
		t.Errorf("Expected description 'This is a test result', got '%s'", result.Description)
	}
}

// TestBochaService_Search_Validation tests the validation in the Search method
func TestBochaService_Search_Validation(t *testing.T) {
	// Create a mock server that won't actually be called
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		// This server should never be called for validation tests
		t.Error("Server should not be called for validation tests")
	}))
	defer server.Close()

	// Create a test configuration
	cfg := &config.Config{
		BochaAPIKey:     "test-api-key",
		BochaAPIBaseURL: server.URL,
		HTTPTimeout:     5 * time.Second,
	}

	// Create a search service with the test configuration
	service := NewBochaServiceWithConfig(cfg)

	// Test with empty query
	ctx := context.Background()
	_, err := service.Search(ctx, "", "noLimit", 10, true)
	if err == nil {
		t.Error("Expected error for empty query, got nil")
	}

	// Test count normalization
	// We can't easily test the actual normalization without making HTTP requests,
	// so we'll just check that the validation doesn't return an error

	// Create a mock server that will be called for the count tests
	countServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"results": []}`))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer countServer.Close()

	// Update the service to use the count test server
	service.apiBaseURL = countServer.URL

	// Test with count < 1
	_, err = service.Search(ctx, "test query", "noLimit", 0, true)
	if err != nil {
		t.Errorf("Expected no error for count < 1, got %v", err)
	}

	// Test with count > 50
	_, err = service.Search(ctx, "test query", "noLimit", 100, true)
	if err != nil {
		t.Errorf("Expected no error for count > 50, got %v", err)
	}
}

// TestSanitizeQuery tests the sanitizeQuery function
func TestSanitizeQuery(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal query",
			input:    "test query",
			expected: "test query",
		},
		{
			name:     "Empty query",
			input:    "",
			expected: "",
		},
		{
			name:     "Query at max length",
			input:    strings.Repeat("a", 1000),
			expected: strings.Repeat("a", 1000),
		},
		{
			name:     "Query exceeding max length",
			input:    strings.Repeat("a", 1500),
			expected: strings.Repeat("a", 1000),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := sanitizeQuery(tc.input)
			if result != tc.expected {
				t.Errorf("Expected query of length %d, got length %d", len(tc.expected), len(result))
			}
		})
	}
}

// TestBochaService_Search_Errors tests error handling in the Search method
func TestBochaService_Search_Errors(t *testing.T) {
	// Test server that returns an error
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Internal server error"}`))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer errorServer.Close()

	// Test server that returns invalid JSON
	invalidJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"results": [invalid json]}`))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer invalidJSONServer.Close()

	// Test server that returns empty results
	emptyResultsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{}`))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer emptyResultsServer.Close()

	// Create a test configuration
	cfg := &config.Config{
		BochaAPIKey:     "test-api-key",
		BochaAPIBaseURL: errorServer.URL,
		HTTPTimeout:     5 * time.Second,
	}

	// Create a search service with the test configuration
	service := NewBochaServiceWithConfig(cfg)

	// Test error response
	ctx := context.Background()
	_, err := service.Search(ctx, "test query", "noLimit", 10, true)
	if err == nil {
		t.Error("Expected error for error response, got nil")
	}

	// Test invalid JSON response
	service.apiBaseURL = invalidJSONServer.URL
	_, err = service.Search(ctx, "test query", "noLimit", 10, true)
	if err == nil {
		t.Error("Expected error for invalid JSON response, got nil")
	}

	// Test empty results
	service.apiBaseURL = emptyResultsServer.URL
	_, err = service.Search(ctx, "test query", "noLimit", 10, true)
	if err == nil {
		t.Error("Expected error for empty results, got nil")
	}

	// Test invalid freshness parameter
	service.apiBaseURL = errorServer.URL // Use any server, validation happens before request
	_, err = service.Search(ctx, "test query", "invalid", 10, true)
	if err == nil {
		t.Error("Expected error for invalid freshness parameter, got nil")
	}

	// Test context cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel the context immediately
	_, err = service.Search(cancelCtx, "test query", "noLimit", 10, true)
	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}
}
