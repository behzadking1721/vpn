@echo off
echo Running VPN Client tests...
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

REM Run unit tests
echo Running unit tests...
go test ./src/... -v
if %errorlevel% neq 0 (
    echo Error: Tests failed
    pause
    exit /b 1
)

echo.
echo All tests passed!
echo.

REM Run the application with version flag
echo Testing application execution...
vpn-client.exe --version
if %errorlevel% neq 0 (
    echo Error: Failed to run application
    pause
    exit /b 1
)

echo.
echo Application test successful!
echo.
pause