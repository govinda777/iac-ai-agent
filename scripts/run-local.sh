#!/bin/bash

# Run IaC AI Agent locally

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check if binary exists
if [ ! -f bin/iac-ai-agent ]; then
    echo "Binary not found. Building..."
    make build
fi

echo "🚀 Starting IaC AI Agent..."
echo "📍 API: http://localhost:${PORT:-8080}"
echo "📊 Health: http://localhost:${PORT:-8080}/health"
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Run the application
./bin/iac-ai-agent -config configs/app.yaml
