@echo off
REM Build script for VPN client application on Windows

echo Building VPN Client Application...

REM Create bin directory if it doesn't exist
if not exist "bin" mkdir bin

REM Build for Windows
echo Building for Windows...
go build -o bin\vpn-client.exe src\main.go

echo Build complete!
echo Binary located in .\bin directory

REM List built binary
dir bin\