package managers

import (
	"fmt"
	"sync"
	"time"

	"vpnclient/internal/logging"
	"vpnclient/internal/notifications"
	"vpnclient/internal/stats"
	"vpnclient/src/core"
)

// ConnectionStatus represents the status of a VPN connection
type ConnectionStatus int

const (
	// Disconnected means no connection is active
	Disconnected ConnectionStatus = iota
	// Connecting means a connection is being established
	Connecting
	// Connected means a connection is active
	Connected
	// Disconnecting means a connection is being terminated
	Disconnecting
	// Error means there was an error with the connection
	Error
)

// ConnectionInfo holds information about the current connection
type ConnectionInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartedAt time.Time `json:"started_at"`
	DataSent  int64     `json:"data_sent"`
	DataRecv  int64     `json:"data_recv"`
}

// ConnectionManager handles VPN connections
type ConnectionManager struct {
	status             ConnectionStatus
	currentServer      *core.Server
	connectionInfo     *ConnectionInfo
	mutex              sync.RWMutex
	statsManager       *stats.StatsManager
	notificationManager *notifications.NotificationManager
	logger             *logging.Logger
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		status: Disconnected,
	}
}

// SetStatsManager sets the stats manager
func (cm *ConnectionManager) SetStatsManager(statsManager *stats.StatsManager) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.statsManager = statsManager
}

// SetNotificationManager sets the notification manager
func (cm *ConnectionManager) SetNotificationManager(notificationManager *notifications.NotificationManager) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.notificationManager = notificationManager
}

// SetLogger sets the logger
func (cm *ConnectionManager) SetLogger(logger *logging.Logger) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.logger = logger
}

// GetStatus returns the current connection status
func (cm *ConnectionManager) GetStatus() ConnectionStatus {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.status
}

// GetCurrentServer returns the currently connected server
func (cm *ConnectionManager) GetCurrentServer() *core.Server {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.currentServer
}

// GetConnectionInfo returns information about the current connection
func (cm *ConnectionManager) GetConnectionInfo() *ConnectionInfo {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	if cm.connectionInfo == nil {
		return nil
	}
	
	// Return a copy to prevent external modification
	info := *cm.connectionInfo
	return &info
}

// Connect establishes a VPN connection to the specified server
func (cm *ConnectionManager) Connect(server *core.Server) error {
	cm.mutex.Lock()
	
	// Check if already connected or connecting
	if cm.status == Connected || cm.status == Connecting {
		cm.mutex.Unlock()
		return fmt.Errorf("already connected or connecting to a server")
	}
	
	// Set status to connecting
	cm.status = Connecting
	cm.currentServer = server
	
	// Initialize connection info
	cm.connectionInfo = &ConnectionInfo{
		ID:        server.ID,
		Name:      server.Name,
		StartedAt: time.Now(),
		DataSent:  0,
		DataRecv:  0,
	}
	
	cm.mutex.Unlock()
	
	// Log connection attempt
	if cm.logger != nil {
		cm.logger.Info("Attempting to connect to server: %s (%s:%d)", server.Name, server.Host, server.Port)
	}
	
	// Send connecting notification
	if cm.notificationManager != nil {
		cm.notificationManager.AddNotification("Connecting", fmt.Sprintf("Connecting to %s", server.Name), notifications.Info)
	}
	
	var err error
	
	// Connect based on protocol
	switch server.Protocol {
	case "wireguard":
		err = cm.connectWireGuard(server)
	case "vmess":
		err = cm.connectVMess(server)
	case "shadowsocks":
		err = cm.connectShadowsocks(server)
	case "trojan":
		err = cm.connectTrojan(server)
	default:
		err = fmt.Errorf("unsupported protocol: %s", server.Protocol)
	}
	
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if err != nil {
		// Set status to error
		cm.status = Error
		cm.currentServer = nil
		cm.connectionInfo = nil
		
		// Log error
		if cm.logger != nil {
			cm.logger.Error("Failed to connect to server %s: %v", server.Name, err)
		}
		
		// Send error notification
		if cm.notificationManager != nil {
			cm.notificationManager.AddNotification("Connection Failed", fmt.Sprintf("Failed to connect to %s: %v", server.Name, err), notifications.Error)
		}
		
		return err
	}
	
	// Set status to connected
	cm.status = Connected
	
	// Log successful connection
	if cm.logger != nil {
		cm.logger.Info("Successfully connected to server: %s", server.Name)
	}
	
	// Send success notification
	if cm.notificationManager != nil {
		cm.notificationManager.AddNotification("Connected", fmt.Sprintf("Successfully connected to %s", server.Name), notifications.Success)
	}
	
	// Start stats collection if stats manager is available
	if cm.statsManager != nil {
		cm.statsManager.StartConnection(cm.connectionInfo.ID, cm.connectionInfo.Name)
	}
	
	return nil
}

