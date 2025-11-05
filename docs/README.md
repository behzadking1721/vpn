# VPN Client Project

## Overview

This project is a comprehensive VPN client solution with multiple interfaces and components:

1. **Backend Service** - Written in Go, handles VPN connections, server management, and core functionality
2. **Desktop UI** - Web-based user interface for desktop usage
3. **Mobile App** - React Native application for mobile devices
4. **CLI Tool** - Command-line interface for advanced users and automation

## Components

### Backend Service
The backend service is written in Go and provides a REST API for all client interfaces. It handles:
- VPN connection management
- Server and subscription management
- Statistics tracking
- Notification system
- Logging

### Desktop UI
The desktop UI is a web-based interface that communicates with the backend service through REST API calls. Features include:
- Server listing and management
- Connection status and controls
- Statistics visualization
- Settings management

### Mobile App
The mobile app is built with React Native and provides a native mobile experience:
- All essential VPN client features
- Native mobile UI/UX
- Push notifications

### CLI Tool (vpnctl)
The command-line interface provides powerful automation capabilities:
- Full control over VPN connections
- Server management
- Statistics and monitoring
- Configuration management
- Scripting support

## Installation

### Prerequisites
- Go 1.21 or higher
- Node.js and npm (for UI development)
- Mobile development environment (for mobile app)

### Building

1. **Backend Service**:
   ```bash
   cd cmd/vpn-client
   go build -o vpn-client
   ```

2. **CLI Tool**:
   ```bash
   cd cmd/vpnctl
   go build -o vpnctl
   ```

3. **Desktop UI**:
   ```bash
   cd ui/desktop
   # UI is served by the backend service
   ```

4. **Mobile App**:
   ```bash
   cd mobile
   npm install
   npx react-native run-android # or run-ios
   ```

## Usage

### Running the Service
```bash
./vpn-client
```
The service will start and listen on `http://localhost:8080`.

### Using the CLI Tool
```bash
# Connect to the best server
./vpnctl connect --best

# List all servers
./vpnctl server list

# View connection status
./vpnctl status

# View logs
./vpnctl logs
```

For detailed CLI usage, see [CLI Documentation](cli.md).

### Using the Desktop UI
Open your browser and navigate to `http://localhost:8080`.

### Using the Mobile App
Launch the mobile app on your device or emulator.

## API Documentation

See [API Documentation](api.md) for details on the REST API.

## Development

### Project Structure
```
vpnclient/
├── cmd/
│   ├── vpn-client/        # Main backend service
│   └── vpnctl/            # CLI tool
├── internal/
│   ├── api/               # REST API handlers
│   ├── database/          # Database abstractions
│   ├── logging/           # Logging system
│   ├── managers/          # Business logic managers
│   ├── notifications/     # Notification system
│   ├── stats/             # Statistics system
│   └── core/              # Core data structures
├── ui/
│   └── desktop/           # Desktop web UI
├── mobile/                # Mobile application
├── docs/                  # Documentation
└── data/                  # Data storage
```

### Testing

Run tests with:
```bash
go test ./...
```

See [Testing Documentation](testing.md) for more details.

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.