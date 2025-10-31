#!/bin/bash

# Cross-platform build script for VPN Client

echo "Building VPN Client for all platforms..."

# Create output directory
mkdir -p dist

# Build for Windows
echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -o dist/vpn-client-windows-amd64.exe ./src

# Build for Linux
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -o dist/vpn-client-linux-amd64 ./src

# Build for macOS
echo "Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -o dist/vpn-client-darwin-amd64 ./src

# Build for macOS ARM64 (Apple Silicon)
echo "Building for macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -o dist/vpn-client-darwin-arm64 ./src

echo "Build process completed!"
echo "Binaries are located in the dist/ directory."