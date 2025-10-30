@echo off
echo Running VPN Client Tests
echo =======================

echo.
echo 1. Running unit tests...
cd src
go test -v ./protocols
if %errorlevel% neq 0 (
    echo Unit tests failed!
    exit /b %errorlevel%
)

echo.
echo 2. Running protocol integration test...
cd ..
go run test_protocols.go
if %errorlevel% neq 0 (
    echo Protocol integration test failed!
    exit /b %errorlevel%
)

echo.
echo All tests passed successfully!