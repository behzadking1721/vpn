# Protocol Integration Guide

This document explains how to integrate real protocol libraries into the VPN client application.

## Overview

The VPN client uses a modular architecture that allows for easy integration of different protocol libraries. Each protocol is implemented as a handler that conforms to the [ProtocolHandler](file:///c:/Users/behza/OneDrive/Documents/vpn/src/protocols/protocol_handler.go#L9-L16) interface.

## Adding a New Protocol

To add support for a new protocol, follow these steps:

### 1. Update Core Types

Add the new protocol to the [ProtocolType](file:///c:/Users/behza/OneDrive/Documents/vpn/src/core/types.go#L9-L16) enum in [src/core/types.go](file:///c:/Users/behza/OneDrive/Documents/vpn/src/core/types.go):

```go
const (
    // ... existing protocols ...
    ProtocolNewProtocol ProtocolType = "newprotocol"
)
```

### 2. Update Protocol Factory

Add the new protocol to the [ProtocolFactory](file:///c:/Users/behza/OneDrive/Documents/vpn/src/protocols/protocol_handler.go#L30-L32) in [src/protocols/protocol_handler.go](file:///c:/Users/behza/OneDrive/Documents/vpn/src/protocols/protocol_handler.go):

```go
func (pf *ProtocolFactory) CreateHandler(protocolType core.ProtocolType) (ProtocolHandler, error) {
    switch protocolType {
    // ... existing cases ...
    case core.ProtocolNewProtocol:
        return NewNewProtocolHandler(), nil
    default:
        return nil, errors.New("unsupported protocol type")
    }
}
```

### 3. Create Protocol Handler

Create a new file `src/protocols/newprotocol_handler.go` with the implementation:

```go
package protocols

import (
    "c:/Users/behza/OneDrive/Documents/vpn/src/core"
    // Import the actual protocol library
)

type NewProtocolHandler struct {
    BaseHandler
    server core.Server
    // Add protocol-specific fields
}

func NewNewProtocolHandler() *NewProtocolHandler {
    handler := &NewProtocolHandler{}
    handler.BaseHandler.protocol = core.ProtocolNewProtocol
    return handler
}

func (nph *NewProtocolHandler) Connect(server core.Server) error {
    // Implementation here
    return nil
}

func (nph *NewProtocolHandler) Disconnect() error {
    // Implementation here
    return nil
}

func (nph *NewProtocolHandler) GetDataUsage() (sent, received int64, err error) {
    // Implementation here
    return 0, 0, nil
}

func (nph *NewProtocolHandler) GetConnectionDetails() (map[string]interface{}, error) {
    // Implementation here
    return nil, nil
}
```

## Protocol Integration Examples

### Shadowsocks Integration

The Shadowsocks integration demonstrates how to use an external library:

1. Add dependency to [go.mod](file:///c:/Users/behza/OneDrive/Documents/vpn/go.mod):
```go
require (
    github.com/shadowsocks/go-shadowsocks2 v0.0.0-20230516033142-213602970b87
)
```

2. Import and use in handler:
```go
import (
    "github.com/shadowsocks/go-shadowsocks2/core"
    "github.com/shadowsocks/go-shadowsocks2/shadowsocks"
)
```

3. Initialize cipher and client:
```go
cipher, err := core.PickCipher(method, []byte(password))
if err != nil {
    return err
}

client, err := shadowsocks.NewClient(host, cipher)
if err != nil {
    return err
}
```

### V2Ray/XRay Integration

For VMess/VLESS protocols:

1. Add dependency:
```go
require (
    github.com/v2ray/v2ray-core v0.0.0-20230601040104-11eba6094346
    // or for XRay:
    // github.com/xtls/xray-core v0.0.0-20230601040104-11eba6094346
)
```

2. Import and use:
```go
import (
    "github.com/v2ray/v2ray-core/app/proxyman"
    "github.com/v2ray/v2ray-core/common/protocol"
    "github.com/v2ray/v2ray-core/common/serial"
    "github.com/v2ray/v2ray-core/proxy/vmess"
    "github.com/v2ray/v2ray-core/proxy/vmess/outbound"
)
```

3. Configure and start:
```go
config := &core.Config{
    // Configuration details
}

server, err := core.New(config)
if err != nil {
    return err
}

err = server.Start()
if err != nil {
    return err
}
```

## Best Practices

### Error Handling

Always provide meaningful error messages:
```go
if server.Method == "" {
    return fmt.Errorf("missing encryption method for Shadowsocks")
}
```

### Resource Management

Ensure proper cleanup of resources:
```go
func (handler *ProtocolHandler) Disconnect() error {
    if handler.client != nil {
        handler.client.Close()
        handler.client = nil
    }
    return nil
}
```

### Data Tracking

Implement data usage tracking when possible:
```go
func (handler *ProtocolHandler) GetDataUsage() (sent, received int64, err error) {
    if handler.client == nil {
        return 0, 0, fmt.Errorf("not connected")
    }
    
    sent = handler.client.GetBytesSent()
    received = handler.client.GetBytesReceived()
    return sent, received, nil
}
```

## Testing Protocol Handlers

Create tests for each protocol handler in `protocols_test.go`:

```go
func TestNewProtocolHandler(t *testing.T) {
    handler := NewNewProtocolHandler()
    
    // Test initial state
    if handler.IsConnected() {
        t.Error("Expected handler to be disconnected initially")
    }
    
    // Test connection
    server := core.Server{
        // Server configuration
    }
    
    err := handler.Connect(server)
    if err != nil {
        t.Errorf("Failed to connect: %v", err)
    }
    
    // Test disconnection
    err = handler.Disconnect()
    if err != nil {
        t.Errorf("Failed to disconnect: %v", err)
    }
}
```

## Security Considerations

1. Validate all input parameters
2. Handle sensitive data (passwords, keys) securely
3. Use secure defaults for encryption methods
4. Implement proper error handling without leaking sensitive information

## Performance Optimization

1. Reuse connections when possible
2. Implement connection pooling
3. Use buffered I/O operations
4. Minimize memory allocations

## Troubleshooting

Common issues and solutions:

1. **Dependency conflicts**: Use Go modules to manage versions
2. **Compilation errors**: Check library documentation for required CGO settings
3. **Connection failures**: Verify server configuration and network access
4. **Performance issues**: Profile the application to identify bottlenecks

## Next Steps

1. Implement real protocol libraries for all supported protocols
2. Add comprehensive unit tests
3. Create integration tests with real servers
4. Document protocol-specific configuration options
5. Implement protocol-specific error handling