.PHONY: build run test lint clean help release release-snapshot run-config sec-scan sec-deps sec-tidy

# Binary name
BINARY_NAME=mcp-search-server

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Default target
.DEFAULT_GOAL := help

# Build the application
build:
	@echo "Building..."
	@$(GOBUILD) -o $(BINARY_NAME)

# Run the application (requires API key)
run: build
	@if [ -z "$(API_KEY)" ]; then \
		echo "Usage: make run API_KEY=your-api-key-here"; \
		echo "Example: make run API_KEY=abc123"; \
		exit 1; \
	fi
	@echo "Running with API key..."
	@BOCHA_API_KEY=$(API_KEY) ./$(BINARY_NAME)

# Run with custom configuration
run-custom: build
	@if [ -z "$(API_KEY)" ]; then \
		echo "Usage: make run-custom API_KEY=your-api-key-here [API_BASE_URL=url] [HTTP_TIMEOUT=timeout] [SERVER_NAME=name] [SERVER_VERSION=version]"; \
		exit 1; \
	fi
	@echo "Running with custom configuration..."
	@BOCHA_API_KEY=$(API_KEY) \
	 BOCHA_API_BASE_URL=$(if $(API_BASE_URL),$(API_BASE_URL),https://api.bochaai.com/v1/ai-search) \
	 HTTP_TIMEOUT=$(if $(HTTP_TIMEOUT),$(HTTP_TIMEOUT),10s) \
	 SERVER_NAME=$(if $(SERVER_NAME),$(SERVER_NAME),"Bocha AI Search Server") \
	 SERVER_VERSION=$(if $(SERVER_VERSION),$(SERVER_VERSION),0.0.1) \
	 ./$(BINARY_NAME)

# Run with config file
run-config: build
	@if [ -z "$(CONFIG_FILE)" ]; then \
		echo "Usage: make run-config CONFIG_FILE=path/to/config.yaml"; \
		echo "Example: make run-config CONFIG_FILE=./config.yaml"; \
		exit 1; \
	fi
	@echo "Running with config file: $(CONFIG_FILE)..."
	@CONFIG_FILE=$(CONFIG_FILE) ./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html

# Update dependencies
deps:
	@echo "Updating dependencies..."
	@$(GOMOD) tidy

# Create a new release using GoReleaser
release:
	@echo "Creating a new release..."
	@goreleaser release --clean

# Create a snapshot release for testing
release-snapshot:
	@echo "Creating a snapshot release..."
	@goreleaser release --snapshot --clean

# Run security scans
sec-scan: sec-deps sec-tidy
	@echo "Running security scans..."
	@govulncheck ./...

# Install security scanning dependencies
sec-deps:
	@echo "Installing security scanning dependencies..."
	@go install golang.org/x/vuln/cmd/govulncheck@latest

# Check for unused dependencies that might introduce vulnerabilities
sec-tidy:
	@echo "Checking for unused dependencies..."
	@go mod tidy
	@echo "Verifying dependencies..."
	@go mod verify

# Show help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run                  Run the server with API key from environment"
	@echo "  run-custom           Run with custom configuration from environment"
	@echo "  run-config           Run with configuration from a file"
	@echo "  build                Build the server binary"
	@echo "  test                 Run tests"
	@echo "  cover                Run tests with coverage"
	@echo "  cover-html           Generate HTML coverage report"
	@echo "  lint                 Run linter"
	@echo "  deps                 Update dependencies"
	@echo "  clean                Clean build artifacts"
	@echo "  sec-scan             Run security vulnerability scans"
	@echo "  sec-deps             Install security scanning dependencies"
	@echo "  sec-tidy             Check for unused dependencies"
	@echo ""
	@echo "Environment variables:"
	@echo "  API_KEY              Bocha AI API key"
	@echo "  API_BASE_URL         Bocha AI API base URL"
	@echo "  HTTP_TIMEOUT         HTTP timeout duration"
	@echo "  SERVER_NAME          Server name"
	@echo "  SERVER_VERSION       Server version"
	@echo "  CONFIG_FILE          Path to configuration file"
	@echo ""
	@echo "Examples:"
	@echo "  make run API_KEY=your-api-key-here"
	@echo "  make run-custom API_KEY=your-api-key-here API_BASE_URL=https://custom-url.com HTTP_TIMEOUT=5s"
	@echo "  make run-config CONFIG_FILE=./config.yaml"
	@echo "  make release-snapshot" 