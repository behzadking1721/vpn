#!/usr/bin/env pwsh

# Android build script for VPN Client

param(
    [string]$Version = "dev"
)

Write-Host "Building Android packages for VPN Client v$Version" -ForegroundColor Blue

# Check if Flutter project exists
if (-not (Test-Path "ui\mobile\flutter_vpn")) {
    Write-Host "Error: Flutter project not found at ui\mobile\flutter_vpn" -ForegroundColor Red
    exit 1
}

Write-Host "Changing directory to Flutter project..." -ForegroundColor Yellow
Set-Location "ui\mobile\flutter_vpn"

# Build APKs for different architectures
Write-Host "Building APKs for different architectures..." -ForegroundColor Yellow
flutter build apk --release --split-per-abi

# Build App Bundle for Google Play
Write-Host "Building App Bundle for Google Play..." -ForegroundColor Yellow
flutter build appbundle --release

Write-Host "Android build process completed!" -ForegroundColor Green
Write-Host "Output files are located in build/app/outputs/" -ForegroundColor Blue

# Return to original directory
Set-Location "..\..\.."