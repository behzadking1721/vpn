# Release Process

This document explains how to create a new release for the VPN Client application using GitHub Actions.

## Prerequisites

1. Ensure all changes are committed and pushed to the repository
2. Decide on the version number following [Semantic Versioning](https://semver.org/) (vX.Y.Z format)

## Creating a Release

### 1. Create and Push a Tag

Create a new tag with the version number and push it to GitHub:

```bash
# Create a new tag (replace X.Y.Z with actual version)
git tag v1.0.0

# Push the tag to GitHub
git push origin v1.0.0
```

This will automatically trigger the GitHub Actions workflow defined in [.github/workflows/release.yml](../../.github/workflows/release.yml).

### 2. Monitor the Release Process

You can monitor the progress of the release process in the [Actions tab](../../actions) on GitHub. The workflow consists of several jobs:

1. **test** - Runs all tests to ensure code quality
2. **build** - Compiles the application for all supported platforms
3. **package** - Creates distributable packages for each platform
4. **release** - Creates a GitHub release and uploads all packages

### 3. Verify the Release

Once the workflow completes successfully:

1. Navigate to the [Releases page](../../releases)
2. Find your new release
3. Verify that all expected assets are attached:
   - Windows ZIP package
   - Linux tar.gz packages (multiple architectures)
   - macOS tar.gz packages (multiple architectures)

## Supported Platforms

The release workflow builds and packages the application for the following platforms:

| Operating System | Architectures | Package Format |
|------------------|---------------|----------------|
| Windows          | amd64         | ZIP            |
| Linux            | amd64, arm64  | tar.gz         |
| macOS            | amd64, arm64  | tar.gz         |

## Manual Release Process

If you need to create a release manually (without using GitHub Actions):

1. Build for all platforms:
   ```bash
   # From project root
   chmod +x scripts/build-all.sh
   ./scripts/build-all.sh 1.0.0
   ```

2. Package the binaries:
   ```bash
   chmod +x scripts/package.sh
   ./scripts/package.sh 1.0.0
   ```

3. Create a release on GitHub and upload the packages from the `packages/` directory.

## Versioning Scheme

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version when you make incompatible API changes
- **MINOR** version when you add functionality in a backward compatible manner
- **PATCH** version when you make backward compatible bug fixes

Example version progression:
- v1.0.0 - Initial release
- v1.0.1 - Bug fixes
- v1.1.0 - New features
- v2.0.0 - Breaking changes

## Troubleshooting

### Release workflow fails

1. Check the Actions tab for detailed error messages
2. Common issues:
   - Tests failing - Fix the code or tests
   - Build errors - Check for platform-specific compilation issues
   - Packaging errors - Verify directory structure and file permissions

### Release created but assets missing

1. Check if the packaging step succeeded
2. Verify that all build artifacts were created correctly
3. Confirm there was enough space for the build process

## Release Artifacts Structure

Each packaged release contains:
```
vpn-client-{platform}_{arch}/
├── vpn-client[.exe]          # Main executable
├── ui/                       # UI files
├── config/                   # Configuration directory
├── data/                     # Data directory
└── logs/                     # Logs directory
```

This structure ensures that users have everything needed to run the application after extracting the package.