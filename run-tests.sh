#!/bin/bash

echo "Running VPN Client tests..."
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

# Run unit tests
echo "Running unit tests..."
if ! go test ./src/... -v; then
    echo "Error: Tests failed"
    exit 1
fi

echo
echo "All tests passed!"
echo

# Run the application with version flag
echo "Testing application execution..."
if ! ./vpn-client --version; then
    echo "Error: Failed to run application"
    exit 1
fi

echo
echo "Application test successful!"