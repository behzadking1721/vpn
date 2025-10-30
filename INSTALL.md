# Installation Guide

This guide provides detailed instructions for installing the VPN Client on various platforms.

## System Requirements

### Windows
- Windows 7 or later (64-bit recommended)
- 200 MB available disk space
- Administrator privileges for installation

### Linux
- Most modern Linux distributions (Ubuntu 18.04+, Debian 10+, Fedora 30+, etc.)
- 200 MB available disk space
- sudo privileges for installation

### macOS
- macOS 10.12 or later
- 200 MB available disk space

## Installation Methods

### Windows Installation

#### Method 1: Using Installer (Recommended)
1. Download the installer (`vpn-client-installer.exe`) from the [releases page](https://github.com/your-org/vpn-client/releases)
2. Double-click the installer to start the installation process
3. Follow the installation wizard:
   - Accept the license agreement
   - Choose installation directory (or use default)
   - Select additional options (desktop shortcut, etc.)
4. Click "Install" and wait for the process to complete
5. Launch the application from Start Menu or Desktop shortcut

#### Method 2: Portable Installation
1. Download the Windows zip archive from the releases page
2. Extract the archive to a folder of your choice
3. Run `vpn-client.exe` directly from the extracted folder

### Linux Installation

#### Method 1: Using Package Manager (Ubuntu/Debian)
```bash
# Download the .deb package
wget https://github.com/your-org/vpn-client/releases/download/v1.0.0/vpn-client_1.0.0_amd64.deb

# Install the package
sudo dpkg -i vpn-client_1.0.0_amd64.deb

# If there are missing dependencies
sudo apt-get install -f
```

#### Method 2: Using Package Manager (Fedora/RHEL/CentOS)
```bash
# Download the .rpm package
wget https://github.com/your-org/vpn-client/releases/download/v1.0.0/vpn-client-1.0.0.x86_64.rpm

# Install the package
sudo rpm -i vpn-client-1.0.0.x86_64.rpm
# or with dnf
sudo dnf install vpn-client-1.0.0.x86_64.rpm
```

#### Method 3: Manual Installation
1. Download the Linux tar.gz archive from the releases page
2. Extract the archive:
   ```bash
   tar -xzf vpn-client-*.tar.gz
   ```
3. Navigate to the extracted directory:
   ```bash
   cd vpn-client-*
   ```
4. Run the application:
   ```bash
   ./vpn-client
   ```

### macOS Installation

#### Method 1: Using DMG Installer
1. Download the DMG file from the releases page
2. Double-click the DMG file to mount it
3. Drag the VPN Client application to your Applications folder
4. Launch the application from Applications folder

#### Method 2: Manual Installation
1. Download the macOS tar.gz archive from the releases page
2. Extract the archive:
   ```bash
   tar -xzf vpn-client-*.tar.gz
   ```
3. Navigate to the extracted directory:
   ```bash
   cd vpn-client-*
   ```
4. Run the application:
   ```bash
   ./vpn-client
   ```

## First-Time Setup

1. Launch the VPN Client application
2. The application will create necessary directories:
   - Configuration directory
   - Data storage directory
   - Log directory
3. Add your VPN servers:
   - Click "Add Server" button
   - Enter server details (address, port, protocol, credentials)
   - Save the server configuration
4. Connect to a server by selecting it and clicking "Connect"

## Troubleshooting

### Windows

#### Issue: "Windows protected your PC" message
Solution: Click "More info" and then "Run anyway"

#### Issue: Application fails to start
Solution: Make sure you have installed Visual C++ Redistributables

### Linux

#### Issue: Permission denied when running the application
Solution: Make the binary executable:
```bash
chmod +x vpn-client
```

#### Issue: Missing dependencies on minimal distributions
Solution: Install required libraries:
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install libgtk-3-0 libwebkit2gtk-4.0-37

# Fedora/RHEL/CentOS
sudo dnf install gtk3 webkit2gtk3
```

### macOS

#### Issue: "Application can't be opened" message
Solution: Allow the application in System Preferences:
1. Go to System Preferences > Security & Privacy
2. Click "Open Anyway" for VPN Client

#### Issue: Application is damaged
Solution: Remove quarantine attribute:
```bash
xattr -d com.apple.quarantine /Applications/VPN\ Client.app
```

## Uninstallation

### Windows

#### Method 1: Using Control Panel
1. Open Control Panel
2. Go to "Programs and Features"
3. Find "VPN Client" in the list
4. Click "Uninstall"

#### Method 2: Using the Uninstaller
1. Navigate to the installation directory
2. Run `uninstall.exe`

### Linux

#### For .deb packages:
```bash
sudo apt-get remove vpn-client
```

#### For .rpm packages:
```bash
sudo rpm -e vpn-client
# or with dnf
sudo dnf remove vpn-client
```

#### For manual installations:
1. Delete the directory where you extracted the application
2. Remove any created configuration files

### macOS

1. Drag the VPN Client application from Applications folder to Trash
2. Empty the Trash

## Configuration

The application stores configuration files in the following locations:

### Windows
```
%APPDATA%\VPN Client\
```

### Linux
```
~/.config/vpn-client/
```

### macOS
```
~/Library/Application Support/VPN Client/
```

## Support

For additional help, please:
1. Check the [FAQ](FAQ.md)
2. Review existing issues on GitHub
3. Create a new issue if your problem is not addressed

## Security Notes

- The application does not collect or transmit any personal data
- All configuration files are stored locally
- Connection credentials are encrypted when stored
- Always download the application from official sources