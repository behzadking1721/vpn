# VPN Client Development Guide

This document provides information for developers who want to contribute to the VPN Client project.

## Project Structure

```
vpn-client/
├── src/                 # Main source code
│   ├── main.go          # Entry point
│   ├── protocols/       # Protocol implementations
│   ├── managers/        # Connection, server, and log managers
│   └── utils/           # Utility functions
├── config/              # Configuration files
├── logs/                # Log files
├── data/                # Temporary data
├── scripts/             # Build and installation scripts
│   ├── build-all.sh     # Cross-platform build script
│   ├── build-windows.bat# Windows build script
│   └── docs/            # Script documentation
├── docs/                # Project documentation
└── dist/                # Built binaries
```

## Building the Project

### Prerequisites

- Go 1.21 or higher
- Git

### Building for Current Platform

```bash
go build -o vpn-client ./src
```

### Cross-Platform Builds

#### Using Bash Script (Linux/macOS)
```bash
chmod +x scripts/build-all.sh
./scripts/build-all.sh
```

#### Using Batch Script (Windows)
```cmd
scripts\build-windows.bat
```

## Configuration

The application uses a JSON configuration file located at `config/settings.json`. The structure is:

```json
{
  "servers": [
    {
      "name": "Server Name",
      "address": "server.example.com",
      "port": 443,
      "protocol": "vless",
      "uuid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "security": "tls"
    }
  ],
  "log_level": "info",
  "auto_connect": false
}
```

## Adding New Protocols

To add a new protocol:

1. Create a new file in `src/protocols/` (e.g., `src/protocols/newprotocol.go`)
2. Implement the protocol handler interface:
   ```go
   type ProtocolHandler interface {
       Connect(config *ServerConfig) error
       Disconnect() error
       IsConnected() bool
       GetStats() *ConnectionStats
   }
   ```
3. Register the protocol in the protocol manager

## Code Standards

- Follow Go coding standards
- Use meaningful variable and function names
- Add comments for exported functions
- Write unit tests for new functionality

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```