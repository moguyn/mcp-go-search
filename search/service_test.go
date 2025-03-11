package search

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"com.moguyn/mcp-go-search/config"
)

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
