# Statistics System Documentation

## Overview

The statistics system provides comprehensive data collection, analysis, and visualization capabilities for the VPN client application. It tracks connection metrics, data usage, session durations, and other key performance indicators to help users understand their VPN usage patterns.

## Architecture

The statistics system consists of several components:

1. **Stats Package** - Core statistics functionality in Go
2. **API Endpoints** - RESTful interface for statistics management
3. **Desktop UI** - Visualization of statistics in the desktop application
4. **Data Storage** - In-memory storage of statistics data

## Stats Package

The stats package is implemented in the [internal/stats](file:///c%3A/Users/behza/OneDrive/Documents/vpn/internal/stats) directory and provides the following key features:

### Data Structures

#### ConnectionStat
Represents statistics for a single ongoing connection:
- `Timestamp` - When the connection started
- `DataSent` - Amount of data sent in bytes
- `DataRecv` - Amount of data received in bytes
- `ServerID` - Identifier of the connected server
- `ServerName` - Name of the connected server

#### SessionStat
Represents statistics for a completed connection session:
- `StartedAt` - When the session started
- `EndedAt` - When the session ended
- `DataSent` - Total data sent during the session
- `DataRecv` - Total data received during the session
- `ServerID` - Identifier of the connected server
- `ServerName` - Name of the connected server

### StatsManager

The StatsManager is the core component that manages all statistics:

#### Features
- Real-time tracking of current connection statistics
- Storage of historical session data
- Calculation of aggregate statistics
- Data querying by time range
- Data formatting for visualization

## API Endpoints

The statistics system exposes the following RESTful API endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/stats/connection` | GET | Get current connection statistics |
| `/api/stats/sessions` | GET | Get all session statistics |
| `/api/stats/summary` | GET | Get a summary of all statistics |
| `/api/stats/daily` | GET | Get daily usage statistics |
| `/api/stats/chart` | GET | Get data formatted for charting |
| `/api/stats/clear` | POST | Clear all statistics |

### GET /api/stats/connection

Get current connection statistics.

**Response:**
```json
{
  "connected": true,
  "timestamp": "2023-01-01T12:00:00Z",
  "data_sent": 1024000,
  "data_recv": 2048000,
  "server_id": "server1",
  "server_name": "Example Server"
}
```

### GET /api/stats/sessions

Get all session statistics.

**Response:**
```json
[
  {
    "started_at": "2023-01-01T10:00:00Z",
    "ended_at": "2023-01-01T11:00:00Z",
    "data_sent": 512000,
    "data_recv": 1024000,
    "server_id": "server1",
    "server_name": "Example Server"
  }
]
```

### GET /api/stats/summary

Get a summary of all statistics.

**Response:**
```json
{
  "total_data_sent": 1536000,
  "total_data_recv": 3072000,
  "current_connection": {
    "timestamp": "2023-01-01T12:00:00Z",
    "data_sent": 1024000,
    "data_recv": 2048000,
    "server_id": "server1",
    "server_name": "Example Server"
  },
  "recent_sessions": [
    {
      "started_at": "2023-01-01T10:00:00Z",
      "ended_at": "2023-01-01T11:00:00Z",
      "data_sent": 512000,
      "data_recv": 1024000,
      "server_id": "server1",
      "server_name": "Example Server"
    }
  ]
}
```

### GET /api/stats/daily

Get daily usage statistics.

**Query Parameters:**
- `days` - Number of days to retrieve (default: 7)

**Response:**
```json
[
  {
    "timestamp": "2023-01-01T00:00:00Z",
    "data_sent": 1536000,
    "data_recv": 3072000
  }
]
```

### GET /api/stats/chart

Get data formatted for charting.

**Query Parameters:**
- `type` - Type of chart data (daily_usage, session_duration, data_comparison)
- `days` - Number of days for daily usage (default: 7)

**Response (daily_usage):**
```json
{
  "type": "daily_usage",
  "labels": ["2023-01-01", "2023-01-02"],
  "datasets": [
    {
      "label": "Data Sent",
      "data": [1536000, 2048000],
      "backgroundColor": "rgba(54, 162, 235, 0.2)",
      "borderColor": "rgba(54, 162, 235, 1)",
      "borderWidth": 1
    },
    {
      "label": "Data Received",
      "data": [3072000, 4096000],
      "backgroundColor": "rgba(255, 99, 132, 0.2)",
      "borderColor": "rgba(255, 99, 132, 1)",
      "borderWidth": 1
    }
  ]
}
```

### POST /api/stats/clear

Clear all statistics.

**Response:**
```json
{
  "status": "success",
  "message": "Statistics cleared successfully"
}
```

## Desktop UI Integration

The desktop UI provides a user-friendly interface for viewing and analyzing statistics:

### Features

- Real-time statistics dashboard
- Interactive charts for data visualization
- Session history table
- Data filtering by time range
- Statistics clearing functionality

### Implementation

The statistics UI is implemented as a separate HTML page ([stats.html](file:///c%3A/Users/behza/OneDrive/Documents/vpn/ui/desktop/stats.html)) with the following components:

1. **Stats Cards** - Key metrics display (total data sent/received, current session, active connections)
2. **Charts** - Interactive visualizations using Chart.js
3. **Sessions Table** - Detailed history of connection sessions
4. **Controls** - Filtering and management options

## Integration Points

### Connection Manager

The connection manager integrates with the statistics system to provide real-time data:

- Starting statistics tracking when a connection is established
- Updating data usage statistics during the connection
- Ending statistics tracking when a connection is terminated

### Data Flow

1. Connection Manager updates statistics through StatsManager
2. StatsManager stores and processes the data
3. API endpoints serve the data to the frontend
4. Desktop UI visualizes the data using charts and tables

## Usage Examples

### Backend Statistics Tracking

```go
// Create stats manager
statsManager := stats.NewStatsManager()

// Start tracking a connection
statsManager.StartConnection(serverID, serverName)

// Update connection statistics
statsManager.UpdateConnection(dataSent, dataRecv)

// End connection tracking
statsManager.EndConnection()
```

### Frontend Statistics Display

```javascript
// Get statistics summary
async function getStatsSummary() {
    const response = await fetch('/api/stats/summary');
    if (response.ok) {
        const summary = await response.json();
        updateStatsDisplay(summary);
    }
}

// Get chart data
async function getChartData(chartType) {
    const response = await fetch(`/api/stats/chart?type=${chartType}`);
    if (response.ok) {
        const chartData = await response.json();
        updateChart(chartData);
    }
}

// Clear statistics
async function clearStatistics() {
    const response = await fetch('/api/stats/clear', { method: 'POST' });
    if (response.ok) {
        refreshStats();
    }
}
```

## Chart Types

The system supports several chart types for data visualization:

### Daily Usage Chart
Shows data sent and received per day over a specified time period.

### Session Duration Chart
Displays the duration of recent connection sessions.

### Data Comparison Chart
Compares total data sent vs. received.

## Data Retention

Statistics are stored in-memory and will be lost when the application is restarted. For production use, a persistent storage solution should be implemented.

## Performance Considerations

- Statistics are updated in real-time with minimal performance impact
- Chart data is aggregated to reduce the amount of data transferred
- Session history is limited to prevent memory issues

## Testing

The statistics system should be tested for the following scenarios:

1. **Data Accuracy** - Verify correct calculation of statistics
2. **Real-time Updates** - Test real-time data updates
3. **Chart Generation** - Verify correct chart data formatting
4. **API Endpoints** - Test all statistics API endpoints
5. **UI Integration** - Test frontend visualization

## Future Improvements

1. **Persistent Storage** - Implement database storage for statistics
2. **Export Functionality** - Add ability to export statistics data
3. **Advanced Analytics** - Add predictive analytics and trends
4. **Custom Reports** - Allow users to create custom reports
5. **Alerts** - Add threshold-based alerts for data usage