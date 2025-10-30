@echo off
setlocal

echo Building VPN Client for Windows...

:: Set version from environment variable or default
set VERSION=%1
if "%VERSION%"=="" set VERSION=0.1.0

:: Create output directory
mkdir dist\windows 2>nul

:: Build for Windows amd64
echo Building for Windows AMD64...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w -H windowsgui -X main.version=%VERSION%" -o dist\windows\vpn-client-amd64.exe ./src

:: Build for Windows 386
echo Building for Windows 386...
set GOOS=windows
set GOARCH=386
go build -ldflags "-s -w -H windowsgui -X main.version=%VERSION%" -o dist\windows\vpn-client-386.exe ./src

:: Build for Windows ARM64
echo Building for Windows ARM64...
set GOOS=windows
set GOARCH=arm64
go build -ldflags "-s -w -H windowsgui -X main.version=%VERSION%" -o dist\windows\vpn-client-arm64.exe ./src

echo Windows builds completed.