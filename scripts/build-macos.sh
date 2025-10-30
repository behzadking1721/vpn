#!/bin/bash

echo "Building VPN Client for macOS..."

# Set version from first argument or default
VERSION=${1:-0.1.0}

# Create output directory
mkdir -p dist/macos

# Build for macOS amd64
echo "Building for macOS AMD64..."
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o dist/macos/vpn-client-amd64 ./src

# Build for macOS ARM64 (Apple Silicon)
echo "Building for macOS ARM64..."
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -X main.version=$VERSION" -o dist/macos/vpn-client-arm64 ./src

echo "macOS builds completed."