// Disconnect terminates the current VPN connection
func (cm *ConnectionManager) Disconnect() error {
	cm.mutex.Lock()
	
	// Check if connected or connecting
	if cm.status != Connected && cm.status != Connecting && cm.status != Error {
		cm.mutex.Unlock()
		return fmt.Errorf("not connected to any server")
	}
	
	// Save server info for notification
	serverName := cm.currentServer.Name
	
	// Set status to disconnecting
	cm.status = Disconnecting
	
	cm.mutex.Unlock()
	
	// Log disconnection attempt
	if cm.logger != nil {
		cm.logger.Info("Attempting to disconnect from server: %s", serverName)
	}
	
	// Send disconnecting notification
	if cm.notificationManager != nil {
		cm.notificationManager.AddNotification("Disconnecting", fmt.Sprintf("Disconnecting from %s", serverName), notifications.Info)
	}
	
	var err error
	
	// Disconnect based on current protocol
	if cm.currentServer != nil {
		switch cm.currentServer.Protocol {
		case "wireguard":
			err = cm.disconnectWireGuard()
		case "vmess":
			err = cm.disconnectVMess()
		case "shadowsocks":
			err = cm.disconnectShadowsocks()
		case "trojan":
			err = cm.disconnectTrojan()
		default:
			err = fmt.Errorf("unsupported protocol: %s", cm.currentServer.Protocol)
		}
	}
	
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	// Stop stats collection if stats manager is available
	if cm.statsManager != nil {
		cm.statsManager.EndConnection()
	}
	
	if err != nil {
		// Set status to error
		cm.status = Error
		
		// Log error
		if cm.logger != nil {
			cm.logger.Error("Failed to disconnect from server %s: %v", serverName, err)
		}
		
		// Send error notification
		if cm.notificationManager != nil {
			cm.notificationManager.AddNotification("Disconnection Failed", fmt.Sprintf("Failed to disconnect from %s: %v", serverName, err), notifications.Error)
		}
		
		return err
	}
	
	// Set status to disconnected
	cm.status = Disconnected
	cm.currentServer = nil
	cm.connectionInfo = nil
	
	// Log successful disconnection
	if cm.logger != nil {
		cm.logger.Info("Successfully disconnected from server: %s", serverName)
	}
	
	// Send success notification
	if cm.notificationManager != nil {
		cm.notificationManager.AddNotification("Disconnected", fmt.Sprintf("Successfully disconnected from %s", serverName), notifications.Success)
	}
	
	return nil
}

// UpdateStats updates the connection statistics
func (cm *ConnectionManager) UpdateStats(sent, recv int64) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if cm.connectionInfo != nil {
		cm.connectionInfo.DataSent = sent
		cm.connectionInfo.DataRecv = recv
		
		// Update stats manager if available
		if cm.statsManager != nil {
			cm.statsManager.UpdateConnection(sent, recv)
		}
	}
}

// connectWireGuard establishes a WireGuard connection
func (cm *ConnectionManager) connectWireGuard(server *core.Server) error {
	// This is a placeholder implementation
	// In a real implementation, this would use the wireguard-go library or wgctrl
	return nil
}

// disconnectWireGuard terminates a WireGuard connection
func (cm *ConnectionManager) disconnectWireGuard() error {
	// This is a placeholder implementation
	// In a real implementation, this would use the wireguard-go library or wgctrl
	return nil
}

// connectVMess establishes a VMess connection
func (cm *ConnectionManager) connectVMess(server *core.Server) error {
	// This is a placeholder implementation
	// In a real implementation, this would use the v2ray-core library
	return nil
}

// disconnectVMess terminates a VMess connection
func (cm *ConnectionManager) disconnectVMess() error {
	// This is a placeholder implementation
	// In a real implementation, this would use the v2ray-core library
	return nil
}

// connectShadowsocks establishes a Shadowsocks connection
func (cm *ConnectionManager) connectShadowsocks(server *core.Server) error {
	// This is a placeholder implementation
	// In a real implementation, this would use the shadowsocks-go library
	return nil
}

// disconnectShadowsocks terminates a Shadowsocks connection
func (cm *ConnectionManager) disconnectShadowsocks() error {
	// This is a placeholder implementation
	// In a real implementation, this would use the shadowsocks-go library
	return nil
}

// connectTrojan establishes a Trojan connection
func (cm *ConnectionManager) connectTrojan(server *core.Server) error {
	// This is a placeholder implementation
	// In a real implementation, this would use the trojan-go library
	return nil
}

// disconnectTrojan terminates a Trojan connection
func (cm *ConnectionManager) disconnectTrojan() error {
	// This is a placeholder implementation
	// In a real implementation, this would use the trojan-go library
	return nil
}