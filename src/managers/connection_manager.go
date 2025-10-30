package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/protocols"
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
	handler      protocols.ProtocolHandler // Handler for the current connection
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

	// Create protocol handler
	factory := protocols.NewProtocolFactory()
	handler, err := factory.CreateHandler(server.Protocol)
	if err != nil {
		cm.mutex.Lock()
		cm.status = core.StatusError
		cm.connInfo.Status = core.StatusError
		cm.connInfo.Error = "Unsupported protocol: " + string(server.Protocol)
		cm.notifyListeners()
		cm.mutex.Unlock()
		return err
	}

	// Store handler for later use
	cm.handler = handler

	// Try to connect using the protocol handler
	err = handler.Connect(server)
	if err != nil {
		cm.mutex.Lock()
		cm.status = core.StatusError
		cm.connInfo.Status = core.StatusError
		cm.connInfo.Error = "Connection failed: " + err.Error()
		cm.notifyListeners()
		cm.mutex.Unlock()
		return err
	}

	// Update status to connected
	cm.mutex.Lock()
	cm.status = core.StatusConnected
	cm.connInfo.Status = core.StatusConnected
	cm.connInfo.LocalIP = "10.0.0.2"
	cm.connInfo.RemoteIP = server.Host
	cm.mutex.Unlock()

	// Start monitoring data usage
	go cm.monitorDataUsage()

	// Notify listeners
	cm.notifyListeners()

	return nil
}

// monitorDataUsage periodically updates data usage statistics
func (cm *ConnectionManager) monitorDataUsage() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.mutex.RLock()
			isConnected := cm.status == core.StatusConnected
			handler := cm.handler
			cm.mutex.RUnlock()

			if !isConnected || handler == nil {
				return
			}

			sent, received, err := handler.GetDataUsage()
			if err == nil {
				cm.UpdateDataUsage(sent, received)
			}
		case <-cm.ctx.Done():
			return
		}
	}
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

	// Disconnect using the protocol handler
	if cm.handler != nil {
		cm.handler.Disconnect()
		cm.handler = nil
	}

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