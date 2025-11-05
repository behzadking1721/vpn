# Logging System Documentation

## Overview

The logging system provides comprehensive logging capabilities for the VPN client application. It enables debugging, monitoring, and troubleshooting by recording important events, errors, and system information.

## Architecture

The logging system consists of several components:

1. **Logger Package** - Core logging functionality in Go
2. **API Endpoints** - RESTful interface for log management
3. **Desktop UI Integration** - Log viewing and management in the desktop application
4. **Log File Storage** - Persistent storage of log messages

## Logger Package

The logger package is implemented in the [internal/logging](file:///c%3A/Users/behza/OneDrive/Documents/vpn/internal/logging) directory and provides the following key features:

### Log Levels

The system supports five log levels:

1. **DEBUG** - Detailed information for diagnosing problems
2. **INFO** - General information about system operation
3. **WARNING** - Warning messages about potential issues
4. **ERROR** - Error messages about serious problems
5. **FATAL** - Critical errors that cause the application to exit

### Logger Configuration

The logger can be configured with the following options:

- **Level** - Minimum log level to record
- **Output** - Destination for log messages (stdout, stderr, or file)
- **Timestamp** - Whether to include timestamps in log messages

### Features

- Thread-safe logging
- File rotation support
- Caller information (file and line number)
- Structured logging with key-value pairs
- Specialized logging functions for different components

## API Endpoints

The logging system exposes the following RESTful API endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/logs` | GET | Retrieve application logs |
| `/api/logs/clear` | POST | Clear all application logs |
| `/api/logs/stats` | GET | Get log statistics |

### GET /api/logs

Retrieve application logs with optional filtering.

**Query Parameters:**
- `limit` - Maximum number of log entries to return
- `level` - Filter logs by level (debug, info, warning, error, fatal)

**Response:**
```json
{
  "logs": [
    "[2023-01-01 12:00:00] [INFO] [server_manager.go:123] Server added: Example Server (example.com:443)",
    "[2023-01-01 12:01:00] [ERROR] [connection.go:456] Connection failed: timeout"
  ],
  "count": 2,
  "file": "data/vpn-client.log"
}
```

### POST /api/logs/clear

Clear all application logs.

**Response:**
```json
{
  "status": "success",
  "message": "Logs cleared successfully"
}
```

### GET /api/logs/stats

Get statistics about the log file.

**Response:**
```json
{
  "file": {
    "path": "data/vpn-client.log",
    "size": 10240,
    "modified": "2023-01-01T12:00:00Z"
  },
  "stats": {
    "DEBUG": 10,
    "INFO": 50,
    "WARNING": 5,
    "ERROR": 2,
    "FATAL": 0
  },
  "total": 67
}
```

## Desktop UI Integration

The desktop UI provides a user-friendly interface for viewing and managing logs:

### Features

- Log level selection
- Log viewing in a modal dialog
- Log clearing with confirmation
- Real-time log display

### Implementation

The desktop UI uses JavaScript to communicate with the logging API endpoints and display logs in a modal dialog with syntax highlighting.

## Log File Storage

Logs are stored in the `data/vpn-client.log` file with the following format:

```
[2023-01-01 12:00:00] [LEVEL] [source_file.go:line_number] Log message
```

### Log Rotation

The system implements basic log rotation by clearing logs when requested through the API. In production environments, external log rotation tools should be used.

## Integration Points

### Connection Manager

The connection manager logs events related to VPN connections:

- Connection attempts
- Successful connections
- Connection failures
- Disconnection events

### Server Manager

The server manager logs events related to server management:

- Server additions
- Server updates
- Server deletions
- Server enable/disable operations
- Ping test results

### Subscription Manager

The subscription manager logs events related to subscriptions:

- Subscription additions
- Subscription updates
- Subscription deletions
- Parsing errors

## Usage Examples

### Backend Logging

```go
// Create logger
logger, err := logging.NewLogger(logging.Config{
    Level:     logging.INFO,
    Output:    "data/vpn-client.log",
    Timestamp: true,
})
if err != nil {
    log.Fatalf("Failed to initialize logger: %v", err)
}
defer logger.Close()

// Log different types of messages
logger.Info("Application started")
logger.Debug("Processing server list")
logger.Warning("Server ping is high: %d ms", ping)
logger.Error("Failed to connect to server: %v", err)
```

### Specialized Logging

```go
// Log connection events
logger.LogConnectionEvent("Connected", serverID, map[string]interface{}{
    "start_time": time.Now().Format(time.RFC3339),
})

// Log server errors
logger.LogServerError(serverID, err, "Failed to update configuration")

// Log subscription events
logger.LogSubscriptionEvent("Updated", subscriptionID, map[string]interface{}{
    "server_count": len(servers),
})
```

### Frontend Log Viewing

```javascript
// Retrieve logs from API
async function viewLogs() {
    const response = await fetch('/api/logs');
    if (response.ok) {
        const data = await response.json();
        showLogModal(data.logs);
    }
}

// Clear logs
async function clearLogs() {
    const response = await fetch('/api/logs/clear', { method: 'POST' });
    if (response.ok) {
        showNotification('Success', 'Logs cleared successfully', 'success');
    }
}
```

## Configuration

### Log Level Configuration

The log level can be configured through the desktop UI or by modifying the logger initialization in the main application.

### Log File Location

Logs are stored in the `data/vpn-client.log` file relative to the application executable.

## Testing

The logging system should be tested for the following scenarios:

1. **Log Message Format** - Verify correct formatting of log messages
2. **Log Levels** - Test filtering by different log levels
3. **File Rotation** - Test log clearing functionality
4. **Error Handling** - Test behavior when log file cannot be written
5. **Performance** - Test logging performance under load

## Best Practices

1. **Use Appropriate Log Levels** - Choose the right level for each message
2. **Include Context** - Add relevant context to log messages
3. **Avoid Sensitive Information** - Never log passwords or other sensitive data
4. **Use Structured Logging** - Use key-value pairs for complex information
5. **Handle Errors** - Always check for and handle logging errors

## Future Improvements

1. **Log Rotation** - Implement automatic log rotation based on size or time
2. **Remote Logging** - Add support for sending logs to remote servers
3. **Log Search** - Add search functionality for large log files
4. **Log Export** - Add support for exporting logs in different formats
5. **Performance Monitoring** - Add performance metrics to logs