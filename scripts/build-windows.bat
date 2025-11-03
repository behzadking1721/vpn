@echo off
setlocal enabledelayedexpansion

echo Building VPN Client for Windows...

REM Change to project root directory
cd /d %~dp0..
echo Current directory: %CD%

REM Get version from git or use default
for /f "tokens=*" %%a in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%a
if "%VERSION%"=="" set VERSION=dev

echo Building version: %VERSION%

REM Create output directory
mkdir dist 2>nul
mkdir dist\windows_amd64 2>nul
mkdir dist\windows_386 2>nul
mkdir dist\windows_arm64 2>nul

REM Build for Windows AMD64
echo Building for Windows AMD64...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w -X main.version=%VERSION%" -o dist\windows_amd64\vpn-client.exe ./cmd/vpn-client

REM Build for Windows 386
echo Building for Windows 386...
set GOOS=windows
set GOARCH=386
go build -ldflags "-s -w -X main.version=%VERSION%" -o dist\windows_386\vpn-client.exe ./cmd/vpn-client

REM Build for Windows ARM64
echo Building for Windows ARM64...
set GOOS=windows
set GOARCH=arm64
go build -ldflags "-s -w -X main.version=%VERSION%" -o dist\windows_arm64\vpn-client.exe ./cmd/vpn-client

echo Windows builds completed.
echo Binaries are located in the dist\ directory.