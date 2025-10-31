# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

## [Unreleased]

### Added
- Protocol handler implementations for VMess, VLESS, Trojan, Reality, Hysteria2, TUIC, SSH, Shadowsocks
- Server management capabilities
- Connection management system
- Configuration management
- Cross-platform support (Windows, Linux, macOS)