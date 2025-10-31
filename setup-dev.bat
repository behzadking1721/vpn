@echo off
echo Setting up VPN Client development environment...
echo.

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

echo Go is installed:
go version
echo.

REM Initialize or tidy Go modules
echo Initializing/Tidying Go modules...
go mod tidy
if %errorlevel% neq 0 (
    echo Error: Failed to initialize/tidy Go modules
    pause
    exit /b 1
)

echo Go modules initialized successfully!
echo.

REM Build the application
echo Building VPN Client...
go build -o vpn-client.exe ./src
if %errorlevel% neq 0 (
    echo Error: Failed to build VPN Client
    pause
    exit /b 1
)

echo VPN Client built successfully!
echo.

echo Development environment setup complete!
echo You can now run the application with: vpn-client.exe
echo.
pause