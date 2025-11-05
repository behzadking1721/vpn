# VPN Client Usage Guide

This document provides instructions on how to use the VPN Client application in different modes: CLI, API, and GUI.

## Table of Contents

1. [CLI Mode](#cli-mode)
2. [API Mode](#api-mode)
3. [GUI Mode](#gui-mode)
4. [Configuration](#configuration)
5. [Managing Servers](#managing-servers)
6. [Managing Subscriptions](#managing-subscriptions)
7. [Connection Management](#connection-management)

## CLI Mode

The VPN Client can be run in command-line interface (CLI) mode for automated tasks and scripting.

### Starting the Application in CLI Mode

To run the application in CLI mode, use the following command:

```bash
# Using the make command
make run-cli

# Directly running the Go application
go run src/main.go --cli

# Or if you have the binary
./vpn-client --cli
```

### CLI Commands

The CLI mode supports the following commands:

- `--help`: Display help information
- `--version`: Display the application version
- `--config <path>`: Specify a custom configuration file path
- `--cli`: Run in CLI mode

### Example CLI Usage

```bash
# Run with default configuration
vpn-client --cli

# Run with custom configuration
vpn-client --cli --config /path/to/config.json

# Display help
vpn-client --cli --help

# Display version
vpn-client --cli --version
```

### CLI Mode Features

In CLI mode, you can perform the following operations:

1. **List Servers**:
   ```bash
   vpn-client --cli servers list
   ```

2. **Connect to a Server**:
   ```bash
   vpn-client --cli connect --server-id server1
   ```

3. **Disconnect**:
   ```bash
   vpn-client --cli disconnect
   ```

4. **Check Status**:
   ```bash
   vpn-client --cli status
   ```

Note: CLI commands may vary based on implementation. Refer to the help command for the most up-to-date information.

## API Mode

The VPN Client can be run as an API server to allow programmatic control.

### Starting the API Server

To start the API server, use one of the following commands:

```bash
# Using the make command
make run-api

# Directly running the Go application
go run src/main.go --api

# Or if you have the binary
./vpn-client --api
```

### API Server Configuration

By default, the API server runs on `localhost:8080`. You can change this by modifying the configuration file.

Example configuration file (`config/settings.json`):

```json
{
  "api": {
    "address": "localhost:8080"
  },
  "data": {
    "directory": "./data"
  }
}
```

### Accessing the API

Once the server is running, you can access the API endpoints. See [API Documentation](api.md) for detailed information about available endpoints.

### Example API Usage

```bash
# Get server status
curl http://localhost:8080/api/status

# List all servers
curl http://localhost:8080/api/servers

# Connect to a server
curl -X POST http://localhost:8080/api/connect \
  -H "Content-Type: application/json" \
  -d '{"server_id": "your-server-id"}'

# Add a new server
curl -X POST http://localhost:8080/api/servers \
  -H "Content-Type: application/json" \
  -d '{
    "id": "new-server",
    "name": "New Server",
    "host": "server.example.com",
    "port": 443,
    "protocol": "vmess",
    "config": {
      "user_id": "user-id-here"
    },
    "enabled": true
  }'

# Test all server pings
curl -X POST http://localhost:8080/api/servers/test-all-ping

# Connect to the fastest server
curl -X POST http://localhost:8080/api/connect/fastest
```

### API Security Considerations

When running the API server in production environments:

1. **Firewall Configuration**: Restrict access to the API port (default 8080) to only trusted sources
2. **Reverse Proxy**: Use a reverse proxy like Nginx or Apache for additional security
3. **Authentication**: Future versions may include authentication; for now, ensure the API is not exposed to public networks
4. **HTTPS**: Use HTTPS in production environments to encrypt data in transit

## GUI Mode

The VPN Client includes a graphical user interface for desktop environments.

### Starting the GUI

To start the GUI, run:

```bash
# Using the make command
make run

# Directly running the Go application
go run src/main.go

# Or if you have the binary
./vpn-client
```

### GUI Features

The GUI provides the following features:

1. **Connection Management**:
   - Connect/disconnect from VPN servers
   - View connection status and statistics
   - Select servers manually or automatically

2. **Server Management**:
   - Add, edit, and delete servers
   - Enable/disable servers
   - Test server ping
   - Find the best server automatically

3. **Subscription Management**:
   - Add subscriptions via URL
   - Automatically import servers from subscriptions
   - Update subscriptions

4. **Dashboard**:
   - View connection history
   - Monitor data usage
   - View server performance metrics

### Using the GUI

1. **Connecting to a Server**:
   - Select a server from the server list
   - Click the "Connect" button
   - Wait for the connection to establish

2. **Adding a Server**:
   - Click the "Add Server" button
   - Fill in the server details
   - Click "Save"

3. **Importing from Subscription**:
   - Go to the subscription section
   - Enter the subscription URL
   - Click "Import"

4. **Testing Server Performance**:
   - Click the "Test All Servers" button to test ping for all enabled servers
   - Use the "Connect Fastest" button to connect to the server with the lowest ping
   - Use the "Connect Best" button to connect to the server with the best overall performance

5. **Viewing Dashboard**:
   - Navigate to the Dashboard tab
   - View connection statistics and server performance metrics

### GUI Keyboard Shortcuts

The GUI supports the following keyboard shortcuts:

- `Ctrl+C`: Connect to selected server
- `Ctrl+D`: Disconnect from current server
- `Ctrl+R`: Refresh server list
- `Ctrl+T`: Test all server pings
- `Ctrl+Q`: Quit the application

## Configuration

The VPN Client uses a JSON configuration file located at `config/settings.json` by default.

### Configuration File Structure

```json
{
  "api": {
    "address": "localhost:8080"
  },
  "data": {
    "directory": "./data"
  }
}
```

### Configuration Options

1. **API Configuration**:
   - `address`: The address and port for the API server (default: localhost:8080)

2. **Data Configuration**:
   - `directory`: The directory where application data is stored (default: ./data)

### Custom Configuration

You can specify a custom configuration file using the `--config` flag:

```bash
vpn-client --config /path/to/custom-config.json
```

### Environment Variables

The VPN Client also supports the following environment variables:

- `VPN_CLIENT_CONFIG`: Path to the configuration file
- `VPN_CLIENT_DATA_DIR`: Directory for storing application data

Environment variables take precedence over configuration file values.

## Managing Servers

Servers are the core components of the VPN Client. You can manage them through the GUI, CLI, or API.

### Adding a Server

Servers can be added in several ways:

1. **Manual Entry**: Enter server details manually in the GUI or through the API
2. **Subscription Import**: Import servers from a subscription URL
3. **QR Code**: Scan a QR code containing server configuration (in GUI mode)

### Server Properties

Each server has the following properties:

- **ID**: Unique identifier for the server
- **Name**: Human-readable name for the server
- **Host**: Server hostname or IP address
- **Port**: Server port number
- **Protocol**: VPN protocol (e.g., vmess, shadowsocks, trojan, wireguard)
- **Config**: Protocol-specific configuration
- **Enabled**: Whether the server is enabled for use
- **Ping**: Server latency in milliseconds

### Supported Protocols

The VPN Client supports the following protocols:

1. **VMess**: A protocol for encrypted communications
2. **Shadowsocks**: A secure socks5 proxy
3. **Trojan**: A protocol that disguises traffic as HTTPS
4. **WireGuard**: A modern, fast VPN protocol (planned support)

### Testing Server Performance

You can test server performance in the following ways:

1. **Ping Test**: Measure basic latency
2. **Speed Test**: Measure connection speed (future feature)
3. **Auto-select**: Let the application choose the best server

### Example Server Configurations

1. **VMess Server**:
   ```json
   {
     "id": "vmess-server",
     "name": "VMess Server",
     "host": "vmess.example.com",
     "port": 443,
     "protocol": "vmess",
     "config": {
       "user_id": "user-id-here",
       "alter_id": 0,
       "security": "auto"
     },
     "enabled": true
   }
   ```

2. **Shadowsocks Server**:
   ```json
   {
     "id": "ss-server",
     "name": "Shadowsocks Server",
     "host": "ss.example.com",
     "port": 8388,
     "protocol": "shadowsocks",
     "config": {
       "method": "aes-128-gcm",
       "password": "password-here"
     },
     "enabled": true
   }
   ```

3. **Trojan Server**:
   ```json
   {
     "id": "trojan-server",
     "name": "Trojan Server",
     "host": "trojan.example.com",
     "port": 443,
     "protocol": "trojan",
     "config": {
       "password": "password-here"
     },
     "enabled": true
   }
   ```

## Managing Subscriptions

Subscriptions allow you to import multiple servers from a single URL.

### Adding a Subscription

1. **In GUI**: Go to the subscription section and enter the URL
2. **In API**: Use the `/api/subscriptions` endpoint
3. **In CLI**: Use the appropriate CLI command (if available)

### Subscription URL Formats

The VPN Client supports the following subscription URL formats:

1. **Base64 Encoded**: A base64-encoded list of server configurations
2. **VMess URLs**: A list of vmess:// URLs
3. **Shadowsocks URLs**: A list of ss:// URLs
4. **Trojan URLs**: A list of trojan:// URLs

### Subscription Updates

Subscriptions can be configured to auto-update. When a subscription is updated, new servers are added and existing servers are updated.

### Example Subscription Usage

```bash
# Add a subscription via API
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Subscription",
    "url": "https://example.com/subscription",
    "auto_update": true
  }'

# Update servers from a subscription
curl -X POST http://localhost:8080/api/subscriptions/sub1/update
```

## Connection Management

The VPN Client provides several ways to manage connections.

### Connection Methods

1. **Manual**: Select a server and connect manually
2. **Fastest**: Connect to the server with the lowest ping
3. **Best**: Connect to the server with the best overall performance

### Connection Status

The application provides real-time connection status information:

- **Status**: Connected/Disconnected/Connecting
- **Uptime**: How long the connection has been active
- **Data Usage**: Amount of data sent and received
- **Current Server**: Details of the connected server

### Example Connection Management

```bash
# Connect to a specific server
curl -X POST http://localhost:8080/api/connect \
  -H "Content-Type: application/json" \
  -d '{"server_id": "server1"}'

# Connect to the fastest server
curl -X POST http://localhost:8080/api/connect/fastest

# Connect to the best server
curl -X POST http://localhost:8080/api/connect/best

# Disconnect
curl -X POST http://localhost:8080/api/disconnect

# Get status
curl http://localhost:8080/api/status

# Get statistics
curl http://localhost:8080/api/stats
```

### Disconnecting

To disconnect, simply click the "Disconnect" button in the GUI or use the appropriate API endpoint.

## Troubleshooting

### Common Issues

1. **Connection Failures**:
   - Check server configuration
   - Verify network connectivity
   - Try a different server
   - Check firewall settings

2. **Performance Issues**:
   - Test different servers
   - Check your internet connection
   - Use the ping test feature
   - Consider using a wired connection instead of Wi-Fi

3. **Subscription Import Failures**:
   - Verify the subscription URL
   - Check if the subscription format is supported
   - Try accessing the URL in a web browser to ensure it's accessible

4. **API Access Issues**:
   - Verify the API server is running
   - Check firewall settings
   - Ensure the correct port is being used

### Logs

The VPN Client writes logs to the `logs/` directory. Check these logs for detailed error information.

### Getting Help

If you encounter issues not covered in this guide:

1. Check the application logs
2. Refer to the API documentation
3. Report issues on the project's GitHub page
4. Contact support (if available)