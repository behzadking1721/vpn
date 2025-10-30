#!/bin/bash

# Build script for VPN client application

echo "Building VPN Client Application..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build for different platforms
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o bin/vpn-client-windows.exe src/main.go

echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o bin/vpn-client-linux src/main.go

echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o bin/vpn-client-macos src/main.go

echo "Build complete!"
echo "Binaries located in ./bin directory"

# List built binaries
ls -la bin/