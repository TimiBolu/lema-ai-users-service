#!/bin/bash

echo "Installing Users-Posts Backend..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}')
if [[ "${GO_VERSION}" < "go1.23" ]]; then
    echo "Error: Go version 1.23.0 or higher is required"
    exit 1
fi

# Install dependencies
echo "Installing dependencies..."
go mod download

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cp .env.example .env || echo "Warning: Could not create .env file"
fi

# Build the application
echo "Building application..."
go build -o server

echo "Installation complete!"
echo "To start the server, run: ./server"
echo "Visit /docs route after starting the server for API documentation"
