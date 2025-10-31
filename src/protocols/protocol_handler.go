package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"errors"
	"time"
)

// BaseHandler provides common functionality for protocol handlers
type BaseHandler struct {
	protocol   core.ProtocolType
	connected  bool
	dataSent   int64
	dataReceived int64
	lastUpdate time.Time
	mutex      sync.RWMutex
}

// ConnectionStats represents connection statistics
type ConnectionStats struct {
	BytesSent     int64
	BytesReceived int64
	Duration      int64 // in seconds
}

// ProtocolHandler defines the interface for all protocol handlers
type ProtocolHandler interface {
	// Connect establishes a connection to the VPN server
	Connect(config *core.ServerConfig) error
	
	// Disconnect terminates the VPN connection
	Disconnect() error
	
	// IsConnected returns true if the VPN is currently connected
	IsConnected() bool
	
	// GetStats returns connection statistics
	GetStats() *ConnectionStats
}

// ProtocolFactory is a function that creates a new protocol handler
type ProtocolFactory func() ProtocolHandler

// protocolFactories stores registered protocol factories
var protocolFactories = make(map[string]ProtocolFactory)

// RegisterProtocol registers a new protocol factory
func RegisterProtocol(name string, factory ProtocolFactory) {
	protocolFactories[name] = factory
}

// CreateProtocol creates a new protocol handler by name
func CreateProtocol(name string) (ProtocolHandler, error) {
	factory, exists := protocolFactories[name]
	if !exists {
		return nil, &UnknownProtocolError{name}
	}
	return factory(), nil
}

// UnknownProtocolError is returned when trying to create an unregistered protocol
type UnknownProtocolError struct {
	Protocol string
}

func (e *UnknownProtocolError) Error() string {
	return "unknown protocol: " + e.Protocol
}

