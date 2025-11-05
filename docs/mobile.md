# Mobile Application Documentation

## Overview

The mobile application for the VPN client is built with React Native, providing a cross-platform solution for both Android and iOS devices. It connects to the same backend API as the desktop application, ensuring consistent functionality across platforms.

## Features

- Cross-platform support (Android and iOS)
- Real-time connection status
- Server browsing and selection
- Subscription management
- Connection statistics
- Multi-language support (English, Persian, Chinese)
- Dark/Light theme support

## Architecture

The mobile application follows a component-based architecture using React Native. It communicates with the backend through REST API calls.

### Main Components

1. **App.js** - Main application component
2. **Navigation** - Tab-based navigation system
3. **Dashboard** - Connection status and statistics
4. **Servers** - Server listing and management
5. **Subscriptions** - Subscription management
6. **Settings** - Application settings

## Setup

### Prerequisites

- Node.js 16 or higher
- Android Studio (for Android development)
- Xcode (for iOS development)
- React Native CLI

### Installation

```bash
cd mobile
npm install
```

### Running the Application

#### Android

```bash
npx react-native run-android
```

#### iOS

```bash
npx react-native run-ios
```

## API Integration

The mobile application connects to the backend API at `http://localhost:8080`. For production deployments, this URL should be updated to point to the actual server.

### Available Endpoints

- `GET /api/servers` - List all servers
- `GET /api/servers/enabled` - List enabled servers
- `POST /api/connect` - Connect to a specific server
- `POST /api/connect/best` - Connect to the best server
- `POST /api/disconnect` - Disconnect from current server
- `GET /api/stats` - Get connection statistics
- `GET /api/subscriptions` - List all subscriptions
- `POST /api/subscriptions` - Add a new subscription
- `DELETE /api/subscriptions/{id}` - Delete a subscription
- `POST /api/servers/test-all-ping` - Test all servers

## UI Components

### Dashboard

The dashboard shows real-time connection statistics including:
- Download usage
- Upload usage
- Connection time
- Quick connect/disconnect buttons

### Servers

The servers tab displays a list of available VPN servers with:
- Server name and location
- Ping latency
- Protocol information
- Connection status
- Connect buttons

### Subscriptions

The subscriptions tab allows users to:
- Add new subscription URLs
- View existing subscriptions
- Update subscriptions
- Delete subscriptions

### Settings

The settings tab provides:
- Language selection (English, Persian, Chinese)
- Theme selection (Dark/Light)
- Other application preferences

## Data Management

The mobile application uses AsyncStorage for local data persistence:
- Theme preferences
- Language settings
- Cached server data

## Internationalization

The application supports three languages:
- English
- Persian (RTL support)
- Chinese

Translations are managed in the main App.js file.

## Styling

The application uses StyleSheet for styling with support for both light and dark themes. All components adapt to the selected theme.

## Performance Considerations

- Virtualized lists for efficient rendering of large server lists
- Connection pooling for API requests
- Caching of frequently accessed data
- Lazy loading of non-critical components

## Testing

To run tests:

```bash
cd mobile
npm test
```

## Building for Production

### Android

```bash
cd mobile
npx react-native build-android --mode=release
```

### iOS

```bash
cd mobile
npx react-native build-ios --mode=release
```

## Troubleshooting

### Common Issues

1. **API Connection Failures**
   - Ensure the backend server is running
   - Check that the API URL is correctly configured
   - Verify network connectivity

2. **Android Build Issues**
   - Make sure Android Studio is properly installed
   - Check that all required SDKs are installed
   - Ensure JAVA_HOME is set correctly

3. **iOS Build Issues**
   - Make sure Xcode is properly installed
   - Check that CocoaPods is installed and up to date
   - Ensure developer account is properly configured

### Debugging

To enable debugging:

1. Shake the device or press Cmd+D (iOS) or Ctrl+M (Android)
2. Select "Debug"
3. Open Chrome Developer Tools at http://localhost:8081/debugger-ui

## Future Improvements

- Push notifications for connection status
- Biometric authentication
- Offline mode with cached data
- Enhanced server sorting and filtering
- Map-based server selection
- Widget support for quick connect/disconnect