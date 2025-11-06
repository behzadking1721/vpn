#!/bin/bash

echo "Running VPN Client Tests"
echo "======================="

echo
echo "1. Running unit tests..."
go test ./internal/managers/... -v
if [ $? -ne 0 ]; then
    echo "Unit tests failed!"
    exit 1
fi

echo
echo "2. Running protocol integration test..."
go test ./internal/... -v
if [ $? -ne 0 ]; then
    echo "Protocol integration test failed!"
    exit 1
fi

echo
echo "All tests passed successfully!"