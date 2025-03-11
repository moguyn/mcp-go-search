# Bocha AI Search MCP Server

[![CI](https://github.com/yourusername/mcp-go-search/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/mcp-go-search/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/mcp-go-search)](https://goreportcard.com/report/github.com/yourusername/mcp-go-search)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Model Context Protocol (MCP) server that provides web search capabilities using the Bocha AI Search API.

## Overview

This server implements the [Model Context Protocol (MCP)](https://modelcontextprotocol.io) to expose the Bocha AI Search API as a tool that can be used by LLM applications. It allows LLMs to search the web and retrieve relevant information.

## Features

- Web search using Bocha AI Search API
- Configurable search parameters (freshness, result count)
- Optional answer generation based on search results
- Clean, formatted search results
- CI/CD with GitHub Actions

## Prerequisites

- Go 1.20 or higher
- A Bocha AI API key

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/mcp-go-search.git
   cd mcp-go-search
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

### Quick Start with run.sh

The easiest way to run the server is using the provided script:

```bash
./run.sh your-bocha-api-key-here
```

### Manual Configuration and Running

1. Set your Bocha AI API key as an environment variable:
   ```bash
   export BOCHA_API_KEY="your-api-key-here"
   ```

2. Build and run the server:
   ```bash
   go build -o mcp-search-server
   ./mcp-search-server
   ```

The server communicates via standard input/output, following the MCP protocol.

### Connecting to an LLM Application

To use this server with an LLM application that supports MCP:

1. Start the server
2. Connect your LLM application to the server's stdin/stdout
3. The search tool will be available to the LLM

### Search Tool Parameters

The search tool accepts the following parameters:

- `query` (string, required): The search query
- `freshness` (string, optional): Filter results by freshness - "noLimit", "day", "week", or "month"
- `count` (number, optional): Number of results to return (1-50)
- `answer` (boolean, optional): Whether to generate an answer based on search results

## Example

Here's an example of how an LLM might use the search tool:

```json
{
  "method": "call_tool",
  "params": {
    "name": "search",
    "arguments": {
      "query": "latest news about artificial intelligence",
      "freshness": "day",
      "count": 5,
      "answer": true
    }
  }
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Linting

This project uses golangci-lint for code quality. To run the linter:

```bash
golangci-lint run
```

### CI/CD

This project uses GitHub Actions for continuous integration:
- Runs golangci-lint for code quality
- Executes all tests
- Builds the application

The workflow runs on every push to main and on pull requests.

## Architecture

The server follows SOLID principles:

- **Single Responsibility**: Each component has a single responsibility
- **Open/Closed**: The design is open for extension but closed for modification
- **Liskov Substitution**: Components can be replaced with their subtypes
- **Interface Segregation**: Interfaces are specific to client needs
- **Dependency Inversion**: High-level modules don't depend on low-level modules

## License

MIT 