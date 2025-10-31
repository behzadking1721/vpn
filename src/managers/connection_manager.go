package managers

import (
	"context"
	"sync"
	"time"
	"c:/Users/behza/OneDrive/Documents/vpn/src/protocols"
)

// ConnectionManager handles VPN connections with protocol-agnostic operations
type ConnectionManager struct {
	status       ConnectionStatus
	connInfo     ConnectionInfo
	mutex        sync.RWMutex
	cancelFunc   context.CancelFunc
	ctx          context.Context
	handler      protocols.ProtocolHandler // Current protocol-specific handler
	isConnected  bool                      // Simplified connection state
}

// ConnectionListener interface for connection status updates
type ConnectionListener interface {
	OnStatusChanged(status ConnectionStatus, info ConnectionInfo)
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	ctx, cancel := context.WithCancel(context.Background())

	cm := &ConnectionManager{
		status:      StatusDisconnected,
		connInfo:    ConnectionInfo{Status: StatusDisconnected},
		ctx:         ctx,
		cancelFunc:  cancel,
		isConnected: false,
	}

	return cm
}

// AddListener adds a connection status listener
func (cm *ConnectionManager) AddListener(listener ConnectionListener) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	// In a real implementation, we would store listeners
	// For now, this is just a placeholder
}


// Connect establishes a connection to a VPN server
func (cm *ConnectionManager) Connect(config *protocols.ServerConfig) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.isConnected {
		return nil // Already connected
	}

	// Create protocol handler
	handler, err := protocols.CreateProtocol(config.Protocol)
	if err != nil {
		// If protocol is not registered, create a mock handler for demonstration
		handler = &MockProtocolHandler{}
	}

	// Attempt to connect
	err = handler.Connect(config)
	if err != nil {
		return err
	}

	// Update connection state
	cm.handler = handler
	cm.isConnected = true
	cm.status = StatusConnected

	// Notify listeners
	cm.connInfo = ConnectionInfo{
		Status:     StatusConnected,
		ServerName: config.Name,
	}

	return nil
}


// Disconnect terminates the current VPN connection
func (cm *ConnectionManager) Disconnect() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if !cm.isConnected {
		return nil // Already disconnected
	}

	cm.status = StatusDisconnecting

	var err error
	if cm.handler != nil {
		err = cm.handler.Disconnect()
	}

	cm.isConnected = false
	cm.status = StatusDisconnected
	cm.handler = nil

	// Reset connection info
	cm.connInfo = ConnectionInfo{
		Status: StatusDisconnected,
	}

	return err
}


// IsConnected returns true if currently connected to a VPN
func (cm *ConnectionManager) IsConnected() bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.isConnected
}


// GetConnectionInfo returns current connection information
func (cm *ConnectionManager) GetConnectionInfo() ConnectionInfo {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.connInfo
}

