# Command Line Interface (CLI) Guide

This document explains how to use the VPN Client's command-line interface.

## Starting the CLI

To start the CLI interface, run:

```bash
cd src
go run main.go --cli
```

## CLI Menu Options

Once the CLI is running, you'll see the main menu with the following options:

### 1. List servers
Displays all configured servers with their details including:
- Server name
- Host and port
- Protocol type
- Enabled/disabled status
- Ping time (if available)

### 2. Add server
Allows you to add a new server configuration:
1. Enter server name
2. Enter host (domain or IP address)
3. Enter port number
4. Select protocol from the list:
   - VMess
   - VLESS
   - Shadowsocks
   - Trojan
   - Reality
   - Hysteria2
   - TUIC
   - SSH
5. Enter protocol-specific settings (method/password for Shadowsocks, encryption/password for others)
6. The server is automatically enabled

### 3. Connect to server
Connects to a configured server:
1. Shows the list of configured servers
2. Prompts for server selection by number
3. Attempts to establish connection to the selected server

### 4. Disconnect
Disconnects from the currently connected server.

### 5. Show status
Displays the current connection status including:
- Connection status (connected/disconnected/connecting)
- Connection start time
- Data sent/received (if available)
- Connected server details

### 6. Ping servers
Performs ping tests on all enabled servers:
- Sends ping requests to each server
- Updates server ping times
- Displays results for each server

### 7. Find fastest server
Automatically finds the server with the lowest ping time:
- Checks all enabled servers
- Displays the fastest server and its ping time

### 8. Exit
Exits the CLI application.

## Example Usage

Here's a typical workflow using the CLI:

1. Start the CLI: `go run main.go --cli`
2. Add a server using option 2
3. List servers using option 1 to verify the server was added
4. Connect to the server using option 3
5. Check connection status using option 5
6. Disconnect using option 4 when finished
7. Exit the CLI using option 8

## Server Configuration

When adding servers, the following information is required:

- **Name**: A descriptive name for the server
- **Host**: The server's domain name or IP address
- **Port**: The port number the server listens on
- **Protocol**: The VPN protocol to use
- **Protocol-specific settings**:
  - For Shadowsocks: Encryption method and password
  - For other protocols: Encryption type and password/UUID

## Data Management

The CLI uses the same underlying managers as the API server:
- Server configurations are managed by ServerManager
- Connection state is managed by ConnectionManager
- Configuration is managed by ConfigManager

All data is persisted to disk and will be available in subsequent sessions.

## Error Handling

The CLI provides informative error messages for common issues:
- Invalid server selection
- Connection failures
- Missing server configurations
- Invalid input values

Error messages are displayed in plain text and the CLI returns to the main menu after an error, allowing you to retry or choose a different option.