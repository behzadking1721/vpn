#!/bin/bash

echo "Building VPN Client for Linux..."

# Set version from first argument or default
VERSION=${1:-0.1.0}

# Create output directory
mkdir -p dist/linux

# Build for Linux amd64
echo "Building for Linux AMD64..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o dist/linux/vpn-client-amd64 ./src

# Build for Linux 386
echo "Building for Linux 386..."
CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -ldflags "-s -w -X main.version=$VERSION" -o dist/linux/vpn-client-386 ./src

# Build for Linux ARM64
echo "Building for Linux ARM64..."
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X main.version=$VERSION" -o dist/linux/vpn-client-arm64 ./src

# Build for Linux ARM
echo "Building for Linux ARM..."
CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -ldflags "-s -w -X main.version=$VERSION" -o dist/linux/vpn-client-arm ./src

echo "Linux builds completed."