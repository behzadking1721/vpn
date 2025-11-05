# Automatic Updater Documentation

## Overview

The automatic updater system provides functionality to automatically update VPN server lists from subscriptions at configurable intervals. This ensures that users always have access to the latest servers without manual intervention.

## Architecture

The updater system consists of several components:

1. **Updater Package** - Core updater functionality in Go
2. **API Endpoints** - RESTful interface for updater management
3. **Desktop UI Integration** - Configuration and control in the desktop application
4. **Background Service** - Periodic execution of update tasks

## Updater Package

The updater package is implemented in the [internal/updater](file:///c%3A/Users/behza/OneDrive/Documents/vpn/internal/updater) directory and provides the following key features:

### Core Components

#### Updater
The main updater component that manages the automatic update process:
- Configurable update intervals
- Enable/disable functionality
- Manual trigger capability
- Integration with server and subscription managers

#### Config
Configuration structure for the updater:
- `Interval` - Update interval as a time.Duration
- `Enabled` - Whether automatic updates are enabled

### Features

- **Automatic Updates** - Periodically updates server lists from subscriptions
- **Manual Trigger** - Allows manual triggering of updates
- **Configurable Intervals** - Supports various update intervals (hourly, daily, etc.)
- **Error Handling** - Gracefully handles update failures
- **Logging** - Comprehensive logging of update activities
- **Notification** - Sends notifications about update results

## API Endpoints

The updater system exposes the following RESTful API endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/updater/status` | GET | Get the current updater status |
| `/api/updater/config` | POST | Update the updater configuration |
| `/api/updater/update` | POST | Manually trigger an update |

### GET /api/updater/status

Get the current updater status.

**Response:**
```json
{
  "enabled": true,
  "interval": "24h0m0s"
}
```

### POST /api/updater/config

Update the updater configuration.

**Request Body:**
```json
{
  "enabled": true,
  "interval": "24h"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Updater configuration updated successfully"
}
```

### POST /api/updater/update

Manually trigger a subscription update.

**Response:**
```json
{
  "status": "success",
  "message": "Subscription update triggered successfully"
}
```

## Desktop UI Integration

The desktop UI provides a user-friendly interface for configuring and controlling automatic updates:

### Features

- Enable/disable automatic updates
- Configure update interval
- Manually trigger updates
- View updater status

### Implementation

The desktop UI uses JavaScript to communicate with the updater API endpoints and provides a settings panel for configuration.

## Configuration

### Update Intervals

The updater supports the following intervals:
- Hourly (1h)
- Every 6 Hours (6h)
- Every 12 Hours (12h)
- Daily (24h)
- Every 3 Days (72h)

### Default Configuration

By default, the updater is configured to:
- Be enabled
- Update every 24 hours

## Integration Points

### Server Manager

The updater integrates with the server manager to:
- Retrieve server lists
- Update server information
- Remove obsolete servers

### Subscription Manager

The updater integrates with the subscription manager to:
- Retrieve subscription lists
- Update subscriptions
- Handle subscription errors

## Usage Examples

### Backend Configuration

```go
// Create updater configuration
config := updater.Config{
    Interval: 24 * time.Hour,
    Enabled:  true,
}

// Create updater
updater := updater.NewUpdater(serverManager, subscriptionManager, config, logger)

// Start automatic updates
updater.Start()

// Stop automatic updates when shutting down
defer updater.Stop()
```

### Frontend Configuration

```javascript
// Enable/disable automatic updates
fetch('/api/updater/config', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({ enabled: true })
});

// Set update interval
fetch('/api/updater/config', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({ interval: '24h' })
});

// Manually trigger update
fetch('/api/updater/update', {
    method: 'POST'
});
```

### Checking Status

```javascript
// Get updater status
fetch('/api/updater/status')
    .then(response => response.json())
    .then(status => {
        console.log('Updater enabled:', status.enabled);
        console.log('Update interval:', status.interval);
    });
```

## Error Handling

The updater system handles various error conditions:

1. **Network Errors** - Failed subscription updates due to network issues
2. **Parse Errors** - Invalid subscription data
3. **Storage Errors** - Database write failures
4. **Configuration Errors** - Invalid configuration parameters

All errors are logged and do not stop the updater service.

## Testing

The updater system should be tested for the following scenarios:

1. **Configuration Changes** - Test changing update intervals and enabling/disabling
2. **Manual Updates** - Test manual trigger functionality
3. **Automatic Updates** - Test periodic update execution
4. **Error Conditions** - Test behavior when updates fail
5. **Integration** - Test integration with server and subscription managers

## Performance Considerations

- Updates are performed in the background without affecting user experience
- Network operations are optimized to minimize bandwidth usage
- Server resources are efficiently managed during update processes

## Security Considerations

- All API endpoints are protected by the same security mechanisms as other endpoints
- Subscription data is validated before processing
- Update operations are logged for audit purposes

## Future Improvements

1. **Smart Update Scheduling** - Schedule updates based on user activity patterns
2. **Bandwidth Throttling** - Limit bandwidth usage during updates
3. **Differential Updates** - Only download changed server data
4. **Update Prioritization** - Prioritize updates for frequently used servers
5. **Offline Support** - Cache updates for offline usage