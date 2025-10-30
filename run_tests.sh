#!/bin/bash

echo "Running VPN Client Tests"
echo "======================="

echo
echo "1. Running unit tests..."
cd src
go test -v ./protocols
if [ $? -ne 0 ]; then
    echo "Unit tests failed!"
    exit 1
fi

echo
echo "2. Running protocol integration test..."
cd ..
go run test_protocols.go
if [ $? -ne 0 ]; then
    echo "Protocol integration test failed!"
    exit 1
fi

echo
echo "All tests passed successfully!"