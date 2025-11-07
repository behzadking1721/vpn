#!/usr/bin/env pwsh

# Cross-platform build script for VPN Client

param(
    [string]$Version = "dev"
)

Write-Host "Building VPN Client for all platforms..." 

# Create output directory
New-Item -ItemType Directory -Force -Path dist | Out-Null

# Build for Windows
Write-Host "Building for Windows AMD64..."
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o dist/vpn-client-windows-amd64.exe ./cmd/vpn-client

# Build for Linux
Write-Host "Building for Linux AMD64..."
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o dist/vpn-client-linux-amd64 ./cmd/vpn-client

# Build for macOS
Write-Host "Building for macOS AMD64..."
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o dist/vpn-client-darwin-amd64 ./cmd/vpn-client

# Build for macOS ARM64 (Apple Silicon)
Write-Host "Building for macOS ARM64..."
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o dist/vpn-client-darwin-arm64 ./cmd/vpn-client

Write-Host "Build process completed!" 
Write-Host "Binaries are located in the dist/ directory."