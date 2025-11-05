# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Protocol handler implementations for VMess, VLESS, Trojan, Reality, Hysteria2, TUIC, SSH, Shadowsocks
- Server management capabilities
- Connection management system
- Configuration management
- Cross-platform support (Windows, Linux, macOS)
- RESTful API for programmatic control
- Graphical user interface (HTML/CSS/JS)
- Command-line interface
- Subscription management with automatic server import
- Smart server selection based on ping and performance
- Comprehensive test suite including unit, integration, and end-to-end tests
- Detailed documentation covering API, usage, testing, and development
- Performance testing guidelines
- Security testing guidelines
- Cross-platform testing guidelines

### Changed
- Enhanced WireGuard implementation with wgctrl library
- Improved subscription parser to handle various formats
- Enhanced server manager with CRUD operations
- Improved connection manager with better status handling
- Enhanced API with comprehensive endpoints
- Improved UI with dashboard and alerts pages
- Enhanced testing with coverage reports

### Fixed
- Various bug fixes in protocol implementations
- Fixed import paths in multiple files
- Fixed server selection algorithms
- Fixed connection handling issues

## [1.0.0] - 2025-10-31

### Added
- Initial release of VPN Client
- Basic Windows executable with version display
- Support for command-line arguments (--version, --help)
- Packaging script for distribution
- Git tag and release preparation

### Changed
- Fixed import paths in main.go to use relative paths
- Updated go.mod module name to match project structure
- Simplified main.go for initial release

### Fixed
- Build script issues with absolute paths
- Dependency resolution issues in go.mod