package protocols

import (
	"sync"
	"time"
)

// ConnectionStats represents connection statistics
type ConnectionStats struct {
	BytesSent     int64
	BytesReceived int64
	Duration      int64 // in seconds
}

// ServerConfig represents a VPN server configuration
type ServerConfig struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	UUID     string `json:"uuid"`
	Security string `json:"security"`
}

// ProtocolHandler defines the interface for all protocol handlers
type ProtocolHandler interface {
	// Connect establishes a connection to the VPN server
	Connect(config *ServerConfig) error

	// Disconnect terminates the VPN connection
	Disconnect() error

	// IsConnected returns true if the VPN is currently connected
	IsConnected() bool

	// GetStats returns connection statistics
	GetStats() *ConnectionStats
}

// BaseHandler provides common functionality for protocol handlers
type BaseHandler struct {
	connected    bool
	dataSent     int64
	dataReceived int64
	lastUpdate   time.Time
	mutex        sync.RWMutex
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

// Initialize registers all available protocols
func Initialize() {
	RegisterProtocol("vless", func() ProtocolHandler {
		return &VLESSHandler{}
	})

	// Add other protocols here as they are implemented
	// RegisterProtocol("vmess", func() ProtocolHandler { return NewVMessHandler() })
	// RegisterProtocol("trojan", func() ProtocolHandler { return NewTrojanHandler() })
	// etc.
}

// UnknownProtocolError is returned when trying to create an unregistered protocol
type UnknownProtocolError struct {
	Protocol string
}

func (e *UnknownProtocolError) Error() string {
	return "unknown protocol: " + e.Protocol
}
