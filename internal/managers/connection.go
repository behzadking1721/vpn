package managers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vpnclient/internal/logging"
	"vpnclient/internal/notifications"
	"vpnclient/internal/stats"
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
	status             ConnectionStatus
	currentServer      *core.Server
	startTime          time.Time
	dataSent           int64
	dataReceived       int64
	mutex              sync.RWMutex
	ctx                context.Context
	cancel             context.CancelFunc
	notificationManager *notifications.NotificationManager
	logger             *logging.Logger
	statsManager       *stats.StatsManager
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	// Create a context for managing connection lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	
	return &ConnectionManager{
		status:  Disconnected,
		ctx:     ctx,
		cancel:  cancel,
		notificationManager: nil,
		logger:  nil,
		statsManager: nil,
	}
}

// SetNotificationManager sets the notification manager
func (cm *ConnectionManager) SetNotificationManager(nm *notifications.NotificationManager) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.notificationManager = nm
}

// SetLogger sets the logger
func (cm *ConnectionManager) SetLogger(logger *logging.Logger) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.logger = logger
}

// SetStatsManager sets the stats manager
func (cm *ConnectionManager) SetStatsManager(statsManager *stats.StatsManager) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.statsManager = statsManager
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

// GetCurrentServer returns the current server
func (cm *ConnectionManager) GetCurrentServer() *core.Server {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.currentServer
}

// GetConnectionInfo returns connection information
func (cm *ConnectionManager) GetConnectionInfo() *core.ConnectionInfo {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	if cm.status != Connected {
		return nil
	}
	
	return &core.ConnectionInfo{
		ID:        cm.currentServer.ID,
		StartedAt: cm.startTime,
		DataSent:  cm.dataSent,
		DataRecv:  cm.dataReceived,
	}
}

// Connect connects to a VPN server
func (cm *ConnectionManager) Connect(server *core.Server) error {
	cm.mutex.Lock()
	
	// Check if we're already connected or connecting
	if cm.status == Connecting || cm.status == Connected {
		cm.mutex.Unlock()
		return fmt.Errorf("already connected or connecting to a server")
	}
	
	// Update status
	cm.status = Connecting
	cm.currentServer = server
	
	// Log connection attempt
	if cm.logger != nil {
		cm.logger.LogConnectionEvent("Connecting", server.ID, map[string]interface{}{
			"server_name": server.Name,
			"server_host": server.Host,
			"server_port": server.Port,
			"protocol":    server.Protocol,
		})
	}
	
	cm.mutex.Unlock()
	
	// Send notification
	if cm.notificationManager != nil {
		cm.notificationManager.AddNotification(
			"Connecting",
			fmt.Sprintf("Connecting to %s", server.Name),
			notifications.Info,
		)
	}
	
	// Perform the connection based on protocol
	var err error
	switch server.Protocol {
	case "vmess":
		err = cm.connectVMess(server)
	case "shadowsocks":
		err = cm.connectShadowsocks(server)
	case "trojan":
		err = cm.connectTrojan(server)
	case "wireguard":
		// Check which WireGuard implementation to use
		if isWireGuardAvailable() {
			err = cm.connectWireGuard(server)
		} else {
			err = fmt.Errorf("wireguard not available in this build")
		}
	default:
		err = fmt.Errorf("unsupported protocol: %s", server.Protocol)
	}
	
	cm.mutex.Lock()
	if err != nil {
		// Connection failed, revert status
		cm.status = Error
		cm.currentServer = nil
		
		// Log connection failure
		if cm.logger != nil {
			cm.logger.LogConnectionEvent("Connection Failed", server.ID, map[string]interface{}{
				"error": err.Error(),
			})
		}
		
		// Send error notification
		if cm.notificationManager != nil {
			cm.notificationManager.AddNotification(
				"Connection Failed",
				fmt.Sprintf("Failed to connect to %s: %v", server.Name, err),
				notifications.Error,
			)
		}
	} else {
		// Connection successful
		cm.status = Connected
		cm.startTime = time.Now()
		cm.dataSent = 0
		cm.dataReceived = 0
		
		// Start tracking statistics
		if cm.statsManager != nil {
			cm.statsManager.StartConnection(server.ID, server.Name)
		}
		
		// Log successful connection
		if cm.logger != nil {
			cm.logger.LogConnectionEvent("Connected", server.ID, map[string]interface{}{
				"start_time": cm.startTime.Format(time.RFC3339),
			})
		}
		
		// Send success notification
		if cm.notificationManager != nil {
			cm.notificationManager.AddNotification(
				"Connected",
				fmt.Sprintf("Successfully connected to %s", server.Name),
				notifications.Success,
			)
		}
	}
	cm.mutex.Unlock()
	
	return err
}

