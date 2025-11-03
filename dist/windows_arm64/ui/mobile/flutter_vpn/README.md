# VPN Client - Flutter UI

A cross-platform mobile VPN client with a modern and user-friendly interface, built with Flutter.

## Features

- Cross-platform support (Android, iOS)
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
- Fastest server auto-selection
- Data usage display
- Subscription link import with deep linking
- QR code import
- IPv6 support
- Ad-free experience
- Dark/Light theme support

## Screenshots

| Server List | Connection Status | Quick Connect | Settings |
|-------------|-------------------|---------------|----------|
| ![Server List](assets/screenshots/server_list.png) | ![Status](assets/screenshots/status.png) | ![Quick Connect](assets/screenshots/quick_connect.png) | ![Settings](assets/screenshots/settings.png) |

## Getting Started

### Prerequisites

- Flutter SDK (2.17.0 or higher)
- Dart SDK
- Android Studio or Xcode for mobile development

### Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   ```

2. Navigate to the Flutter project directory:
   ```
   cd ui/mobile/flutter_vpn
   ```

3. Install dependencies:
   ```
   flutter pub get
   ```

4. Run the app:
   ```
   flutter run
   ```

## Project Structure

```
lib/
├── main.dart                 # App entry point
├── models/                   # Data models
│   ├── server.dart           # Server model
│   └── connection_status.dart # Connection status model
├── providers/                # State management providers
│   └── app_provider.dart     # Main app state provider
├── services/                 # Business logic services
│   ├── connection_service.dart # Connection management
│   └── server_service.dart   # Server management
├── utils/                    # Utility functions
│   └── format_utils.dart     # Data formatting utilities
├── screens/                  # Screen widgets
│   ├── home_screen.dart      # Main screen with bottom navigation
│   ├── settings_screen.dart  # Settings screen
│   ├── qr_scanner_screen.dart # QR code scanner
│   ├── add_server_screen.dart # Add new server
│   ├── server_details_screen.dart # Server details
│   └── subscription_screen.dart # Subscription import
└── widgets/                  # Reusable widgets
    ├── server_list.dart      # Server list widget
    ├── connection_status.dart # Connection status widget
    └── quick_connect.dart    # Quick connect widget
```

## Dependencies

- `flutter`: UI toolkit
- `provider`: State management
- `http`: For API communication
- `qr_code_scanner`: For QR code scanning
- `shared_preferences`: For local data storage
- `flutter_secure_storage`: For secure data storage

## Building for Production

### Android

1. Create a keystore for signing the app
2. Update `android/key.properties` with your keystore information
3. Build the app:
   ```
   flutter build apk --release
   ```

### iOS

1. Update iOS app settings in Xcode
2. Build the app:
   ```
   flutter build ios --release
   ```

## Architecture

The app follows a clean architecture pattern with separation of concerns:

1. **Models**: Data structures representing the app's entities
2. **Providers**: State management using Provider package
3. **Services**: Business logic and data operations
4. **Screens**: Full-page UI components
5. **Widgets**: Reusable UI components
6. **Utils**: Helper functions and utilities

## Security

- Uses `flutter_secure_storage` for sensitive data
- Implements secure communication protocols
- No analytics or tracking
- Open-source for transparency

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to all the open-source projects that made this possible
- Special thanks to the Flutter community