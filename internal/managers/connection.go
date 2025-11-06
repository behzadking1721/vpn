package managers

import (
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

// String returns a string representation of the ConnectionStatus
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
	status       ConnectionStatus
	currentServer *core.Server
	startTime    time.Time
	dataSent     int64
	dataReceived int64
	mutex        sync.RWMutex
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
	status := cm.GetStatus()
	switch status {
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

// GetCurrentServer returns the currently connected server
func (cm *ConnectionManager) GetCurrentServer() *core.Server {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.currentServer
}

// GetDataUsage returns current data usage
func (cm *ConnectionManager) GetDataUsage() (int64, int64) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	// Simulate data usage increase for demo purposes
	if cm.status == Connected {
		// In a real implementation, this would come from the actual connection
		// For demo, we'll just simulate increasing data usage
		duration := time.Since(cm.startTime).Seconds()
		sent := int64(duration * 1024 * 1024) // 1MB/s
		received := int64(duration * 2 * 1024 * 1024) // 2MB/s
		
		return sent, received
	}
	
	return cm.dataSent, cm.dataReceived
}

// GetUptime returns connection uptime in seconds
func (cm *ConnectionManager) GetUptime() time.Duration {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	if cm.status == Connected {
		return time.Since(cm.startTime)
	}
	
	return 0
}

// Connect connects to a VPN server
func (cm *ConnectionManager) Connect(server *core.Server) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.status == Connected {
		return fmt.Errorf("already connected to a server")
	}

	cm.status = Connecting
	cm.currentServer = server
	cm.startTime = time.Now()
	cm.dataSent = 0
	cm.dataReceived = 0

	// Check if this is a WireGuard server
	if server != nil && server.Protocol == "wireguard" {
		err := cm.connectWireGuard(server)
		if err != nil {
			cm.status = Error
			return err
		}
	} else {
		// Simulate connection process
		// In a real implementation, this would involve:
		// 1. Initializing the appropriate protocol handler
		// 2. Establishing the connection
		// 3. Setting up routing rules
		// 4. Starting data transfer monitoring

		// For demo purposes, we'll just simulate a successful connection
		time.Sleep(100 * time.Millisecond)
	}

	cm.status = Connected

	return nil
}

// Disconnect disconnects from the current VPN server
func (cm *ConnectionManager) Disconnect() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.status != Connected {
		return fmt.Errorf("not connected to any server")
	}

	cm.status = Disconnecting

	// Check if this is a WireGuard connection
	if cm.currentServer != nil && cm.currentServer.Protocol == "wireguard" {
		err := cm.disconnectWireGuard(cm.currentServer)
		if err != nil {
			cm.status = Error
			return err
		}
	} else {
		// Simulate disconnection process
		// In a real implementation, this would involve:
		// 1. Tearing down the connection
		// 2. Cleaning up routing rules
		// 3. Stopping data transfer monitoring

		// For demo purposes, we'll just simulate a successful disconnection
		time.Sleep(100 * time.Millisecond)
	}
	
	cm.status = Disconnected
	cm.currentServer = nil
	cm.startTime = time.Time{}

	return nil
}

// connectWireGuard connects to a WireGuard server
func (cm *ConnectionManager) connectWireGuard(server *core.Server) error {
	// This will be implemented in wireguard_wgctrl.go with build tag
	return fmt.Errorf("wireguard support not compiled in this build")
}

// disconnectWireGuard disconnects from a WireGuard server
func (cm *ConnectionManager) disconnectWireGuard(server *core.Server) error {
	// This will be implemented in wireguard_wgctrl.go with build tag
	return fmt.Errorf("wireguard support not compiled in this build")
}
