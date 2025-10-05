#!/bin/bash

# Setup script for IaC AI Agent

set -e

echo "🚀 Setting up IaC AI Agent..."

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✅ Go version: $(go version)"

# Check Python installation (for Checkov)
if ! command -v python3 &> /dev/null; then
    echo "⚠️  Python3 is not installed. Checkov will not be available."
else
    echo "✅ Python version: $(python3 --version)"
    
    # Install Checkov
    echo "📦 Installing Checkov..."
    pip3 install --upgrade checkov || echo "⚠️  Failed to install Checkov"
fi

# Check Terraform installation
if ! command -v terraform &> /dev/null; then
    echo "⚠️  Terraform is not installed. Some features may not work."
else
    echo "✅ Terraform version: $(terraform version | head -n1)"
fi

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod download
go mod verify

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cat > .env << EOF
# LLM Configuration
LLM_PROVIDER=openai
LLM_API_KEY=your-api-key-here
LLM_MODEL=gpt-4

# GitHub Configuration
GITHUB_TOKEN=your-github-token-here
GITHUB_WEBHOOK_SECRET=your-webhook-secret-here

# Analysis Configuration
CHECKOV_ENABLED=true
IAM_ANALYSIS_ENABLED=true
COST_OPTIMIZATION_ENABLED=true

# Server Configuration
PORT=8080
LOG_LEVEL=info
EOF
    echo "⚠️  Please edit .env file with your credentials"
fi

# Build the application
echo "🔨 Building application..."
make build

echo ""
echo "✅ Setup complete!"
echo ""
echo "📋 Next steps:"
echo "   1. Edit .env file with your credentials"
echo "   2. Run 'make run' to start the application"
echo "   3. Visit http://localhost:8080 to verify it's running"
echo ""
echo "📚 For more information, see docs/README.md"
