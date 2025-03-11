.PHONY: build run test lint clean help

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
	 SERVER_VERSION=$(if $(SERVER_VERSION),$(SERVER_VERSION),1.0.0) \
	 ./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)

# Update dependencies
deps:
	@echo "Updating dependencies..."
	@$(GOMOD) tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Build and run the application (requires API_KEY)"
	@echo "  run-custom   - Run with custom configuration options"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linter"
	@echo "  clean        - Remove build artifacts"
	@echo "  deps         - Update dependencies"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make run API_KEY=your-api-key-here"
	@echo "  make run-custom API_KEY=your-api-key-here API_BASE_URL=https://custom-url.com HTTP_TIMEOUT=5s" 