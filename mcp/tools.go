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
		// Extract parameters from the request
		query, ok := request.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("query parameter is required and must be a string"), nil
		}

		// Extract optional parameters with defaults
		freshness := "noLimit"
		if f, ok := request.Params.Arguments["freshness"].(string); ok && f != "" {
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
			return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
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
