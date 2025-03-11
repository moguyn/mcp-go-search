#!/bin/bash

# Check if API key is provided
if [ -z "$1" ]; then
  echo "Usage: ./run.sh <BOCHA_API_KEY>"
  echo "Example: ./run.sh your-api-key-here"
  exit 1
fi

# Set environment variables
export BOCHA_API_KEY="$1"

# Optional environment variables (uncomment to override defaults)
# export BOCHA_API_BASE_URL="https://api.bochaai.com/v1/ai-search"
# export HTTP_TIMEOUT="10s"
# export SERVER_NAME="Bocha AI Search Server"
# export SERVER_VERSION="1.0.0"

# Build and run the server
go build -o mcp-search-server
./mcp-search-server 