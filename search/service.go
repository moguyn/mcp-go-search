package search

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"com.moguyn/mcp-go-search/config"
)

// WebSearchRequest represents the request structure for the Bocha Web Search API
type WebSearchRequest struct {
	Query     string `json:"query"`
	Freshness string `json:"freshness"`
	Count     int    `json:"count"`
	Summary   bool   `json:"summary"`
}

// WebPageResult represents a single web page result from the Bocha Web Search API
type WebPageResult struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	DisplayURL       string `json:"displayUrl"`
	Snippet          string `json:"snippet"`
	SiteName         string `json:"siteName,omitempty"`
	SiteIcon         string `json:"siteIcon,omitempty"`
	DateLastCrawled  string `json:"dateLastCrawled,omitempty"`
	CachedPageURL    any    `json:"cachedPageUrl"`
	Language         any    `json:"language"`
	IsFamilyFriendly any    `json:"isFamilyFriendly"`
	IsNavigational   any    `json:"isNavigational"`
}

// WebPages represents the web pages section of the search response
type WebPages struct {
	WebSearchURL          string          `json:"webSearchUrl"`
	TotalEstimatedMatches int             `json:"totalEstimatedMatches"`
	Value                 []WebPageResult `json:"value"`
	SomeResultsRemoved    bool            `json:"someResultsRemoved"`
}

// ImageResult represents a single image result from the Bocha Web Search API
type ImageResult struct {
	WebSearchURL       any    `json:"webSearchUrl"`
	Name               any    `json:"name"`
	ThumbnailURL       string `json:"thumbnailUrl"`
	DatePublished      any    `json:"datePublished"`
	ContentURL         string `json:"contentUrl"`
	HostPageURL        string `json:"hostPageUrl"`
	ContentSize        any    `json:"contentSize"`
	EncodingFormat     any    `json:"encodingFormat"`
	HostPageDisplayURL string `json:"hostPageDisplayUrl"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	Thumbnail          any    `json:"thumbnail"`
}

// Images represents the images section of the search response
type Images struct {
	ID               any           `json:"id"`
	ReadLink         any           `json:"readLink"`
	WebSearchURL     any           `json:"webSearchUrl"`
	Value            []ImageResult `json:"value"`
	IsFamilyFriendly any           `json:"isFamilyFriendly"`
}

// QueryContext represents the query context section of the search response
type QueryContext struct {
	OriginalQuery string `json:"originalQuery"`
}

// Data represents the data section of the search response
type Data struct {
	Type         string       `json:"_type"`
	QueryContext QueryContext `json:"queryContext"`
	WebPages     WebPages     `json:"webPages"`
	Images       Images       `json:"images,omitempty"`
	Videos       any          `json:"videos"`
}

// WebSearchResponse represents the response structure from the Bocha Web Search API
type WebSearchResponse struct {
	Code  int    `json:"code"`
	LogID string `json:"log_id"`
	Msg   any    `json:"msg"`
	Data  Data   `json:"data"`
}

// Service defines the interface for search operations
type Service interface {
	Search(ctx context.Context, query string, freshness string, count int, summary bool) (*WebSearchResponse, error)
}

// BochaService implements the Service interface for Bocha Web Search API
type BochaService struct {
	apiKey      string
	apiBaseURL  string
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

// NewBochaService creates a new instance of the BochaService
func NewBochaService() *BochaService {
	return NewBochaServiceWithConfig(config.New())
}

// NewBochaServiceWithConfig creates a new instance of the BochaService with the provided configuration
func NewBochaServiceWithConfig(cfg *config.Config) *BochaService {
	// Create a secure transport with modern TLS configuration
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		ForceAttemptHTTP2: true,
		MaxIdleConns:      100,
		IdleConnTimeout:   90 * time.Second,
	}

	// Create a rate limiter that allows 10 requests per second with a burst of 20
	limiter := rate.NewLimiter(rate.Limit(10), 20)

	return &BochaService{
		apiKey:     cfg.BochaAPIKey,
		apiBaseURL: cfg.BochaAPIBaseURL,
		httpClient: &http.Client{
			Timeout:   cfg.HTTPTimeout,
			Transport: transport,
		},
		rateLimiter: limiter,
	}
}

// Search performs a search using the Bocha Web Search API
func (s *BochaService) Search(ctx context.Context, query string, freshness string, count int, summary bool) (*WebSearchResponse, error) {
	// Apply rate limiting
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	// Validate inputs
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Sanitize the query to prevent potential injection attacks
	query = sanitizeQuery(query)

	// Validate freshness parameter if provided
	if freshness != "" && freshness != "noLimit" && freshness != "day" && freshness != "week" && freshness != "month" && freshness != "oneYear" {
		return nil, fmt.Errorf("invalid freshness value: %q, must be one of: noLimit, day, week, month, oneYear", freshness)
	}

	if count < 1 {
		count = 1
	} else if count > 50 {
		count = 50
	}

	// Create the request payload
	reqBody := WebSearchRequest{
		Query:     query,
		Freshness: freshness,
		Count:     count,
		Summary:   summary,
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
	req.Header.Set("User-Agent", "BochaWebSearchMCPServer/1.0")

	// Send the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Bocha API: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body with a size limit to prevent memory exhaustion
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024)) // 10MB limit
	if err != nil {
		return nil, fmt.Errorf("failed to read Bocha API response body: %w", err)
	}

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		// Try to extract error message from response if possible
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
			return nil, fmt.Errorf("bocha api error (status %d): %s", resp.StatusCode, errorResp.Error)
		}

		// Don't return the full response body in case of error to avoid leaking sensitive information
		return nil, fmt.Errorf("bocha api returned status code %d", resp.StatusCode)
	}

	// Parse the response
	var searchResp WebSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse bocha api response: %w", err)
	}

	// Validate response
	if searchResp.Data.WebPages.Value == nil {
		return nil, fmt.Errorf("bocha api returned empty or invalid response")
	}

	return &searchResp, nil
}

// sanitizeQuery performs basic sanitization on the search query
// to prevent potential injection attacks
func sanitizeQuery(query string) string {
	// This is a simple implementation - in a production environment,
	// you might want to use a more sophisticated sanitization library

	// Limit query length to prevent DoS attacks
	const maxQueryLength = 1000
	if len(query) > maxQueryLength {
		query = query[:maxQueryLength]
	}

	return query
}
