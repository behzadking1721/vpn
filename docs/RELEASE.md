# Release Process

This document describes the process for creating and publishing releases of the VPN Client.

## Versioning

We follow [Semantic Versioning](https://semver.org/) for version numbers in the format `MAJOR.MINOR.PATCH`:

- **MAJOR** version when you make incompatible API changes
- **MINOR** version when you add functionality in a backwards compatible manner
- **PATCH** version when you make backwards compatible bug fixes

## Release Process

### 1. Preparation

1. Ensure all changes for the release are merged into the `main` branch
2. Update the [CHANGELOG.md](file:///c%3A/Users/behza/OneDrive/Documents/vpn/CHANGELOG.md) with all notable changes
3. Update the version in the codebase if needed
4. Run all tests to ensure stability:
   ```bash
   make test
   ```

### 2. Create Release Tag

1. Create and push a new tag:
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```
   
   Or using GitHub CLI:
   ```bash
   gh release create v1.0.0 \
     --title "Version 1.0.0" \
     --notes "See [CHANGELOG.md](CHANGELOG.md) for changes" \
     --draft=false \
     --prerelease=false
   ```

### 3. Automated Build and Release

The GitHub Actions workflow in `.github/workflows/release.yml` will automatically:

1. Build the application for all platforms (Windows, Linux, macOS)
2. Build the CLI tool for all platforms
3. Create release artifacts
4. Publish the release on GitHub

### 4. Manual Release Steps (if needed)

If the automated release fails, you can manually create the builds:

1. Build for all platforms:
   ```bash
   make build-all
   make build-cli-all
   ```

2. Create the release directory structure:
   ```bash
   make release-structure VERSION=1.0.0
   ```

3. Manually upload artifacts to GitHub release

### 5. Post-Release

1. Verify all artifacts are uploaded correctly
2. Test downloads on different platforms
3. Update documentation if needed
4. Announce release (if applicable)

## Build System

### Makefile Targets

The project uses a Makefile for building and managing the release process:

```bash
# Build for all platforms
make build-all

# Build CLI for all platforms
make build-cli-all

# Clean build artifacts
make clean

# Run tests
make test

# Create release directory structure
make release-structure VERSION=1.0.0
```

### Build Variables

The build system supports the following variables:

- `VERSION` - Release version (default: dev)
- `BUILD_TIME` - Build timestamp (set automatically)
- `GIT_COMMIT` - Git commit hash (set automatically)

These variables are embedded into the binaries at build time.

## Platform-Specific Considerations

### Windows

- Builds produce `.exe` files
- Consider creating an installer using Inno Setup or NSIS
- Sign binaries for security and to avoid antivirus warnings

### Linux

- Builds produce executable binaries
- Consider creating packages for popular distributions (`.deb`, `.rpm`)
- Consider creating AppImage for universal Linux support

### macOS

- Builds produce executable binaries
- Consider creating a `.dmg` installer
- Sign binaries and request notarization from Apple

## Security Considerations

1. Sign all binaries for supported platforms
2. Verify dependencies before building
3. Scan binaries for vulnerabilities
4. Use reproducible builds when possible

### Signing Binaries

#### Windows
```bash
signtool sign /fd SHA256 /tr http://timestamp.digicert.com vpn-client.exe
```

#### macOS
```bash
codesign --deep --force --verify --verbose --sign "Developer ID Application: Your Name" vpn-client.app
```

#### Linux
```bash
gpg --detach-sign vpn-client
```

## Troubleshooting

### Build Failures

1. Ensure Go 1.21+ is installed
2. Run `make deps` to update dependencies
3. Check for platform-specific build issues

### Release Failures

1. Verify GitHub Actions permissions
2. Check that the tag format is correct (`v*`)
3. Ensure `GITHUB_TOKEN` secret is configured

### Common Issues

1. **Missing dependencies**: Run `make deps`
2. **Permission denied**: Check file permissions
3. **Platform-specific errors**: Ensure cross-compilation targets are supported