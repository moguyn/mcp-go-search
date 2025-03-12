# Bocha AI Search MCP Server

[![CI](https://github.com/moguyn/mcp-go-search/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/mcp-go-search/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/moguyn/mcp-go-search/graph/badge.svg?token=J7QC7MFP0D)](https://codecov.io/gh/moguyn/mcp-go-search)
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
- Enhanced security features:
  - API key protection
  - Input validation and sanitization
  - Rate limiting to prevent abuse
  - Secure error handling
  - TLS 1.2+ support

## Prerequisites

- Go 1.20 or higher
- A Bocha AI API key
- Make (for using the Makefile)

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

### Quick Start with Make

The easiest way to run the server is using the provided Makefile:

```bash
make run API_KEY=your-bocha-api-key-here
```

For a list of all available make targets:

```bash
make help
```

### Running with Custom Configuration

You can customize the server configuration:

```bash
make run-custom API_KEY=your-api-key-here API_BASE_URL=https://custom-url.com HTTP_TIMEOUT=5s SERVER_NAME="Custom Server" SERVER_VERSION="2.0.0"
```

### Running with a Configuration File

You can also use a YAML configuration file:

1. Create a config.yaml file (see config.yaml.example for reference):
   ```yaml
   # Bocha AI Search Server Configuration
   
   # API configuration
   bocha_api_key: "your-api-key-here"
   bocha_api_base_url: "https://api.bochaai.com/v1/ai-search"
   http_timeout: "10s"
   
   # Server configuration
   server_name: "Bocha AI Search Server"
   server_version: "1.0.0"
   ```

2. Run the server with the config file:
   ```bash
   make run-config CONFIG_FILE=./config.yaml
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
make test
```

### Linting

This project uses golangci-lint for code quality. To run the linter:

```bash
make lint
```

### Cleaning Build Artifacts

```bash
make clean
```

### Updating Dependencies

```bash
make deps
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

## Contributing

We welcome contributions to the Bocha AI Search MCP Server! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) file for detailed guidelines on how to contribute.

Here's a quick overview:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Commit your changes** (`git commit -m 'Add some amazing feature'`)
4. **Push to the branch** (`git push origin feature/amazing-feature`)
5. **Open a Pull Request**

Please ensure your code follows our style guidelines and passes all tests before submitting a PR.

## License

MIT 

## Security

This project implements several security features to protect your API keys and prevent abuse:

- API key masking in logs
- Input validation and sanitization
- Rate limiting
- Request timeouts
- Secure error handling to prevent information leakage
- TLS 1.2+ support

For detailed security guidelines, please see the [SECURITY.md](SECURITY.md) file.

### API Key Security

It is strongly recommended to use environment variables for your API key rather than configuration files:

```bash
export BOCHA_API_KEY="your-api-key-here"
```

If you must use a configuration file, ensure it is not committed to version control and has restricted permissions:

```bash
chmod 600 config.yaml
``` 