@echo off
echo Building VPN Client for Windows...

REM Create output directory
mkdir dist 2>nul

REM Build for Windows AMD64
echo Building for Windows AMD64...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -o dist\vpn-client-windows-amd64.exe ./src

REM Build for Windows 386
echo Building for Windows 386...
set GOOS=windows
set GOARCH=386
go build -o dist\vpn-client-windows-386.exe ./src

REM Build for Windows ARM64
echo Building for Windows ARM64...
set GOOS=windows
set GOARCH=arm64
go build -o dist\vpn-client-windows-arm64.exe ./src

echo Windows builds completed.
echo Binaries are located in the dist\ directory.