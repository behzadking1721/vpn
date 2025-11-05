# Notification System Documentation

## Overview

The notification system provides real-time feedback to users about VPN connection status, server updates, and other important events. It supports both in-app notifications and system-level notifications for desktop and mobile platforms.

## Architecture

The notification system consists of several components:

1. **Backend Notification Manager** - Core notification handling in Go
2. **API Endpoints** - RESTful interface for notification management
3. **Desktop UI Notifications** - Browser-based notifications for the web interface
4. **Mobile Push Notifications** - Native push notifications for mobile apps

## Backend Implementation

### Notification Manager

The backend notification system is implemented in the [internal/notifications](file:///c%3A/Users/behza/OneDrive/Documents/vpn/internal/notifications) package with the following key components:

- `Notification` struct - Represents a single notification
- `NotificationManager` - Manages notification lifecycle
- Methods for adding, retrieving, and managing notifications

### Notification Types

The system supports four notification types:

1. **Info** - General information messages
2. **Success** - Success confirmation messages
3. **Warning** - Warning messages that require attention
4. **Error** - Error messages indicating problems

### Features

- Maximum notification limit (100 notifications)
- Read/unread status tracking
- Timestamp tracking
- Notification clearing functionality

## API Endpoints

The notification system exposes the following RESTful API endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/notifications` | GET | Retrieve all notifications |
| `/api/notifications/unread` | GET | Retrieve unread notifications |
| `/api/notifications/read` | POST | Mark a notification as read |
| `/api/notifications/read-all` | POST | Mark all notifications as read |
| `/api/notifications/clear` | POST | Clear all notifications |
| `/api/notifications/clear-read` | POST | Clear read notifications |

## Desktop Implementation

The desktop UI uses both in-app notifications and browser-based system notifications:

1. **In-App Notifications** - Custom notification UI displayed within the application
2. **System Notifications** - Browser Web Notifications API for system-level alerts

### Features

- Automatic request for notification permissions
- Notification timeout (5 seconds)
- Manual dismissal capability
- Visual differentiation by notification type

## Mobile Implementation

The mobile application uses the React Native Push Notification library for both local and push notifications:

### Features

- Channel-based notification organization (Android)
- Configurable notification settings
- Local notification support
- Notification persistence

### Configuration

Mobile notifications are configured with the following settings:

- Channel ID: `vpn-notifications`
- Channel Name: `VPN Notifications`
- Sound: Default system sound
- Vibration: Enabled

## Integration Points

### Connection Manager

The connection manager integrates with the notification system to provide real-time connection status updates:

- Connection initiated
- Connection successful
- Connection failed
- Disconnection initiated
- Disconnection successful
- Disconnection failed

### Server Manager

The server manager provides notifications for server-related operations:

- Server added
- Server updated
- Server deleted
- Ping test results

### Subscription Manager

The subscription manager sends notifications for subscription operations:

- Subscription added
- Subscription updated
- Subscription deleted
- Subscription parsing errors

## Usage Examples

### Backend Notification Creation

```go
// Create notification manager
notificationManager := notifications.NewNotificationManager(100)

// Add a success notification
notification := notificationManager.AddNotification(
    "Connection Success", 
    "Successfully connected to server", 
    notifications.Success
)
```

### Frontend Notification Display

```javascript
// Show in-app notification
showNotification("Connection Success", "Successfully connected to server", "success");

// Show system notification (if permissions granted)
if ("Notification" in window && Notification.permission === "granted") {
    new Notification("Connection Success", {
        body: "Successfully connected to server"
    });
}
```

### Mobile Notification Display

```javascript
// Show local push notification
PushNotification.localNotification({
    channelId: "vpn-notifications",
    title: "Connection Success",
    message: "Successfully connected to server",
    playSound: true,
    soundName: "default"
});
```

## Configuration

### Desktop

Desktop notifications use the browser's Notification API and require user permission. Users are prompted for permission when the application loads.

### Mobile

Mobile notifications require platform-specific configuration:

#### Android

Add the following to `AndroidManifest.xml`:

```xml
<uses-permission android:name="android.permission.VIBRATE" />
<uses-permission android:name="android.permission.RECEIVE_BOOT_COMPLETED"/>

<receiver android:name="com.dieam.reactnativepushnotification.modules.RNPushNotificationBootEventReceiver">
    <intent-filter>
        <action android:name="android.intent.action.BOOT_COMPLETED" />
    </intent-filter>
</receiver>
```

#### iOS

Enable push notifications in the iOS project settings and add the following to `AppDelegate.m`:

```objc
#import <UserNotifications/UserNotifications.h>
#import <RNCPushNotificationIOS.h>

// In didFinishLaunchingWithOptions
UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
center.delegate = self;
```

## Testing

The notification system should be tested for the following scenarios:

1. **Permission Handling** - Test behavior when notification permissions are granted/denied
2. **Notification Display** - Verify notifications appear correctly in all supported formats
3. **API Endpoints** - Test all notification management endpoints
4. **Integration Points** - Verify notifications are sent from all integrated components
5. **Edge Cases** - Test behavior when maximum notification limit is reached

## Future Improvements

1. **Notification Categories** - Add support for categorizing notifications
2. **Notification Actions** - Add actionable notifications with buttons
3. **Notification History** - Persistent storage of notification history
4. **Custom Sounds** - Support for custom notification sounds
5. **Rich Notifications** - Support for images and rich content in notifications
6. **Scheduling** - Support for scheduled notifications