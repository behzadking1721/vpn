#!/usr/bin/env pwsh

# Packaging script for VPN client - creates installers for various platforms

param(
    [string]$Version = "dev"
)

# If version is still default, try to get from git
if ($Version -eq "dev") {
    $gitVersion = git describe --tags --always --dirty 2>$null
    if ($gitVersion) {
        $Version = $gitVersion
    }
}

# Base directories
$DIST_DIR = "dist"
$PACKAGE_DIR = "packages"
$RELEASE_DIR = "release"

# Clean previous packages
Write-Host "Cleaning previous packages..." -ForegroundColor Blue
if (Test-Path $PACKAGE_DIR) { Remove-Item $PACKAGE_DIR -Recurse -Force }
if (Test-Path $RELEASE_DIR) { Remove-Item $RELEASE_DIR -Recurse -Force }
New-Item -ItemType Directory -Force -Path $PACKAGE_DIR | Out-Null
New-Item -ItemType Directory -Force -Path $RELEASE_DIR | Out-Null

Write-Host "Packaging VPN Client v$Version" -ForegroundColor Blue

# Check if builds exist
if (-not (Test-Path $DIST_DIR)) {
    Write-Host "Error: Build directory '$DIST_DIR' does not exist. Run build script first." -ForegroundColor Red
    exit 1
}

# Create ZIP archives for all platforms
Write-Host "Creating ZIP archives..." -ForegroundColor Yellow

Get-ChildItem $DIST_DIR | Where-Object { $_.PSIsContainer } | ForEach-Object {
    $platform = $_.Name
    $zipName = "vpn-client-$platform-$Version.zip"
    
    Write-Host "Creating archive for $platform..." -ForegroundColor Blue
    Compress-Archive -Path "$DIST_DIR\$platform\*" -DestinationPath "$PACKAGE_DIR\$zipName" -Force
    Write-Host "OK Created $zipName" -ForegroundColor Green
}

# Package Android apps if they exist
if (Test-Path "ui\mobile\flutter_vpn\build\app\outputs") {
    Write-Host "Packaging Android apps..." -ForegroundColor Yellow
    
    # Create android directory in packages
    New-Item -ItemType Directory -Force -Path "$PACKAGE_DIR\android" | Out-Null
    
    $androidOutputs = "ui\mobile\flutter_vpn\build\app\outputs"
    
    # Copy APKs
    if (Test-Path "$androidOutputs\flutter-apk") {
        Get-ChildItem "$androidOutputs\flutter-apk\app-*-release.apk" | ForEach-Object {
            $newName = "vpn-client-$Version-$(($_.Name -replace 'app-', '') -replace '-release.apk', '.apk')"
            Copy-Item $_.FullName "$PACKAGE_DIR\android\$newName"
            Write-Host "OK Created $newName" -ForegroundColor Green
        }
    }
    
    # Copy App Bundle
    if (Test-Path "$androidOutputs\bundle\release\app-release.aab") {
        $aabName = "vpn-client-$Version-appbundle.aab"
        Copy-Item "$androidOutputs\bundle\release\app-release.aab" "$PACKAGE_DIR\android\$aabName"
        Write-Host "OK Created $aabName" -ForegroundColor Green
    }
    
    Write-Host "Android packages created successfully!" -ForegroundColor Green
}

Write-Host "All packaging completed successfully!" -ForegroundColor Green
Write-Host "Packages located in: $PACKAGE_DIR" -ForegroundColor Blue