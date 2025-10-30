# Protocol Testing Guide

This document explains how to test protocol integration in the VPN client application.

## Overview

Testing protocol integration is crucial to ensure that all supported protocols work correctly with the application. This guide covers both unit testing and integration testing approaches.

## Unit Testing

Unit tests verify that each protocol handler works correctly in isolation.

### Running Unit Tests

To run all protocol tests:

```bash
cd src
go test ./protocols
```

To run with verbose output:

```bash
cd src
go test -v ./protocols
```

### Writing Protocol Tests

Each protocol should have tests covering:

1. Handler creation
2. Connection establishment
3. Data usage tracking
4. Connection details retrieval
5. Disconnection

Example test structure:

```go
func TestProtocolHandler(t *testing.T) {
    handler := NewProtocolHandler()
    
    // Test initial state
    if handler.IsConnected() {
        t.Error("Expected handler to be disconnected initially")
    }
    
    // Test connection
    server := core.Server{
        // Configuration
    }
    
    err := handler.Connect(server)
    if err != nil {
        t.Errorf("Failed to connect: %v", err)
    }
    
    // Test functionality
    // ...
    
    // Test disconnection
    err = handler.Disconnect()
    if err != nil {
        t.Errorf("Failed to disconnect: %v", err)
    }
}
```

## Integration Testing

Integration tests verify that protocol handlers work correctly with other components.

### Test Script

The project includes a test script [test_protocols.go](file:///c:/Users/behza/OneDrive/Documents/vpn/test_protocols.go) that tests all protocol handlers:

```bash
go run test_protocols.go
```

This script:
1. Creates a protocol factory
2. Tests each protocol handler
3. Attempts to connect to a sample server
4. Checks data usage and connection details
5. Disconnects from the server

### Manual Testing

For manual testing with real servers:

1. Add a real server configuration
2. Run the CLI: `go run main.go --cli`
3. Connect to the server
4. Monitor connection status
5. Check data usage
6. Disconnect

## Continuous Integration

The project should include CI configuration that:

1. Runs unit tests on each commit
2. Runs integration tests nightly
3. Tests on multiple platforms (Windows, Linux, macOS)
4. Checks for security vulnerabilities

Example GitHub Actions workflow:

```yaml
name: Protocol Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.21
    - name: Run tests
      run: |
        cd src
        go test -v ./protocols
```

## Test Coverage

Aim for high test coverage:

1. **Handler Creation**: 100%
2. **Connection Logic**: 90%+
3. **Error Handling**: 90%+
4. **Data Tracking**: 80%+
5. **Edge Cases**: 70%+

To check test coverage:

```bash
cd src
go test -cover ./protocols
```

For detailed coverage report:

```bash
cd src
go test -coverprofile=coverage.out ./protocols
go tool cover -html=coverage.out
```

## Mocking External Dependencies

When testing protocol handlers that use external libraries:

1. Create interfaces for external dependencies
2. Use mocking frameworks like [gomock](https://github.com/golang/mock)
3. Mock network calls and library functions

Example:

```go
//go:generate mockgen -source=shadowsocks_lib.go -destination=mocks/shadowsocks_mock.go

type ShadowsocksClient interface {
    Connect() error
    Disconnect() error
    GetBytesSent() int64
    GetBytesReceived() int64
}
```

## Performance Testing

Test protocol performance:

1. Connection establishment time
2. Data transfer rates
3. Memory usage
4. CPU usage

Example benchmark test:

```go
func BenchmarkProtocolConnect(b *testing.B) {
    handler := NewProtocolHandler()
    server := createSampleServer()
    
    for i := 0; i < b.N; i++ {
        handler.Connect(server)
        handler.Disconnect()
    }
}
```

Run benchmarks:

```bash
cd src
go test -bench=. ./protocols
```

## Security Testing

Security aspects to test:

1. Input validation
2. Buffer overflow protection
3. Memory safety
4. Encryption correctness

## Troubleshooting

Common testing issues and solutions:

### Dependency Issues
```
cannot find module providing package
```
Solution: Run `go mod tidy` to update dependencies.

### Network Issues
Tests fail due to network connectivity.
Solution: Use mocking for network-dependent tests.

### Race Conditions
Tests fail intermittently.
Solution: Run with race detector: `go test -race ./protocols`

## Best Practices

1. **Isolate Tests**: Each test should be independent
2. **Use Table-Driven Tests**: For testing multiple scenarios
3. **Clean Up Resources**: Always disconnect and clean up
4. **Test Edge Cases**: Empty values, invalid configurations
5. **Mock External Services**: Don't rely on real servers for unit tests
6. **Measure Performance**: Regularly benchmark protocol performance
7. **Automate Testing**: Use CI/CD to run tests automatically

## Next Steps

1. Implement comprehensive unit tests for all protocol handlers
2. Set up continuous integration
3. Add performance benchmarks
4. Implement security testing
5. Create test documentation for contributors