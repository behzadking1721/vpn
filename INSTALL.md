# VPN Client Installation Guide

## 🖥️ Desktop Installation

### Windows

#### Installer (Recommended)
1. Download `vpn-client-windows-VERSION.exe`
2. Run the installer and follow the installation wizard
3. Launch the application from Start Menu or Desktop shortcut

#### Portable
1. Download `vpn-client-windows-VERSION.zip`
2. Extract to a folder of your choice
3. Run `vpn-client.exe`

### macOS

#### Installer (Recommended)
1. Download `vpn-client-macos-VERSION.pkg`
2. Run the installer and follow the installation wizard
3. Launch the application from Applications folder

#### Portable
1. Download `vpn-client-darwin-VERSION.tar.gz`
2. Extract to a folder of your choice
3. Run `vpn-client` binary

### Linux

#### Debian/Ubuntu (DEB Package)
1. Download `vpn-client_VERSION_amd64.deb`
2. Install using:
   ```bash
   sudo dpkg -i vpn-client_VERSION_amd64.deb
   ```
3. Launch from applications menu or run `vpn-client` in terminal

#### Red Hat/Fedora (RPM Package)
1. Download `vpn-client-VERSION.x86_64.rpm`
2. Install using:
   ```bash
   sudo rpm -i vpn-client-VERSION.x86_64.rpm
   ```
3. Launch from applications menu or run `vpn-client` in terminal

#### Portable
1. Download `vpn-client-linux-VERSION.tar.gz`
2. Extract to a folder of your choice
3. Run `./vpn-client`

## 📱 Android Installation

### From Google Play Store (Recommended)
[![Google Play](https://play.google.com/intl/en_us/badges/static/images/badges/en_badge_web_generic.png)](https://play.google.com/store/apps/details?id=com.yourapp.vpnclient)

### Manual Installation (APK)
1. Download the appropriate APK file:
   - `vpn-client-VERSION-arm64.apk` for modern Android devices (64-bit ARM)
   - `vpn-client-VERSION-arm.apk` for older Android devices (32-bit ARM)
   - `vpn-client-VERSION-x64.apk` for emulators or x86-based devices
2. Enable "Install from unknown sources" in your Android settings
3. Open the downloaded APK file and follow installation prompts

### System Requirements
- Android 8.0 (Oreo) or higher
- Minimum 50MB of free storage space
- Internet connection for server list updates

## 🔧 Post-Installation Setup

1. Launch the VPN Client application
2. Add your VPN servers using one of these methods:
   - Enter server details manually
   - Import subscription links
   - Scan QR codes
3. Select a server and connect

## 🔐 Permissions

The application requires the following permissions:

### Desktop
- Network access: To establish VPN connections
- System tray access: To run in background

### Android
- INTERNET: To establish VPN connections
- ACCESS_NETWORK_STATE: To monitor network connectivity
- FOREGROUND_SERVICE: To maintain VPN service in background
- WAKE_LOCK: To prevent device from sleeping during connection
- VIBRATE: For notification alerts

## 🔄 Updates

The application supports automatic updates. You can also manually check for updates in the application settings.

## 🆘 Troubleshooting

If you encounter issues:

1. Check that your system meets the minimum requirements
2. Ensure you have proper permissions to run the application
3. Check firewall/antivirus settings
4. For connection issues, verify server details
5. Consult the FAQ or create an issue on GitHub

For more help, visit our [GitHub Issues](https://github.com/your-username/vpn-client/issues) page.