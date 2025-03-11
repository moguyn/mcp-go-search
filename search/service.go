package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"com.moguyn/mcp-go-search/config"
)

// Request represents the request structure for the Bocha AI Search API
type Request struct {
	Query     string `json:"query"`
	Freshness string `json:"freshness"`
	Count     int    `json:"count"`
	Answer    bool   `json:"answer"`
	Stream    bool   `json:"stream"`
}

// Result represents a single search result from the Bocha AI Search API
type Result struct {
	Title         string `json:"title"`
	URL           string `json:"url"`
	Description   string `json:"description"`
	DatePublished string `json:"datePublished,omitempty"`
}

// Response represents the response structure from the Bocha AI Search API
type Response struct {
	Results []Result `json:"results"`
	Answer  string   `json:"answer,omitempty"`
}

// Service defines the interface for search operations
type Service interface {
	Search(ctx context.Context, query string, freshness string, count int, answer bool) (*Response, error)
}

// BochaService implements the Service interface for Bocha AI Search API
type BochaService struct {
	apiKey     string
	apiBaseURL string
	httpClient *http.Client
}

// NewBochaService creates a new instance of the BochaService
func NewBochaService() *BochaService {
	return NewBochaServiceWithConfig(config.New())
}

// NewBochaServiceWithConfig creates a new instance of the BochaService with the provided configuration
func NewBochaServiceWithConfig(cfg *config.Config) *BochaService {
	return &BochaService{
		apiKey:     cfg.BochaAPIKey,
		apiBaseURL: cfg.BochaAPIBaseURL,
		httpClient: &http.Client{
			Timeout: cfg.HTTPTimeout,
		},
	}
}

// Search performs a search using the Bocha AI Search API
func (s *BochaService) Search(ctx context.Context, query string, freshness string, count int, answer bool) (*Response, error) {
	// Validate inputs
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	if count < 1 {
		count = 1
	} else if count > 50 {
		count = 50
	}

	// Create the request payload
	reqBody := Request{
		Query:     query,
		Freshness: freshness,
		Count:     count,
		Answer:    answer,
		Stream:    false, // We don't support streaming in this implementation
	}

	// Convert the request to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiBaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	// Send the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Bocha API: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Bocha API response body: %w", err)
	}

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bocha API returned status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var searchResp Response
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse Bocha API response: %w", err)
	}

	return &searchResp, nil
}
