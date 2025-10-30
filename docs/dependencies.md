# Project Dependencies

This document outlines the dependencies used in the VPN Client project.

## Go Modules

The project uses Go modules for dependency management. The main dependencies are:

### Gorilla Mux
- **Purpose**: HTTP request multiplexer for the API server
- **Version**: v1.8.0
- **Usage**: Routing HTTP requests in the API server
- **License**: BSD-3-Clause

### (Future) V2Ray/XRay Core
- **Purpose**: Handle VMess/VLESS protocols
- **Usage**: Will be integrated to provide real protocol support
- **License**: MIT

### (Future) Shadowsocks-go
- **Purpose**: Handle Shadowsocks protocol
- **Usage**: Will be integrated to provide Shadowsocks support
- **License**: MIT

### (Future) Other protocol libraries
- **Purpose**: Handle other VPN protocols (Trojan, Reality, etc.)
- **Usage**: Will be integrated as needed
- **License**: Various

## Development Tools

### Go
- **Purpose**: Main programming language
- **Version**: 1.21+
- **Usage**: Core application development
- **License**: BSD-3-Clause

### Make
- **Purpose**: Build automation
- **Usage**: Simplify build processes across platforms
- **License**: GNU General Public License

### Git
- **Purpose**: Version control
- **Usage**: Source code management
- **License**: GNU General Public License

## Testing Tools

### Go Test
- **Purpose**: Unit testing framework
- **Usage**: Testing protocol handlers and managers
- **License**: BSD-3-Clause

## Future UI Dependencies

### Flutter
- **Purpose**: Cross-platform mobile UI development
- **Usage**: Mobile application interface
- **License**: BSD-3-Clause

### Electron
- **Purpose**: Cross-platform desktop UI development
- **Usage**: Desktop application interface
- **License**: MIT

### Tauri
- **Purpose**: Alternative cross-platform desktop UI development
- **Usage**: Desktop application interface (alternative to Electron)
- **License**: MIT or Apache 2.0

## Build Dependencies

### Build Scripts
- **Windows**: Batch files (.bat)
- **Unix/Linux/macOS**: Shell scripts (.sh)
- **Cross-platform**: Makefile

## Runtime Dependencies

### System Libraries
- None required for the core Go application

### Network Access
- Required for connecting to VPN servers
- Required for fetching subscription links

## Installation

### Go Dependencies
```bash
go mod tidy
```

This command will automatically download and install all required Go dependencies.

### Development Environment
1. Install Go 1.21 or later
2. Clone the repository
3. Run `go mod tidy` to install dependencies
4. Build and run the application

## Updating Dependencies

To update dependencies:
```bash
go get -u
go mod tidy
```

## Dependency Licenses

All dependencies are open-source with permissive licenses:
- MIT License
- BSD-3-Clause License
- Apache 2.0 License

These licenses allow for commercial use, modification, distribution, and patent use, with minimal restrictions.

## Security Considerations

All dependencies should be regularly updated to address security vulnerabilities:
- Monitor security advisories for Go modules
- Keep development tools updated
- Review dependency licenses for compliance