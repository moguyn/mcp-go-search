package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"com.moguyn/mcp-go-search/search"
)

// MockSearchService is a mock implementation of the search.Service interface
type MockSearchService struct {
	SearchFunc func(ctx context.Context, query string, freshness string, count int, summary bool) (*search.WebSearchResponse, error)
}

// Search calls the mock SearchFunc
func (m *MockSearchService) Search(ctx context.Context, query string, freshness string, count int, summary bool) (*search.WebSearchResponse, error) {
	return m.SearchFunc(ctx, query, freshness, count, summary)
}

func TestNewSearchTool(t *testing.T) {
	mockService := &MockSearchService{}
	tool := NewSearchTool(mockService)

	if tool == nil {
		t.Fatal("NewSearchTool returned nil")
	}

	if tool.searchService != mockService {
		t.Errorf("Expected searchService to be %v, got %v", mockService, tool.searchService)
	}
}

func TestDefinition(t *testing.T) {
	mockService := &MockSearchService{}
	tool := NewSearchTool(mockService)
	definition := tool.Definition()

	// Check tool name
	if definition.Name != "search" {
		t.Errorf("Expected tool name 'search', got '%s'", definition.Name)
	}

	// Check that InputSchema is not empty
	if len(definition.InputSchema.Properties) == 0 {
		t.Error("Expected InputSchema to have properties")
	}

	// Check required parameters
	requiredParams := definition.InputSchema.Required
	if len(requiredParams) == 0 || requiredParams[0] != "query" {
		t.Error("Expected 'query' to be a required parameter")
	}
}

func TestHandler(t *testing.T) {
	// Create test cases
	testCases := []struct {
		name           string
		request        mcp.CallToolRequest
		mockResponse   *search.WebSearchResponse
		mockError      error
		expectedResult *mcp.CallToolResult
		expectedError  error
	}{
		{
			name: "Basic search",
			request: mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{
						"query": "test query",
					},
				},
			},
			mockResponse: &search.WebSearchResponse{
				Results: []search.WebSearchResult{
					{
						Title:       "Test Result",
						URL:         "https://example.com",
						Description: "This is a test result",
					},
				},
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "Search with all parameters",
			request: mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{
						"query":     "test query",
						"freshness": "day",
						"count":     float64(5),
						"summary":   true,
					},
				},
			},
			mockResponse: &search.WebSearchResponse{
				Results: []search.WebSearchResult{
					{
						Title:         "Test Result",
						URL:           "https://example.com",
						Description:   "This is a test result",
						DatePublished: "2023-01-01T12:00:00Z",
					},
				},
				Summary: "This is a test summary",
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "Missing query parameter",
			request: mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{},
				},
			},
			mockResponse:  nil,
			mockError:     nil,
			expectedError: nil, // We handle this error internally
		},
		{
			name: "Empty query parameter",
			request: mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{
						"query": "",
					},
				},
			},
			mockResponse:  nil,
			mockError:     nil,
			expectedError: nil, // We handle this error internally
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock search service
			mockService := &MockSearchService{}
			mockService.SearchFunc = func(ctx context.Context, query string, freshness string, count int, summary bool) (*search.WebSearchResponse, error) {
				return tc.mockResponse, tc.mockError
			}

			// Create the search tool
			tool := NewSearchTool(mockService)
			handler := tool.Handler()

			// Call the handler
			result, err := handler(context.Background(), tc.request)

			// Check error
			if (err == nil && tc.expectedError != nil) || (err != nil && tc.expectedError == nil) {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}

			// For error cases, check that result indicates an error
			if tc.mockError != nil && result != nil {
				if !result.IsError {
					t.Error("Expected IsError to be true for error case")
				}
			}

			// For missing/empty query cases
			query, ok := tc.request.Params.Arguments["query"].(string)
			if !ok || query == "" {
				if result == nil || !result.IsError {
					t.Error("Expected IsError to be true for missing/empty query")
				}
			}

			// For successful cases with mock response
			if tc.mockResponse != nil && tc.mockError == nil && result != nil {
				// Check that result doesn't indicate an error
				if result.IsError {
					t.Error("Expected IsError to be false for successful case")
				}

				// Check that result content contains expected data
				if len(result.Content) == 0 {
					t.Error("Expected non-empty content in result")
				} else {
					// Get the text content
					var resultText string
					for _, content := range result.Content {
						if textContent, ok := content.(mcp.TextContent); ok {
							resultText += textContent.Text
						}
					}

					// Check that result text contains the title of the first result
					if len(tc.mockResponse.Results) > 0 {
						if !strings.Contains(resultText, tc.mockResponse.Results[0].Title) {
							t.Errorf("Expected result text to contain '%s'", tc.mockResponse.Results[0].Title)
						}
					}

					// Check if summary is included when requested
					if tc.request.Params.Arguments["summary"] == true && tc.mockResponse.Summary != "" {
						if !strings.Contains(resultText, "Summary:") {
							t.Errorf("Expected result to contain summary, but it didn't: %s", resultText)
						}
						if !strings.Contains(resultText, tc.mockResponse.Summary) {
							t.Errorf("Expected result to contain summary text '%s', but it didn't: %s", tc.mockResponse.Summary, resultText)
						}
					}
				}
			}
		})
	}
}

func TestFormatFreshness(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"day", "Past 24 hours"},
		{"week", "Past week"},
		{"month", "Past month"},
		{"oneYear", "Past year"},
		{"noLimit", "No time limit"},
		{"", "No time limit"},
		{"invalid", "No time limit"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := formatFreshness(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"2023-01-01T12:00:00Z", "January 1, 2023"},
		{"2023-01-01", "January 1, 2023"},
		{"invalid", "invalid"}, // Should return original string for invalid format
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := formatDate(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestSanitizeErrorMessage(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No sensitive information",
			input:    "Simple error message",
			expected: "Simple error message",
		},
		{
			name:     "Contains API key in Bearer token",
			input:    "Error with Authorization: Bearer abc123secret456token789",
			expected: "Error with Authorization: Bearer [REDACTED]",
		},
		{
			name:     "Contains URL with http",
			input:    "Failed to connect to http://api.example.com/v1/endpoint",
			expected: "Failed to connect to [URL REDACTED]",
		},
		{
			name:     "Contains URL with https",
			input:    "Failed to connect to https://api.example.com/v1/endpoint",
			expected: "Failed to connect to [URL REDACTED]",
		},
		{
			name:     "Contains both Bearer token and URL",
			input:    "Error with Authorization: Bearer abc123secret456token789 when connecting to https://api.example.com",
			expected: "Error with Authorization: Bearer [REDACTED] when connecting to [URL REDACTED]",
		},
		{
			name:     "Bearer token at end of string",
			input:    "Error with Authorization: Bearer abc123secret456token789",
			expected: "Error with Authorization: Bearer [REDACTED]",
		},
		{
			name:     "URL at end of string",
			input:    "Failed to connect to https://api.example.com",
			expected: "Failed to connect to [URL REDACTED]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := sanitizeErrorMessage(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}
