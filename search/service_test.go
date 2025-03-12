package search

import (
	"context"
	"encoding/json"
	"io"
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
	// Mock server response
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

		// Read and parse request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		var req WebSearchRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		// Check request parameters
		if req.Query != "test query" {
			t.Errorf("Expected query 'test query', got %s", req.Query)
		}

		if req.Freshness != "noLimit" {
			t.Errorf("Expected freshness 'noLimit', got %s", req.Freshness)
		}

		if req.Count != 10 {
			t.Errorf("Expected count 10, got %d", req.Count)
		}

		if !req.Summary {
			t.Errorf("Expected summary to be true")
		}

		// Return a mock response
		resp := WebSearchResponse{
			Results: []WebSearchResult{
				{
					Title:         "Test Result 1",
					URL:           "https://example.com/1",
					Description:   "This is test result 1",
					DatePublished: "2023-01-01T12:00:00Z",
				},
				{
					Title:         "Test Result 2",
					URL:           "https://example.com/2",
					Description:   "This is test result 2",
					DatePublished: "2023-01-02T12:00:00Z",
				},
			},
			Summary: "This is a summary of the search results.",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create a configuration with the test server URL
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

	if len(response.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(response.Results))
	}

	if response.Results[0].Title != "Test Result 1" {
		t.Errorf("Expected first result title 'Test Result 1', got %s", response.Results[0].Title)
	}

	if response.Results[1].URL != "https://example.com/2" {
		t.Errorf("Expected second result URL 'https://example.com/2', got %s", response.Results[1].URL)
	}

	if response.Summary != "This is a summary of the search results." {
		t.Errorf("Expected summary 'This is a summary of the search results.', got %s", response.Summary)
	}
}

// TestBochaService_Search_Validation tests the validation in the Search method
func TestBochaService_Search_Validation(t *testing.T) {
	// Create a mock server that returns a valid response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"results": [
				{
					"title": "Test Result",
					"url": "https://example.com",
					"description": "This is a test result"
				}
			],
			"summary": "This is a test summary"
		}`))
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
	ctx := context.Background()

	// Test empty query
	_, err := service.Search(ctx, "", "noLimit", 10, true)
	if err == nil {
		t.Error("Expected error for empty query, got nil")
	} else if err.Error() != "search query cannot be empty" {
		t.Errorf("Expected error message 'search query cannot be empty', got '%s'", err.Error())
	}

	// Test count validation (too low)
	_, err = service.Search(ctx, "test query", "noLimit", 0, true)
	if err != nil {
		t.Errorf("Expected no error for count 0 (should be adjusted to 1), got %v", err)
	}

	// Test count validation (too high)
	_, err = service.Search(ctx, "test query", "noLimit", 100, true)
	if err != nil {
		t.Errorf("Expected no error for count 100 (should be adjusted to 50), got %v", err)
	}

	// Test freshness validation
	_, err = service.Search(ctx, "test query", "invalid", 10, true)
	if err == nil {
		t.Error("Expected error for invalid freshness, got nil")
	} else if err.Error() != "invalid freshness value: \"invalid\", must be one of: noLimit, day, week, month, oneYear" {
		t.Errorf("Expected error message about invalid freshness, got '%s'", err.Error())
	}

	// Test context cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel the context immediately
	_, err = service.Search(cancelCtx, "test query", "noLimit", 10, true)
	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
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
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Test error message"}`))
	}))
	defer errorServer.Close()

	// Create a test configuration with the error server
	errorCfg := &config.Config{
		BochaAPIKey:     "test-api-key",
		BochaAPIBaseURL: errorServer.URL,
		HTTPTimeout:     5 * time.Second,
	}

	// Create a search service with the error configuration
	errorService := NewBochaServiceWithConfig(errorCfg)

	// Test with error response
	ctx := context.Background()
	_, err := errorService.Search(ctx, "test query", "noLimit", 10, true)
	if err == nil {
		t.Error("Expected error for error response, got nil")
	} else if err.Error() != "bocha api error (status 400): Test error message" {
		t.Errorf("Expected error message 'bocha api error (status 400): Test error message', got '%s'", err.Error())
	}

	// Test with invalid JSON response
	invalidJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json`))
	}))
	defer invalidJSONServer.Close()

	// Create a test configuration with the invalid JSON server
	invalidJSONCfg := &config.Config{
		BochaAPIKey:     "test-api-key",
		BochaAPIBaseURL: invalidJSONServer.URL,
		HTTPTimeout:     5 * time.Second,
	}

	// Create a search service with the invalid JSON configuration
	invalidJSONService := NewBochaServiceWithConfig(invalidJSONCfg)

	// Test with invalid JSON response
	_, err = invalidJSONService.Search(ctx, "test query", "noLimit", 10, true)
	if err == nil {
		t.Error("Expected error for invalid JSON response, got nil")
	}

	// Test with empty results
	emptyResultsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"results": []}`))
	}))
	defer emptyResultsServer.Close()

	// Create a test configuration with the empty results server
	emptyResultsCfg := &config.Config{
		BochaAPIKey:     "test-api-key",
		BochaAPIBaseURL: emptyResultsServer.URL,
		HTTPTimeout:     5 * time.Second,
	}

	// Create a search service with the empty results configuration
	emptyResultsService := NewBochaServiceWithConfig(emptyResultsCfg)

	// Test with empty results
	_, err = emptyResultsService.Search(ctx, "test query", "noLimit", 10, true)
	if err != nil {
		t.Errorf("Expected no error for empty results, got %v", err)
	}
}
