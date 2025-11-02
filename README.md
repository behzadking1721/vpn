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

## Installation

### Windows

Download the latest Windows installer from the [releases page](https://github.com/your-org/vpn-client/releases) and run it.

Alternatively, you can download the portable ZIP archive, extract it, and run `vpn-client.exe`.

### Linux

#### Debian/Ubuntu (DEB package)

Download the `.deb` package from the releases page and install it:

```bash
sudo dpkg -i vpn-client_*.deb
```

#### Other Linux distributions

Download the appropriate tar.gz archive for your architecture from the releases page:

```bash
tar -xzf vpn-client-linux_amd64-*.tar.gz
cd vpn-client-linux_amd64
./vpn-client
```

### macOS

Download the macOS ZIP archive from the releases page, extract it, and run the application.

Note: You may need to allow the application in System Preferences > Security & Privacy.

## Building from Source

### Building Prerequisites

- Go 1.21 or later
- GCC (for CGO-enabled builds)
- NSIS (for Windows installer creation, optional)
- zip/tar (for packaging)

### Building

To build for all supported platforms:

```bash
chmod +x scripts/build-all.sh
./scripts/build-all.sh
```

To build for a specific platform:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/vpn-client ./src

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/vpn-client.exe ./src

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/vpn-client ./src
```

### Packaging

To create packages for distribution:

```bash
chmod +x scripts/package.sh
./scripts/package.sh
```

## Deployment & Release Process

This project follows a structured release process:

1. **Versioning**: We use Semantic Versioning (SemVer) - vX.Y.Z
2. **Build**: Cross-compile for all supported platforms
3. **Package**: Create platform-specific packages (ZIP, tar.gz, DEB, etc.)
4. **Release**: Create GitHub release with all packages

### Creating a Release

1. Ensure all changes are committed and pushed
2. Update the version in relevant files:

   ```bash
   ./scripts/release.sh v1.0.0
   ```

3. Push the version tag:

   ```bash
   git push origin v1.0.0
   ```

4. The GitHub Actions workflow will automatically:
   - Run tests
   - Build for all platforms
   - Create packages
   - Create a GitHub release

### Supported Platforms

| Platform | Architectures       | Packaging              |
|----------|---------------------|------------------------|
| Windows  | amd64, 386, arm64   | Installer (NSIS), ZIP  |
| Linux    | amd64, 386, arm64, arm | DEB, tar.gz, ZIP    |
| macOS    | amd64, arm64        | ZIP                    |

### Phase 3: User Interface

- Develop mobile UI with Flutter
- Develop desktop UI with Electron/Tauri
- Create UI components:
  - Server list view
  - Connection status panel
  - Settings screen
  - Statistics dashboard

## Project Structure

```text
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

### Windows Build

```batch
build.bat
```

### Linux/macOS Build

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

Then open your browser to <http://localhost:8080> to access the web UI.

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

### Setting up Development Environment

#### Windows Setup

Run the setup script:

```batch
setup-dev.bat
```

#### Linux/macOS Setup

Run the setup script:

```bash
chmod +x setup-dev.sh
./setup-dev.sh
```

### Running Tests

#### Windows Tests

```batch
run-tests.bat
```

#### Linux/macOS Tests

```bash
chmod +x run-tests.sh
./run-tests.sh
```

See [Testing Guide](docs/testing-guide.md) for detailed testing instructions.

## UI Prototypes

The project includes HTML prototypes of both desktop and mobile UIs:
- Desktop: [ui/desktop/index.html](ui/desktop/index.html)
- Mobile: [ui/mobile/index.html](ui/mobile/index.html)

These prototypes demonstrate the intended UI design and can be opened directly in a web browser.

## Contributing

We welcome contributions to this open-source project. Please see our roadmap in the [docs/roadmap.md](docs/roadmap.md) file for planned features and improvements.

For development guidelines, see [Development Guide](docs/development.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
