# Frequently Asked Questions (FAQ)

## General Questions

### What is this VPN client?
This is a multi-platform VPN client that supports various VPN protocols including VMess, Shadowsocks, Trojan, VLESS, Reality, Hysteria, TUIC, SSH, and WireGuard. It provides a unified interface to manage and connect to different VPN servers.

### Is this application free?
Yes, this is an open-source application and is completely free to use. The source code is available on GitHub under the MIT license.

### Which platforms are supported?
The application supports:
- Windows 7 and later
- Linux (most distributions)
- macOS 10.12 and later

### How do I add a VPN server?
1. Launch the application
2. Click the "Add Server" button
3. Enter the server details:
   - Server name (for your reference)
   - Address and port
   - Protocol (VMess, Shadowsocks, etc.)
   - Authentication credentials
4. Save the configuration

### How do I connect to a VPN server?
1. Select a server from the server list
2. Click the "Connect" button
3. Wait for the connection to establish
4. The connection status will show as "Connected" when successful

## Technical Questions

### Does the application log my internet activity?
No, the application does not log your internet activity. It only stores connection logs locally for troubleshooting purposes, which can be cleared at any time.

### How is my data protected?
- All stored server configurations are encrypted
- Connection credentials are never transmitted to third parties
- The application uses secure implementations of VPN protocols
- No analytics or telemetry data is collected

### Can I use this with my existing VPN service?
Yes, as long as your VPN service supports one of the protocols supported by this client (VMess, Shadowsocks, Trojan, VLESS, Reality, Hysteria, TUIC, SSH, or WireGuard).

### How do I update the application?
The application has an automatic update mechanism:
1. The application checks for updates on startup
2. If an update is available, you'll be notified
3. Follow the prompts to download and install the update

You can also manually download updates from the releases page on GitHub.

## Troubleshooting

### The application won't start
Try these solutions:
1. Check that your system meets the minimum requirements
2. Restart your computer
3. Reinstall the application
4. Check the logs in the log directory for error messages

### I can't connect to my VPN server
Try these solutions:
1. Verify your server configuration (address, port, credentials)
2. Check your internet connection
3. Try connecting to a different server
4. Check if your firewall or antivirus is blocking the connection
5. Review the connection logs for error messages

### The connection is slow
Try these solutions:
1. Try connecting to a different server
2. Check your internet connection speed
3. Check if other applications are using bandwidth
4. Try a different VPN protocol
5. Check server ping times in the server list

### I'm having issues on Windows
Common Windows issues and solutions:
1. "Windows protected your PC" - Click "More info" then "Run anyway"
2. Missing Visual C++ Redistributables - Download and install from Microsoft
3. Antivirus blocking - Add the application to your antivirus whitelist

### I'm having issues on Linux
Common Linux issues and solutions:
1. Permission denied - Make the binary executable with `chmod +x vpn-client`
2. Missing dependencies - Install required libraries (gtk3, webkit2gtk)
3. Display server issues - Make sure you're running X11 or have Wayland support

### I'm having issues on macOS
Common macOS issues and solutions:
1. "Application can't be opened" - Allow in System Preferences > Security & Privacy
2. "Application is damaged" - Remove quarantine attribute with `xattr -d com.apple.quarantine`
3. Gatekeeper issues - Use `sudo xattr -rd com.apple.quarantine /Applications/VPN\ Client.app`

## Security Questions

### Is my data encrypted?
Yes, all stored server configurations are encrypted using industry-standard encryption algorithms. Connection data is encrypted by the VPN protocols themselves.

### Can third parties access my data?
No, the application does not transmit any data to third parties. All data remains on your local device unless you explicitly configure it otherwise.

### How often is the application updated?
We release updates regularly to:
- Fix bugs and security issues
- Add new features
- Improve performance
- Add support for new protocols

Major updates are typically released every few months, with minor updates as needed.

### How can I verify the application's integrity?
You can verify the integrity of downloaded files using checksums provided with each release. Additionally, all releases are signed, and you can verify the signatures using GPG.

## Development Questions

### Can I contribute to the project?
Yes! We welcome contributions. Please:
1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

### How do I build the application from source?
See the [README.md](README.md) file for build instructions.

### How do I report a bug?
Please create an issue on GitHub with:
1. A clear description of the problem
2. Steps to reproduce the issue
3. Your operating system and application version
4. Any relevant log files

### How do I request a feature?
Please create an issue on GitHub with:
1. A clear description of the feature
2. Why you think it would be useful
3. Any implementation suggestions (optional)

## Legal Questions

### Is this application legal to use?
The application itself is legal to use. However, the legality of using VPN services depends on your jurisdiction. Please check local laws and regulations.

### Is there a warranty?
This software is provided "as is", without warranty of any kind, express or implied. See the [LICENSE](LICENSE) file for details.

### Can I use this in commercial settings?
Yes, the MIT license allows for commercial use. See the [LICENSE](LICENSE) file for details.

### Who owns the copyright?
The copyright is held by the contributors to the project. See the [LICENSE](LICENSE) file for details.