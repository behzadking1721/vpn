# Cross-Platform VPN Client

A multi-platform VPN application with a simple and user-friendly interface similar to Hiddify, supporting various protocols and subscription links.

## Features

- Cross-platform support (Android, iOS, Windows, macOS, Linux)
- Support for multiple protocols:
  - VMess
  - VLESS
  - Trojan
  - Reality
  - Hysteria2
  - TUIC
  - SSH
  - Shadowsocks
- Server management capabilities
- User and connection settings management
- Fastest server auto-selection (LowestPing)
- Data usage display
- Subscription link import with deep linking
- QR code import
- IPv6 support
- Ad-free experience
- Open-source and secure

## Prerequisites

Before you can build and run this project, you need to install the following dependencies:

### Go (Golang)
- **Windows**: Download from [golang.org](https://golang.org/dl/)
- **macOS**: `brew install go` or download from [golang.org](https://golang.org/dl/)
- **Linux**: `sudo apt install golang` (Ubuntu/Debian) or equivalent for your distribution

### Flutter (for mobile UI - future implementation)
- Download from [flutter.dev](https://flutter.dev/docs/get-started/install)

### Node.js and npm (for desktop UI with Electron - future implementation)
- Download from [nodejs.org](https://nodejs.org/)

## Technical Architecture

The application is built using a modular approach with shared core logic and platform-specific UI implementations.

### Core Components

1. **Protocol Handler**: Manages different VPN protocols
2. **Server Manager**: Handles server connections and configurations
3. **Subscription Manager**: Processes subscription links
4. **Connection Manager**: Controls VPN connection lifecycle
5. **UI Layer**: Platform-specific user interfaces

### Tech Stack

- Core: Go for performance and cross-platform compatibility
- Mobile: Flutter for cross-platform mobile development (planned)
- Desktop: Electron or Tauri for cross-platform desktop apps (planned)
- Backend: Optional cloud infrastructure for sync capabilities

## Current Status

✅ **Phase 1: Core Functionality Completed**
- Basic architecture with modular design
- Data models and interfaces
- Server management system
- Connection management
- Protocol handler framework
- Configuration management
- All protocol handlers with placeholder implementations

## Next Steps

### Phase 2: Enhanced Features (In Progress)
- Implement real protocol handlers using open-source libraries
- Complete subscription parsing functionality
- Implement QR code scanning
- Add data usage tracking
- Implement IPv6 support toggle
- Advanced settings and profile management

### Phase 3: User Interface
- Develop mobile UI with Flutter
- Develop desktop UI with Electron/Tauri
- Create UI components:
  - Server list view
  - Connection status panel
  - Settings screen
  - Statistics dashboard

## Project Structure

```
vpn/
├── src/
│   ├── api/           # REST API for UI integration
│   ├── cli/           # Command-line interface
│   ├── core/          # Core data models and interfaces
│   ├── managers/      # Business logic managers
│   ├── protocols/     # Protocol-specific implementations
│   ├── utils/         # Utility functions
│   └── main.go        # Application entry point
├── ui/
│   ├── mobile/        # Mobile UI (Flutter planned)
│   │   └── index.html # Mobile UI prototype
│   └── desktop/       # Desktop UI (Electron/Tauri planned)
│       └── index.html # Desktop UI prototype
├── docs/
│   ├── architecture.md      # Architecture documentation
│   ├── cli_guide.md         # CLI usage guide
│   ├── dependencies.md      # Project dependencies
│   ├── protocol_integration.md # Protocol integration guide
│   ├── protocol_testing.md  # Protocol testing guide
│   ├── roadmap.md           # Feature roadmap
│   └── ui_integration.md    # UI integration guide
├── assets/            # Images and other assets
├── bin/               # Compiled binaries
├── build.sh           # Build script for Unix-like systems
├── build.bat          # Build script for Windows
├── Makefile           # Build automation
├── go.mod             # Go module definition
└── test_protocols.go  # Protocol integration test script
```

## Building the Project

### Using Make (recommended)
```bash
make deps    # Install dependencies
make build   # Build the application
```

### Windows
```cmd
build.bat
```

### Linux/macOS
```bash
chmod +x build.sh
./build.sh
```

## Running the Application

The application supports three modes of operation:

### Command-line interface (CLI) mode
```bash
cd src
go run main.go --cli
```

See [CLI Guide](docs/cli_guide.md) for detailed usage instructions.

### API server mode
```bash
cd src
go run main.go --api
```

Then open your browser to http://localhost:8080 to access the web UI.

### Protocol testing mode
```bash
go run test_protocols.go
```

See [Protocol Testing Guide](docs/protocol_testing.md) for detailed testing instructions.

### Help
```bash
cd src
go run main.go --help
```

## API Endpoints

When running in API mode, the application exposes the following endpoints:

- `GET /api/servers` - List all servers
- `POST /api/servers` - Add a new server
- `GET /api/servers/{id}` - Get server details
- `PUT /api/servers/{id}` - Update a server
- `DELETE /api/servers/{id}` - Remove a server
- `POST /api/connect` - Connect to a server
- `POST /api/disconnect` - Disconnect from current server
- `GET /api/status` - Get connection status
- `GET /api/stats` - Get connection statistics
- `GET /api/config` - Get application configuration
- `PUT /api/config` - Update application configuration

## Protocol Integration

The application supports multiple VPN protocols through a modular handler system. See [Protocol Integration Guide](docs/protocol_integration.md) for details on how to integrate new protocols.

Currently implemented protocols:
- VMess (with enhanced implementation)
- VLESS
- Shadowsocks (with real library integration)
- Trojan
- Reality
- Hysteria2
- TUIC
- SSH

## Testing

The project includes comprehensive tests for all components:

### Unit Tests
```bash
cd src
go test ./...
```

### Protocol Integration Tests
```bash
go run test_protocols.go
```

See [Protocol Testing Guide](docs/protocol_testing.md) for detailed testing instructions.

## UI Prototypes

The project includes HTML prototypes of both desktop and mobile UIs:
- Desktop: [ui/desktop/index.html](ui/desktop/index.html)
- Mobile: [ui/mobile/index.html](ui/mobile/index.html)

These prototypes demonstrate the intended UI design and can be opened directly in a web browser.

## Contributing

We welcome contributions to this open-source project. Please see our roadmap in the [docs/roadmap.md](docs/roadmap.md) file for planned features and improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.