package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"context"
	"errors"
	"sync"
	"time"
)

// ConnectionManager handles VPN connection lifecycle
type ConnectionManager struct {
	status       core.ConnectionStatus
	currentServer *core.Server
	connInfo     core.ConnectionInfo
	mutex        sync.RWMutex
	cancelFunc   context.CancelFunc
	ctx          context.Context
	listeners    []ConnectionListener
}

// ConnectionListener interface for connection status updates
type ConnectionListener interface {
	OnStatusChanged(status core.ConnectionStatus, info core.ConnectionInfo)
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	cm := &ConnectionManager{
		status:    core.StatusDisconnected,
		connInfo:  core.ConnectionInfo{Status: core.StatusDisconnected},
		listeners: make([]ConnectionListener, 0),
	}

	return cm
}

// AddListener adds a connection status listener
func (cm *ConnectionManager) AddListener(listener ConnectionListener) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.listeners = append(cm.listeners, listener)
}

// RemoveListener removes a connection status listener
func (cm *ConnectionManager) RemoveListener(listener ConnectionListener) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for i, l := range cm.listeners {
		if l == listener {
			cm.listeners = append(cm.listeners[:i], cm.listeners[i+1:]...)
			break
		}
	}
}

// notifyListeners notifies all listeners of status changes
func (cm *ConnectionManager) notifyListeners() {
	cm.mutex.RLock()
	listeners := make([]ConnectionListener, len(cm.listeners))
	copy(listeners, cm.listeners)
	connInfo := cm.connInfo
	cm.mutex.RUnlock()

	for _, listener := range listeners {
		listener.OnStatusChanged(cm.status, connInfo)
	}
}

// Connect connects to a specified server
func (cm *ConnectionManager) Connect(server core.Server) error {
	cm.mutex.Lock()

	// Check if already connected or connecting
	if cm.status == core.StatusConnected || cm.status == core.StatusConnecting {
		cm.mutex.Unlock()
		return errors.New("already connected or connecting")
	}

	// Update status
	cm.status = core.StatusConnecting
	cm.currentServer = &server
	cm.connInfo = core.ConnectionInfo{
		Status:    core.StatusConnecting,
		ServerID:  server.ID,
		StartTime: time.Now(),
	}

	// Notify listeners
	cm.notifyListeners()
	cm.mutex.Unlock()

	// In a real implementation, this would establish the actual VPN connection
	// based on the protocol type and server configuration
	// For now, we'll simulate the connection process

	// Simulate connection delay
	time.Sleep(1 * time.Second)

	// Update status to connected
	cm.mutex.Lock()
	cm.status = core.StatusConnected
	cm.connInfo.Status = core.StatusConnected
	cm.connInfo.LocalIP = "10.0.0.2"
	cm.connInfo.RemoteIP = server.Host
	cm.mutex.Unlock()

	// Notify listeners
	cm.notifyListeners()

	return nil
}

// Disconnect disconnects from the current server
func (cm *ConnectionManager) Disconnect() error {
	cm.mutex.Lock()

	// Check if already disconnected
	if cm.status == core.StatusDisconnected || cm.status == core.StatusDisconnecting {
		cm.mutex.Unlock()
		return errors.New("not connected")
	}

	// Update status
	cm.status = core.StatusDisconnecting
	cm.connInfo.Status = core.StatusDisconnecting

	// Notify listeners
	cm.notifyListeners()
	cm.mutex.Unlock()

	// In a real implementation, this would tear down the VPN connection
	// For now, we'll simulate the disconnection process

	// Simulate disconnection delay
	time.Sleep(500 * time.Millisecond)

	// Update status to disconnected
	cm.mutex.Lock()
	cm.status = core.StatusDisconnected
	cm.connInfo.Status = core.StatusDisconnected
	cm.connInfo.EndTime = time.Now()
	cm.currentServer = nil
	cm.mutex.Unlock()

	// Notify listeners
	cm.notifyListeners()

	return nil
}

// GetStatus returns the current connection status
func (cm *ConnectionManager) GetStatus() core.ConnectionStatus {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.status
}

// GetCurrentServer returns the currently connected server
func (cm *ConnectionManager) GetCurrentServer() *core.Server {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	// Return a copy to prevent external modification
	if cm.currentServer != nil {
		server := *cm.currentServer
		return &server
	}
	
	return nil
}

// GetConnectionInfo returns detailed connection information
func (cm *ConnectionManager) GetConnectionInfo() core.ConnectionInfo {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	// Return a copy
	info := cm.connInfo
	return info
}

// UpdateDataUsage updates the data sent/received statistics
func (cm *ConnectionManager) UpdateDataUsage(sent, received int64) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	cm.connInfo.DataSent = sent
	cm.connInfo.DataReceived = received
	
	// Notify listeners of data usage update
	cm.notifyListeners()
}