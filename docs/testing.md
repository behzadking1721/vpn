# VPN Client Testing Guide

This document provides comprehensive information about how to run tests for the VPN Client application. Testing is a critical part of maintaining code quality and ensuring the application functions correctly.

## Table of Contents

1. [Overview](#overview)
2. [Test Structure](#test-structure)
3. [Running Tests](#running-tests)
4. [Test Types](#test-types)
5. [Writing Tests](#writing-tests)
6. [Performance Testing](#performance-testing)
7. [Security Testing](#security-testing)
8. [Cross-Platform Testing](#cross-platform-testing)
9. [Continuous Integration](#continuous-integration)

## Overview

The VPN Client uses Go's built-in testing framework for unit and integration tests. Tests are organized by packages and can be run individually or all together.

## Test Structure

The tests are organized in the following structure:

```
internal/
  managers/
    *_test.go        # Unit and integration tests for managers package
tests/
  test_*.go          # Integration tests and test utilities
```

### Test Files

- `*_test.go`: Files containing unit tests, following Go's convention
- `connection_e2e_test.go`: End-to-end tests for VPN connection functionality
- `server_manager_test.go`: Tests for server management functionality
- `subscription_parser_test.go`: Tests for subscription parsing functionality
- `ui_test.go`: Tests for UI elements
- `simple_test.go`: Simple test examples
- `manager_test.go`: Additional tests for connection managers

## Running Tests

### Prerequisites

Before running tests, ensure you have:

1. Go installed (version 1.16 or higher)
2. All dependencies installed (run `go mod tidy`)

### Running All Tests

To run all tests in the project:

```bash
# Navigate to the project root
cd /path/to/vpn-client

# Run all tests with verbose output
go test ./internal/managers -v

# Run all tests with coverage
go test ./internal/managers -cover
```

### Running Specific Test Packages

You can run tests for specific packages:

```bash
# Run only manager tests
go test ./internal/managers -v

# Run tests with coverage profile
go test ./internal/managers -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

### Running Specific Tests

You can run specific tests by name using the `-run` flag:

```bash
# Run only VPN connection end-to-end tests
go test ./internal/managers -run TestVPNConnectionEndToEnd -v

# Run only server CRUD tests
go test ./internal/managers -run TestServerCRUD -v

# Run only subscription parser tests
go test ./internal/managers -run TestSubscriptionParser -v

# Run only UI tests
go test ./internal/managers -run TestUI -v
```

### Running Tests with Coverage

To analyze test coverage:

```bash
# Run tests with coverage
go test ./internal/managers -coverprofile=coverage.out

# View coverage in HTML format
go tool cover -html=coverage.out

# View coverage as a function
go tool cover -func=coverage.out
```

### Using Makefile (if available)

If you have `make` installed, you can use the provided Makefile commands:

```bash
# Run all tests
make test

# Run manager tests
make test-managers

# Run protocol tests
make test-protocols
```

Note: On Windows systems, you might need to use `mingw32-make` or run the Go commands directly.

## Test Types

### Unit Tests

Unit tests verify the functionality of individual functions and methods. These tests are fast and don't require external dependencies.

Examples:
- Testing connection status transitions
- Testing server CRUD operations
- Testing subscription parsing

### Integration Tests

Integration tests verify that multiple components work together correctly.

Examples:
- End-to-end VPN connection tests
- Server management with database operations

### UI Tests

UI tests verify that the HTML files contain the expected elements.

Examples:
- Testing that UI files contain required elements
- Testing navigation between UI pages

## Writing Tests

### Test Structure

Tests should follow Go's testing conventions:

```go
func TestFunctionName(t *testing.T) {
    // Test setup
    
    // Execute function under test
    
    // Verify results
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Test Naming

- Use `Test` prefix followed by the function name being tested
- For table-driven tests, use descriptive subtest names
- Keep test names descriptive but concise

### Test Organization

1. **Setup**: Initialize required objects and data
2. **Execution**: Call the function being tested
3. **Verification**: Check that results match expectations
4. **Cleanup**: Clean up any resources if necessary

### Example Test

```go
func TestAddServer(t *testing.T) {
    // Setup
    store := setupTestDB(t)
    sm := NewServerManager(store)
    
    server := &core.Server{
        ID:       "test-server",
        Name:     "Test Server",
        Host:     "example.com",
        Port:     8080,
        Protocol: "vmess",
        Config: map[string]interface{}{
            "user_id": "test_user_id",
        },
        Enabled: true,
    }
    
    // Execution
    err := sm.AddServer(server)
    
    // Verification
    if err != nil {
        t.Fatalf("Failed to add server: %v", err)
    }
    
    retrievedServer, err := sm.GetServer(server.ID)
    if err != nil {
        t.Fatalf("Failed to get server: %v", err)
    }
    
    if retrievedServer.ID != server.ID {
        t.Errorf("Expected server ID %s, got %s", server.ID, retrievedServer.ID)
    }
}
```

### Mocking Dependencies

For tests that require external dependencies, use mocking or create test doubles:

```go
// setupTestDB creates a temporary database for testing
func setupTestDB(t *testing.T) database.Store {
    // Create a temporary directory for test data
    tempDir := t.TempDir()
    
    // Create a JSON store for testing
    store, err := database.NewJSONStore(tempDir)
    if err != nil {
        t.Fatalf("Failed to create test store: %v", err)
    }
    
    return store
}
```

## Performance Testing

Performance testing ensures that the VPN client can handle expected loads and performs efficiently under various conditions.

### Load Testing

Load testing verifies that the application can handle expected user loads:

1. **Connection Load Testing**:
   - Test concurrent connections to multiple servers
   - Measure connection establishment time under load
   - Verify resource usage (CPU, memory) during high connection loads

2. **Server Management Load Testing**:
   - Test managing large numbers of servers (1000+ servers)
   - Measure response times for server listing operations
   - Verify database performance with large datasets

Example load test script (using a tool like [hey](https://github.com/rakyll/hey)):

```bash
# Test API server with 100 requests, 10 concurrent
hey -n 100 -c 10 http://localhost:8080/api/servers
```

### Stress Testing

Stress testing pushes the application beyond normal operational capacity:

1. **Extreme Connection Testing**:
   - Test rapid connect/disconnect cycles
   - Test with maximum number of servers
   - Test with malformed or invalid server configurations

2. **Resource Exhaustion Testing**:
   - Test behavior when system resources are limited
   - Test with low memory conditions
   - Test with high CPU usage scenarios

### Performance Metrics

Key performance metrics to monitor:

- Connection establishment time
- Data transfer rates
- Memory usage
- CPU usage
- API response times
- Database query times

## Security Testing

Security testing ensures that the VPN client properly protects user data and communications.

### Penetration Testing

Penetration testing simulates attacks to identify vulnerabilities:

1. **API Security Testing**:
   - Test for unauthorized access to API endpoints
   - Test input validation for all API parameters
   - Test for injection vulnerabilities (SQL, command, etc.)

2. **Data Security Testing**:
   - Verify that sensitive data is properly encrypted
   - Test access controls for user data
   - Verify secure storage of credentials

3. **Network Security Testing**:
   - Test for information leakage through network traffic
   - Verify that VPN tunnels are properly encrypted
   - Test for man-in-the-middle attack vulnerabilities

### Security Testing Tools

Recommended security testing tools:

1. **OWASP ZAP**: For web application security testing
2. **Nmap**: For network discovery and security auditing
3. **Burp Suite**: For security testing of web applications
4. **Bandit**: For finding common security issues in Python code (if applicable)

### Security Test Examples

1. **API Input Validation**:
   ```bash
   # Test API with malicious input
   curl -X POST http://localhost:8080/api/servers \
     -H "Content-Type: application/json" \
     -d '{"id": "test", "name": "<script>alert(1)</script>"}'
   ```

2. **Authentication Testing**:
   - Test accessing protected endpoints without authentication
   - Test weak password policies (if applicable)
   - Test session management vulnerabilities

## Cross-Platform Testing

Cross-platform testing ensures that the VPN client works correctly on all supported operating systems.

### Supported Platforms

The VPN Client supports the following platforms:

1. **Windows** (Windows 10 and later)
2. **Linux** (Ubuntu 18.04+, CentOS 7+, Debian 9+)
3. **macOS** (macOS 10.14 and later)

### Platform-Specific Testing

Each platform requires specific testing considerations:

#### Windows Testing

1. **Windows Service Testing**:
   - Test running as a Windows service
   - Test automatic startup
   - Test service recovery options

2. **Windows UI Testing**:
   - Test desktop integration
   - Test notification area integration
   - Test Windows Firewall compatibility

#### Linux Testing

1. **Systemd Service Testing**:
   - Test systemd service installation
   - Test automatic startup
   - Test logging integration

2. **Distribution Testing**:
   - Test on Ubuntu, CentOS, and Debian
   - Test package installation (DEB/RPM)
   - Test desktop environment integration (GNOME, KDE)

#### macOS Testing

1. **macOS Application Testing**:
   - Test application bundle functionality
   - Test Gatekeeper compatibility
   - Test system preferences integration

2. **Permission Testing**:
   - Test network extension permissions
   - Test keychain integration
   - Test notarization requirements

### Cross-Platform Test Scenarios

1. **Installation Testing**:
   - Test installation on each platform
   - Test upgrade scenarios
   - Test uninstallation and cleanup

2. **Configuration Testing**:
   - Test configuration file compatibility
   - Test environment variable handling
   - Test command-line argument parsing

3. **Functionality Testing**:
   - Test all core features on each platform
   - Test platform-specific features
   - Test file system integration

4. **Performance Testing**:
   - Test performance consistency across platforms
   - Test resource usage on each platform
   - Test network performance

### Cross-Platform Testing Tools

Recommended tools for cross-platform testing:

1. **Virtual Machines**: VMware, VirtualBox for testing different OS versions
2. **Containerization**: Docker for testing Linux environments
3. **Cloud Testing**: AWS, Azure, or Google Cloud for testing on different platforms
4. **CI/CD Systems**: GitHub Actions, GitLab CI for automated cross-platform testing

## Continuous Integration

The project can be configured to run tests automatically in a CI environment.

### GitHub Actions

Example GitHub Actions workflow for running tests:

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.16
    - name: Test
      run: go test ./internal/managers -v
```

### Test Reporting

Tests output results in standard Go test format. For CI environments, you can generate JUnit-style reports using tools like `go-junit-report`:

```bash
# Install go-junit-report
go install github.com/jstemmer/go-junit-report@latest

# Run tests and generate JUnit report
go test ./internal/managers -v | go-junit-report > report.xml
```

## Troubleshooting

### Common Issues

1. **Missing Dependencies**:
   ```bash
   # Run to install dependencies
   go mod tidy
   ```

2. **Permission Issues**:
   Ensure you have write permissions to the project directory, especially for tests that create temporary files.

3. **Test Failures**:
   - Check the error message for details
   - Run the specific failing test with `-v` flag for verbose output
   - Ensure all required services are running

4. **Coverage Issues**:
   If coverage is lower than expected, check if all code paths are covered by tests.

### Debugging Tests

To debug tests, you can:

1. Add print statements:
   ```go
   t.Logf("Debug info: %v", variable)
   ```

2. Run a single test:
   ```bash
   go test ./internal/managers -run TestSpecificTest -v
   ```

3. Use Go's debugger (Delve):
   ```bash
   dlv test ./internal/managers -- -test.run TestSpecificTest
   ```

## Best Practices

1. **Keep tests fast**: Avoid unnecessary setup or long-running operations
2. **Make tests independent**: Each test should be able to run independently
3. **Use table-driven tests**: For testing multiple cases of the same function
4. **Name tests clearly**: Test names should clearly indicate what is being tested
5. **Test edge cases**: Include tests for error conditions and boundary values
6. **Maintain test coverage**: Aim for high test coverage, especially for critical functionality
7. **Keep tests maintainable**: Avoid duplication and keep tests simple to understand

## Conclusion

Testing is an integral part of the development process for the VPN Client. Regularly running tests ensures that new changes don't break existing functionality and helps maintain code quality. Always run tests before submitting changes and add new tests when adding functionality.