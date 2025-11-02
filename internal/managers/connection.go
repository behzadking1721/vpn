package managers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vpnclient/src/core"
)

// ConnectionStatus represents the connection status
type ConnectionStatus int

const (
	// Disconnected represents a disconnected state
	Disconnected ConnectionStatus = iota
	// Connecting represents a connecting state
	Connecting
	// Connected represents a connected state
	Connected
	// Disconnecting represents a disconnecting state
	Disconnecting
	// Error represents an error state
	Error
)

// String converts ConnectionStatus to string
func (s ConnectionStatus) String() string {
	switch s {
	case Disconnected:
		return "Disconnected"
	case Connecting:
		return "Connecting"
	case Connected:
		return "Connected"
	case Disconnecting:
		return "Disconnecting"
	case Error:
		return "Error"
	default:
		return "Unknown"
	}
}

// ConnectionManager manages VPN connections
type ConnectionManager struct {
	status      ConnectionStatus
	currentServer *core.Server
	connectedAt   time.Time
	dataSent      int64
	dataReceived  int64
	mutex         sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		status: Disconnected,
	}
}

// GetStatus returns the current connection status
func (cm *ConnectionManager) GetStatus() ConnectionStatus {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.status
}

// GetStatusString returns the current connection status as a string
func (cm *ConnectionManager) GetStatusString() string {
	return cm.GetStatus().String()
}

// GetCurrentServer returns the currently connected server
func (cm *ConnectionManager) GetCurrentServer() *core.Server {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.currentServer
}

// GetConnectionInfo returns connection information
func (cm *ConnectionManager) GetConnectionInfo() *core.ConnectionInfo {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	return &core.ConnectionInfo{
		ID:        "",
		StartedAt: cm.connectedAt,
		DataSent:  cm.dataSent,
		DataRecv:  cm.dataReceived,
	}
}

// UpdateDataUsage updates data usage statistics
func (cm *ConnectionManager) UpdateDataUsage(sent, received int64) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	cm.dataSent += sent
	cm.dataReceived += received
}

// GetDataUsage returns current data usage
func (cm *ConnectionManager) GetDataUsage() (sent, received int64) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.dataSent, cm.dataReceived
}

// GetUptime returns connection uptime in seconds
func (cm *ConnectionManager) GetUptime() int64 {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	if cm.status != Connected || cm.connectedAt.IsZero() {
		return 0
	}
	
	return int64(time.Since(cm.connectedAt).Seconds())
}

// Connect attempts to connect to a server
func (cm *ConnectionManager) Connect(config interface{}) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	// If already connected, disconnect first
	if cm.status == Connected {
		cm.status = Disconnecting
		cm.status = Disconnected
		cm.currentServer = nil
		cm.dataSent = 0
		cm.dataReceived = 0
	}
	
	// Parse server from config
	var server *core.Server
	if s, ok := config.(*core.Server); ok {
		server = s
	} else {
		return fmt.Errorf("invalid server configuration")
	}
	
	cm.status = Connecting
	cm.currentServer = server
	
	// Simulate connection process
	time.Sleep(100 * time.Millisecond)
	
	cm.status = Connected
	cm.connectedAt = time.Now()
	cm.dataSent = 0
	cm.dataReceived = 0
	
	// Create context for connection lifecycle
	cm.ctx, cm.cancel = context.WithCancel(context.Background())
	
	return nil
}

// Disconnect disconnects from the server
func (cm *ConnectionManager) Disconnect() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if cm.status != Connected {
		return fmt.Errorf("not connected")
	}
	
	cm.status = Disconnecting
	
	// Cancel context if exists
	if cm.cancel != nil {
		cm.cancel()
	}
	
	// Simulate disconnection process
	time.Sleep(100 * time.Millisecond)
	
	cm.status = Disconnected
	cm.currentServer = nil
	cm.connectedAt = time.Time{}
	
	return nil
}