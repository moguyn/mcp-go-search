package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"com.moguyn/mcp-go-search/search"
)

// SearchTool provides the search functionality as an MCP tool
type SearchTool struct {
	searchService search.Service
}

// NewSearchTool creates a new search tool with the provided search service
func NewSearchTool(searchService search.Service) *SearchTool {
	return &SearchTool{
		searchService: searchService,
	}
}

// Definition returns the MCP tool definition
func (t *SearchTool) Definition() mcp.Tool {
	return mcp.NewTool("search",
		mcp.WithDescription("Search the web using Bocha AI Search API"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search query"),
		),
		mcp.WithString("freshness",
			mcp.Description("Filter results by freshness (noLimit, day, week, month)"),
			mcp.Enum("noLimit", "day", "week", "month"),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of results to return (1-50)"),
		),
		mcp.WithBoolean("answer",
			mcp.Description("Whether to generate an answer based on search results"),
		),
	)
}

// Handler returns the MCP tool handler function
func (t *SearchTool) Handler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create a timeout context to prevent long-running searches
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// Extract parameters from the request
		query, ok := request.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("query parameter is required and must be a string"), nil
		}

		// Validate query length to prevent abuse
		if len(query) > 1000 {
			return mcp.NewToolResultError("query is too long (maximum 1000 characters)"), nil
		}

		// Extract optional parameters with defaults
		freshness := "noLimit"
		if f, ok := request.Params.Arguments["freshness"].(string); ok && f != "" {
			// Validate freshness parameter
			if f != "noLimit" && f != "day" && f != "week" && f != "month" {
				return mcp.NewToolResultError(fmt.Sprintf("invalid freshness value: %q, must be one of: noLimit, day, week, month", f)), nil
			}
			freshness = f
		}

		count := 10
		if c, ok := request.Params.Arguments["count"].(float64); ok {
			count = int(c)
			// Ensure count is within valid range
			if count < 1 {
				count = 1
			} else if count > 50 {
				count = 50
			}
		}

		answer := false
		if a, ok := request.Params.Arguments["answer"].(bool); ok {
			answer = a
		}

		// Perform the search
		response, err := t.searchService.Search(ctx, query, freshness, count, answer)
		if err != nil {
			// Handle context cancellation
			if ctx.Err() == context.DeadlineExceeded {
				return mcp.NewToolResultError("Search timed out after 30 seconds"), nil
			}

			// Sanitize error message to prevent leaking sensitive information
			errMsg := sanitizeErrorMessage(err.Error())
			return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", errMsg)), nil
		}

		// Format the results
		var resultBuilder strings.Builder

		// Add search metadata
		resultBuilder.WriteString(fmt.Sprintf("Search Query: \"%s\"\n", query))
		resultBuilder.WriteString(fmt.Sprintf("Freshness: %s\n", formatFreshness(freshness)))
		resultBuilder.WriteString(fmt.Sprintf("Results: %d\n\n", len(response.Results)))

		// Add answer if available
		if answer && response.Answer != "" {
			resultBuilder.WriteString("Answer:\n")
			resultBuilder.WriteString(response.Answer)
			resultBuilder.WriteString("\n\n")
		}

		// Add search results
		resultBuilder.WriteString("Search Results:\n")
		resultBuilder.WriteString("==============\n\n")

		for i, result := range response.Results {
			resultBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, result.Title))
			resultBuilder.WriteString(fmt.Sprintf("   URL: %s\n", result.URL))

			// Add date published if available
			if result.DatePublished != "" {
				resultBuilder.WriteString(fmt.Sprintf("   Published: %s\n", formatDate(result.DatePublished)))
			}

			resultBuilder.WriteString(fmt.Sprintf("   %s\n\n", result.Description))
		}

		return mcp.NewToolResultText(resultBuilder.String()), nil
	}
}

// formatFreshness returns a human-readable string for the freshness parameter
func formatFreshness(freshness string) string {
	switch freshness {
	case "day":
		return "Past 24 hours"
	case "week":
		return "Past week"
	case "month":
		return "Past month"
	default:
		return "No time limit"
	}
}

// formatDate attempts to format the date in a more readable format
func formatDate(dateStr string) string {
	// Try to parse the date
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	} {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Format("January 2, 2006")
		}
	}

	// Return the original string if parsing fails
	return dateStr
}

// sanitizeErrorMessage removes potentially sensitive information from error messages
func sanitizeErrorMessage(errMsg string) string {
	// Remove any API keys that might be in the error message
	// This is a simple implementation - in a production environment,
	// you might want to use a more sophisticated approach
	if strings.Contains(errMsg, "Bearer ") {
		parts := strings.Split(errMsg, "Bearer ")
		if len(parts) > 1 {
			// Find the end of the token
			tokenEnd := strings.IndexAny(parts[1], " \t\n\r\",;:)")
			if tokenEnd != -1 {
				parts[1] = "[REDACTED]" + parts[1][tokenEnd:]
				errMsg = strings.Join(parts, "Bearer ")
			} else {
				// If we can't find the end of the token, it might be at the end of the string
				parts[1] = "[REDACTED]"
				errMsg = strings.Join(parts, "Bearer ")
			}
		}
	}

	// Remove any URLs that might contain sensitive information
	if strings.Contains(errMsg, "http") {
		// Simple regex-like replacement for URLs
		for _, prefix := range []string{"http://", "https://"} {
			if idx := strings.Index(errMsg, prefix); idx != -1 {
				start := idx
				end := start + len(prefix)
				// Find the end of the URL
				for end < len(errMsg) && !strings.ContainsAny(string(errMsg[end]), " \t\n\r\",;:)") {
					end++
				}
				if end > start+len(prefix) {
					errMsg = errMsg[:start] + "[URL REDACTED]" + errMsg[end:]
				}
			}
		}
	}

	return errMsg
}
