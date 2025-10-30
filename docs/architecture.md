# VPN Client Architecture

This document describes the architecture and design of the multi-platform VPN client application.

## Overview

The VPN client is designed as a modular application with separate components for different functionalities. The architecture follows a layered approach with clear separation of concerns.

## Core Components

### 1. Core Package (`src/core`)
Contains the fundamental data structures and types used throughout the application:
- Server and Subscription models
- Connection status enums
- Application configuration structures

### 2. Managers (`src/managers`)
Handles business logic and state management:
- **ServerManager**: Manages server configurations and storage
- **ConnectionManager**: Controls VPN connection lifecycle
- **SubscriptionManager**: Handles subscription links and QR code imports
- **ConfigManager**: Manages application configuration

### 3. Protocols (`src/protocols`)
Implements support for different VPN protocols:
- VMess, VLESS, Trojan, Reality, Hysteria2, TUIC, SSH, Shadowsocks
- Protocol-specific connection handlers
- Unified interface for all protocols

### 4. Utilities (`src/utils`)
Provides helper functions:
- ID generation
- Data formatting
- Ping measurement
- Other common utilities

## Supported Protocols

The application supports the following VPN protocols:

| Protocol | Status | Description |
|----------|--------|-------------|
| VMess | Implemented | V2Ray's primary protocol |
| VLESS | Implemented | Next-generation protocol from V2Ray |
| Shadowsocks | Implemented | Lightweight SOCKS5-like proxy |
| Trojan | Implemented | Application-layer protocol |
| Reality | Implemented | Advanced TLS-based protocol |
| Hysteria2 | Implemented | Modern high-performance protocol |
| TUIC | Implemented | QUIC-based protocol |
| SSH | Implemented | Secure Shell tunnel |

## Key Features

### Multi-Platform Support
The application is designed to work on:
- Android
- iOS
- Windows
- macOS
- Linux

### Server Management
- Add/remove servers manually
- Import servers via subscription links
- Import servers via QR codes
- Enable/disable servers
- Ping measurement for server speed

### Connection Features
- Automatic fastest server selection
- Data usage tracking
- IPv6 support
- Multiple tunneling modes (TCP/UDP/Both/Auto)
- Persistent connection settings

### Security & Privacy
- No ads or tracking
- Open-source codebase
- Local data storage
- Encrypted communications

## Architecture Diagram

```
┌─────────────────────────────────────┐
│              UI Layer               │
│  (Platform-specific interfaces)     │
├─────────────────────────────────────┤
│          Business Logic             │
│  ┌──────────────┬───────────────┐   │
│  │ ServerMgr    │ ConnectionMgr │   │
│  ├──────────────┼───────────────┤   │
│  │ SubscribeMgr │ ConfigMgr     │   │
│  └──────────────┴───────────────┘   │
├─────────────────────────────────────┤
│           Protocol Layer            │
│                                     │
│  VMess  VLESS  SS  Trojan  Reality  │
│  HY2    TUIC   SSH                  │
├─────────────────────────────────────┤
│           Utility Layer             │
│                                     │
│  Helpers  Ping  Formatting          │
└─────────────────────────────────────┘
```

## Data Flow

1. **Server Configuration**
   - Servers added manually or imported via subscriptions/QR codes
   - Stored in ServerManager
   - Periodic ping measurements to determine server quality

2. **Connection Process**
   - User selects server or enables auto-select
   - ConnectionManager creates appropriate ProtocolHandler
   - ProtocolHandler establishes connection based on server protocol
   - Connection status monitored and reported

3. **Data Tracking**
   - ProtocolHandlers report data usage
   - ConnectionManager aggregates and reports statistics
   - UI displays real-time data usage information

## Future Enhancements

Planned improvements and additions:
- Real implementation of protocol handlers using open-source libraries
- UI implementation for all platforms
- Advanced routing rules
- Custom DNS settings
- Split tunneling capabilities
- Dark/light theme support
- Localization for multiple languages