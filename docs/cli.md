# CLI Tool Documentation (vpnctl)

## Overview

The `vpnctl` command-line interface provides a powerful way to control and manage your VPN client. It offers full programmatic access to all VPN client features, making it ideal for automation, scripting, and advanced users who prefer the command line.

## Installation

To build the CLI tool:

```bash
cd cmd/vpnctl
go build -o vpnctl
```

Optionally, you can install it system-wide:

```bash
go install
```

## Command Structure

The CLI follows a consistent command structure:

```bash
vpnctl [command] [subcommand] [flags] [arguments]
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | Output in JSON format |
| `--verbose` | `-v` | Verbose output |
| `--quiet` | `-q` | Quiet mode |
| `--config` | | Configuration file path |
| `--timeout` | | Timeout in seconds |
| `--help` | `-h` | Help for command |
| `--version` | | Version information |

## Core Commands

### Connection Management

#### connect
Connect to a VPN server.

```bash
# Connect to a specific server
vpnctl connect server-name

# Connect to the best server
vpnctl connect --best

# Connect to the fastest server
vpnctl connect --fastest
```

#### disconnect
Disconnect from the current VPN server.

```bash
vpnctl disconnect
```

#### status
Show current VPN connection status.

```bash
# Show status in human-readable format
vpnctl status

# Show status in JSON format
vpnctl status --json
```

#### restart
Restart the VPN service.

```bash
vpnctl restart
```

### Server Management

#### server list
List all available VPN servers.

```bash
vpnctl server list
```

#### server add
Add a VPN subscription from URL.

```bash
vpnctl server add https://example.com/subscription
```

#### server remove
Remove a VPN server.

```bash
vpnctl server remove server-name
```

#### server test
Test speed of all servers.

```bash
vpnctl server test
```

#### server update
Update server list from subscriptions.

```bash
vpnctl server update
```

#### server favorites
Manage favorite servers.

```bash
vpnctl server favorites
```

### Statistics and Monitoring

#### stats
Show data usage statistics.

```bash
# Show statistics in human-readable format
vpnctl stats

# Show statistics in JSON format
vpnctl stats --json
```

#### logs
View application logs.

```bash
# View recent logs
vpnctl logs

# Follow logs in real-time
vpnctl logs --follow
```

#### speedtest
Perform a network speed test.

```bash
vpnctl speedtest
```

#### traffic
Show real-time traffic statistics.

```bash
vpnctl traffic
```

### Configuration Management

#### config show
Show current configuration.

```bash
vpnctl config show
```

#### config set
Set a configuration value.

```bash
vpnctl config set key value
```

#### config reset
Reset configuration to defaults.

```bash
vpnctl config reset
```

#### config import
Import configuration from file.

```bash
vpnctl config import config.json
```

#### config export
Export configuration to file.

```bash
vpnctl config export config.json
```

## Usage Examples

### Basic Connection Workflow

```bash
# List available servers
vpnctl server list

# Connect to the best server
vpnctl connect --best

# Check connection status
vpnctl status

# View real-time traffic
vpnctl traffic

# Disconnect
vpnctl disconnect
```

### Automation Script Example

```bash
#!/bin/bash

# Connect to VPN
vpnctl connect --best

# Perform speed test
vpnctl speedtest

# Run some network tasks
# ... your commands here ...

# Disconnect
vpnctl disconnect
```

### JSON Output for Scripting

```bash
# Get connection status as JSON
vpnctl status --json | jq '.connected'

# Get server list as JSON
vpnctl server list --json | jq '.[] | select(.status=="online")'
```

## Output Formats

### Human-Readable Format
Commands output formatted, colored text that's easy for humans to read:

```
🟢 Connected to: USA - New York
⏱️  Duration: 12:34:56
📊 Download: 125 MB ↑ | Upload: 45 MB ↓
🌐 IP: 192.168.1.100
```

### JSON Format
Use the `--json` flag for machine-readable JSON output:

```json
{
  "connected": true,
  "server": "USA - New York",
  "duration": "12:34:56",
  "upload": "125 MB",
  "download": "45 MB",
  "ip": "192.168.1.100"
}
```

### Table Format
Some commands display data in tables for easy scanning:

```
┌──────┬────────────┬────────┬─────────┬──────────┐
│ Flag │ Name       │ Ping   │ Speed   │ Status   │
├──────┼────────────┼────────┼─────────┼──────────┤
│ 🇺🇸  │ USA-NY     │ 45ms   │ 85 Mbps │ 🟢 Online │
│ 🇩🇪  │ Germany    │ 65ms   │ 72 Mbps │ 🟢 Online │
│ 🇯🇵  │ Japan      │ 120ms  │ 45 Mbps │ 🟡 Slow   │
└──────┴────────────┴────────┴─────────┴──────────┘
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Misuse of shell builtins |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `VPNCTL_CONFIG` | Path to configuration file |
| `VPNCTL_TIMEOUT` | Default timeout in seconds |

## Troubleshooting

### Common Issues

1. **Permission denied**: Make sure the CLI tool has appropriate permissions
2. **Connection failed**: Check if the VPN service is running
3. **Server not found**: Verify server name or update server list

### Getting Help

All commands support the `--help` flag for detailed usage information:

```bash
vpnctl --help
vpnctl connect --help
vpnctl server --help
```

## Advanced Features

### Shell Completions

Generate shell completion scripts:

```bash
# Bash
source <(vpnctl completion bash)

# Zsh
source <(vpnctl completion zsh)

# Fish
vpnctl completion fish | source
```

### Configuration File

The CLI tool can use a configuration file (default: `$HOME/.vpnctl.yaml`):

```yaml
server: "default-server"
auto_connect: true
log_level: "info"
timeout: 30
```

## API Integration

The CLI tool communicates with the VPN client backend service through its REST API. This architecture allows:

1. Multiple CLI instances to control the same backend
2. Consistent behavior across all interfaces (CLI, Desktop, Mobile)
3. Remote management capabilities

## Performance Considerations

- CLI commands are lightweight and execute quickly
- Network operations (connect, speedtest) may take longer
- Real-time monitoring commands update at regular intervals
- JSON output is optimized for parsing by other tools

## Security

- All communication with the backend service is through localhost
- Configuration files should have appropriate permissions
- Sensitive information is handled by the backend service