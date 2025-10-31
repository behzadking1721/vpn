#!/bin/bash

echo "Setting up VPN Client development environment..."
echo

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Error: Go is not installed or not in PATH"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo "Go is installed:"
go version
echo

# Initialize or tidy Go modules
echo "Initializing/Tidying Go modules..."
if ! go mod tidy; then
    echo "Error: Failed to initialize/tidy Go modules"
    exit 1
fi

echo "Go modules initialized successfully!"
echo

# Build the application
echo "Building VPN Client..."
if ! go build -o vpn-client ./src; then
    echo "Error: Failed to build VPN Client"
    exit 1
fi

echo "VPN Client built successfully!"
echo

echo "Development environment setup complete!"
echo "You can now run the application with: ./vpn-client"