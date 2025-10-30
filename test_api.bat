@echo off
echo Starting VPN Client API Server...
echo =================================

echo To test the API, you can use the following curl commands:
echo.
echo List servers:
echo curl http://localhost:8080/api/servers
echo.
echo Get connection status:
echo curl http://localhost:8080/api/status
echo.
echo Get configuration:
echo curl http://localhost:8080/api/config
echo.

echo Starting server on port 8080...
echo Press Ctrl+C to stop the server
echo.

cd src
go run main.go --api