// Disconnect disconnects from the current VPN server
func (cm *ConnectionManager) Disconnect() error {
	cm.mutex.Lock()
	
	// Check if we're connected or connecting
	if cm.status != Connected && cm.status != Connecting {
		cm.mutex.Unlock()
		return fmt.Errorf("not connected to any server")
	}
	
	// Update status
	cm.status = Disconnecting
	server := cm.currentServer
	
	// Log disconnection attempt
	if cm.logger != nil {
		cm.logger.LogConnectionEvent("Disconnecting", server.ID, nil)
	}
	
	cm.mutex.Unlock()
	
	// Send notification
	if cm.notificationManager != nil {
		cm.notificationManager.AddNotification(
			"Disconnecting",
			fmt.Sprintf("Disconnecting from %s", server.Name),
			notifications.Info,
		)
	}
	
	// Perform the disconnection based on protocol
	var err error
	switch server.Protocol {
	case "vmess":
		err = cm.disconnectVMess(server)
	case "shadowsocks":
		err = cm.disconnectShadowsocks(server)
	case "trojan":
		err = cm.disconnectTrojan(server)
	case "wireguard":
		if isWireGuardAvailable() {
			err = cm.disconnectWireGuard(server)
		} else {
			err = fmt.Errorf("wireguard not available in this build")
		}
	default:
		err = fmt.Errorf("unsupported protocol: %s", server.Protocol)
	}
	
	cm.mutex.Lock()
	if err != nil {
		// Disconnection failed, set error status
		cm.status = Error
		
		// Log disconnection failure
		if cm.logger != nil {
			cm.logger.LogConnectionEvent("Disconnection Failed", server.ID, map[string]interface{}{
				"error": err.Error(),
			})
		}
		
		// Send error notification
		if cm.notificationManager != nil {
			cm.notificationManager.AddNotification(
				"Disconnection Failed",
				fmt.Sprintf("Failed to disconnect from %s: %v", server.Name, err),
				notifications.Error,
			)
		}
	} else {
		// Disconnection successful
		cm.status = Disconnected
		cm.currentServer = nil
		cm.startTime = time.Time{}
		cm.dataSent = 0
		cm.dataReceived = 0
		
		// End statistics tracking
		if cm.statsManager != nil {
			cm.statsManager.EndConnection()
		}
		
		// Log successful disconnection
		if cm.logger != nil {
			cm.logger.LogConnectionEvent("Disconnected", server.ID, nil)
		}
		
		// Send success notification
		if cm.notificationManager != nil {
			cm.notificationManager.AddNotification(
				"Disconnected",
				fmt.Sprintf("Successfully disconnected from %s", server.Name),
				notifications.Success,
			)
		}
	}
	cm.mutex.Unlock()
	
	// Cancel any ongoing context
	if cm.cancel != nil {
		cm.cancel()
	}
	
	// Create a new context for future connections
	cm.ctx, cm.cancel = context.WithCancel(context.Background())
	
	return err
}

// UpdateStats updates the connection statistics
func (cm *ConnectionManager) UpdateStats(sent, received int64) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	cm.dataSent += sent
	cm.dataReceived += received
	
	// Update statistics manager
	if cm.statsManager != nil {
		cm.statsManager.UpdateConnection(sent, received)
	}
}

// connectVMess connects to a VMess server
func (cm *ConnectionManager) connectVMess(server *core.Server) error {
	// Implementation would go here
	// This is a placeholder for actual VMess connection logic
	
	// Log VMess connection attempt
	if cm.logger != nil {
		cm.logger.Debug("Attempting VMess connection to %s:%d", server.Host, server.Port)
	}
	
	return nil
}

// disconnectVMess disconnects from a VMess server
func (cm *ConnectionManager) disconnectVMess(server *core.Server) error {
	// Implementation would go here
	// This is a placeholder for actual VMess disconnection logic
	
	// Log VMess disconnection
	if cm.logger != nil {
		cm.logger.Debug("Disconnecting from VMess server %s:%d", server.Host, server.Port)
	}
	
	return nil
}

// connectShadowsocks connects to a Shadowsocks server
func (cm *ConnectionManager) connectShadowsocks(server *core.Server) error {
	// Implementation would go here
	// This is a placeholder for actual Shadowsocks connection logic
	
	// Log Shadowsocks connection attempt
	if cm.logger != nil {
		cm.logger.Debug("Attempting Shadowsocks connection to %s:%d", server.Host, server.Port)
	}
	
	return nil
}

// disconnectShadowsocks disconnects from a Shadowsocks server
func (cm *ConnectionManager) disconnectShadowsocks(server *core.Server) error {
	// Implementation would go here
	// This is a placeholder for actual Shadowsocks disconnection logic
	
	// Log Shadowsocks disconnection
	if cm.logger != nil {
		cm.logger.Debug("Disconnecting from Shadowsocks server %s:%d", server.Host, server.Port)
	}
	
	return nil
}

// connectTrojan connects to a Trojan server
func (cm *ConnectionManager) connectTrojan(server *core.Server) error {
	// Implementation would go here
	// This is a placeholder for actual Trojan connection logic
	
	// Log Trojan connection attempt
	if cm.logger != nil {
		cm.logger.Debug("Attempting Trojan connection to %s:%d", server.Host, server.Port)
	}
	
	return nil
}

// disconnectTrojan disconnects from a Trojan server
func (cm *ConnectionManager) disconnectTrojan(server *core.Server) error {
	// Implementation would go here
	// This is a placeholder for actual Trojan disconnection logic
	
	// Log Trojan disconnection
	if cm.logger != nil {
		cm.logger.Debug("Disconnecting from Trojan server %s:%d", server.Host, server.Port)
	}
	
	return nil
}

// isWireGuardAvailable checks if WireGuard is available in this build
func isWireGuardAvailable() bool {
	// This will be overridden by build tags
	return false
}