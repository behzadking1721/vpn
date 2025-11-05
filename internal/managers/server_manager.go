package managers

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"vpnclient/internal/database"
	"vpnclient/internal/logging"
	"vpnclient/internal/notifications"
	"vpnclient/src/core"
)

// ServerManager manages VPN servers
type ServerManager struct {
	store              *database.ServerStoreWrapper
	subStore           *database.SubscriptionStoreWrapper
	cache              map[string]*core.Server
	cacheMux           sync.RWMutex
	cacheTTL           time.Duration
	notificationManager *notifications.NotificationManager
	logger             *logging.Logger
}

// NewServerManager creates a new server manager
func NewServerManager(store database.Store) *ServerManager {
	return &ServerManager{
		store:    database.NewServerStore(store),
		subStore: database.NewSubscriptionStore(store),
		cache:    make(map[string]*core.Server),
		cacheTTL: 5 * time.Minute, // Cache TTL of 5 minutes
	}
}

// SetNotificationManager sets the notification manager
func (sm *ServerManager) SetNotificationManager(nm *notifications.NotificationManager) {
	sm.notificationManager = nm
}

// SetLogger sets the logger
func (sm *ServerManager) SetLogger(logger *logging.Logger) {
	sm.logger = logger
}

// AddServer adds a new server
func (sm *ServerManager) AddServer(server *core.Server) error {
	// Validate server
	if err := sm.validateServer(server); err != nil {
		return err
	}

	// Generate ID if not provided
	if server.ID == "" {
		server.ID = sm.generateServerID(server)
	}

	// Set defaults
	if server.Config == nil {
		server.Config = make(map[string]interface{})
	}
	if !server.Enabled {
		server.Enabled = true // Default to enabled
	}

	// Add to store
	if err := sm.store.AddServer(server); err != nil {
		return err
	}

	// Update cache
	sm.cacheMux.Lock()
	sm.cache[server.ID] = server
	sm.cacheMux.Unlock()
	
	// Log server addition
	if sm.logger != nil {
		sm.logger.Info("Server added: %s (%s:%d)", server.Name, server.Host, server.Port)
	}
	
	// Send notification
	if sm.notificationManager != nil {
		sm.notificationManager.AddNotification(
			"Server Added",
			fmt.Sprintf("New server %s added successfully", server.Name),
			notifications.Success,
		)
	}

	return nil
}

// GetServer retrieves a server by ID
func (sm *ServerManager) GetServer(id string) (*core.Server, error) {
	if id == "" {
		return nil, fmt.Errorf("server ID is required")
	}

	// Check cache first
	sm.cacheMux.RLock()
	if server, exists := sm.cache[id]; exists {
		sm.cacheMux.RUnlock()
		return server, nil
	}
	sm.cacheMux.RUnlock()

	// Get from store
	server, err := sm.store.GetServer(id)
	if err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to get server %s: %v", id, err)
		}
		return nil, err
	}

	// Update cache
	sm.cacheMux.Lock()
	sm.cache[id] = server
	sm.cacheMux.Unlock()

	return server, nil
}

// GetAllServers retrieves all servers
func (sm *ServerManager) GetAllServers() ([]*core.Server, error) {
	servers, err := sm.store.GetAllServers()
	if err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to get all servers: %v", err)
		}
		return nil, err
	}
	
	// Log successful retrieval
	if sm.logger != nil {
		sm.logger.Debug("Retrieved %d servers", len(servers))
	}
	
	return servers, nil
}

// GetEnabledServers retrieves only enabled servers
func (sm *ServerManager) GetEnabledServers() ([]*core.Server, error) {
	servers, err := sm.store.GetAllServers()
	if err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to get enabled servers: %v", err)
		}
		return nil, err
	}

	// Filter only enabled servers
	enabledServers := make([]*core.Server, 0)
	for _, server := range servers {
		if server.Enabled {
			enabledServers = append(enabledServers, server)
		}
	}
	
	// Log successful retrieval
	if sm.logger != nil {
		sm.logger.Debug("Retrieved %d enabled servers out of %d total", len(enabledServers), len(servers))
	}

	return enabledServers, nil
}

// UpdateServer updates an existing server
func (sm *ServerManager) UpdateServer(server *core.Server) error {
	if err := sm.validateServer(server); err != nil {
		return err
	}
	
	// Update in store
	if err := sm.store.UpdateServer(server); err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to update server %s: %v", server.ID, err)
		}
		return err
	}

	// Update cache
	sm.cacheMux.Lock()
	sm.cache[server.ID] = server
	sm.cacheMux.Unlock()
	
	// Log server update
	if sm.logger != nil {
		sm.logger.Info("Server updated: %s (%s:%d)", server.Name, server.Host, server.Port)
	}
	
	// Send notification
	if sm.notificationManager != nil {
		sm.notificationManager.AddNotification(
			"Server Updated",
			fmt.Sprintf("Server %s updated successfully", server.Name),
			notifications.Success,
		)
	}

	return nil
}

// DeleteServer deletes a server by ID
func (sm *ServerManager) DeleteServer(id string) error {
	if id == "" {
		return fmt.Errorf("server ID is required")
	}
	
	// Get server name for notification
	server, err := sm.GetServer(id)
	if err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to get server %s for deletion: %v", id, err)
		}
		return err
	}
	
	// Delete from store
	if err := sm.store.DeleteServer(id); err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to delete server %s: %v", id, err)
		}
		return err
	}

	// Remove from cache
	sm.cacheMux.Lock()
	delete(sm.cache, id)
	sm.cacheMux.Unlock()
	
	// Log server deletion
	if sm.logger != nil {
		sm.logger.Info("Server deleted: %s (%s:%d)", server.Name, server.Host, server.Port)
	}
	
	// Send notification
	if sm.notificationManager != nil {
		sm.notificationManager.AddNotification(
			"Server Deleted",
			fmt.Sprintf("Server %s deleted successfully", server.Name),
			notifications.Success,
		)
	}

	return nil
}

// EnableServer enables a server
func (sm *ServerManager) EnableServer(id string) error {
	server, err := sm.GetServer(id)
	if err != nil {
		return err
	}
	server.Enabled = true
	
	// Log server enable
	if sm.logger != nil {
		sm.logger.Info("Server enabled: %s (%s:%d)", server.Name, server.Host, server.Port)
	}
	
	return sm.UpdateServer(server)
}

// DisableServer disables a server
func (sm *ServerManager) DisableServer(id string) error {
	server, err := sm.GetServer(id)
	if err != nil {
		return err
	}
	server.Enabled = false
	
	// Log server disable
	if sm.logger != nil {
		sm.logger.Info("Server disabled: %s (%s:%d)", server.Name, server.Host, server.Port)
	}
	
	return sm.UpdateServer(server)
}

// GetBestServer returns the best server based on ping and other factors
func (sm *ServerManager) GetBestServer() (*core.Server, error) {
	servers, err := sm.GetEnabledServers()
	if err != nil {
		return nil, err
	}

	if len(servers) == 0 {
		// Log warning
		if sm.logger != nil {
			sm.logger.Warning("No enabled servers available for best server selection")
		}
		return nil, fmt.Errorf("no enabled servers available")
	}

	// Find the server with the lowest ping
	var bestServer *core.Server
	lowestPing := -1

	for _, server := range servers {
		// If this is the first server or it has a lower ping
		if lowestPing == -1 || (server.Ping > 0 && server.Ping < lowestPing) {
			bestServer = server
			lowestPing = server.Ping
		}
	}

	// If no server has a valid ping, return the first one
	if bestServer == nil {
		bestServer = servers[0]
		
		// Log that we're selecting the first server due to lack of ping data
		if sm.logger != nil {
			sm.logger.Info("Selecting first server as no ping data available: %s", bestServer.Name)
		}
	} else {
		// Log best server selection
		if sm.logger != nil {
			sm.logger.Info("Selected best server: %s (ping: %d ms)", bestServer.Name, bestServer.Ping)
		}
	}

	return bestServer, nil
}

// UpdatePing updates the ping for a server
func (sm *ServerManager) UpdatePing(id string, ping int) error {
	server, err := sm.GetServer(id)
	if err != nil {
		return err
	}

	server.Ping = ping
	
	// Log ping update
	if sm.logger != nil {
		sm.logger.Debug("Updated ping for server %s: %d ms", server.Name, ping)
	}
	
	return sm.UpdateServer(server)
}

// TestServerPing tests the ping for a server
func (sm *ServerManager) TestServerPing(id string) (int, error) {
	server, err := sm.GetServer(id)
	if err != nil {
		return 0, err
	}

	ping, err := sm.testPing(server)
	if err != nil {
		// Log ping test failure
		if sm.logger != nil {
			sm.logger.Error("Ping test failed for server %s: %v", server.Name, err)
		}
		return 0, err
	}

	// Update the server's ping
	server.Ping = ping
	if err := sm.UpdateServer(server); err != nil {
		return 0, err
	}
	
	// Log successful ping test
	if sm.logger != nil {
		sm.logger.Info("Ping test completed for server %s: %d ms", server.Name, ping)
	}
	
	// Send notification for successful ping test
	if sm.notificationManager != nil && err == nil {
		sm.notificationManager.AddNotification(
			"Ping Test Complete",
			fmt.Sprintf("Server %s has ping %d ms", server.Name, ping),
			notifications.Info,
		)
	}

	return ping, nil
}

// TestAllServersPing tests the ping for all servers
func (sm *ServerManager) TestAllServersPing() error {
	servers, err := sm.GetAllServers()
	if err != nil {
		return err
	}

	// Log start of ping test
	if sm.logger != nil {
		sm.logger.Info("Starting ping test for %d servers", len(servers))
	}

	// Test ping for each server concurrently
	var wg sync.WaitGroup
	errChan := make(chan error, len(servers))
	successCount := 0
	failCount := 0
	
	var successMutex sync.Mutex

	for _, server := range servers {
		wg.Add(1)
		go func(s *core.Server) {
			defer wg.Done()
			ping, err := sm.testPing(s)
			if err != nil {
				// Don't return error, just log it
				errChan <- fmt.Errorf("failed to ping server %s: %v", s.ID, err)
				successMutex.Lock()
				failCount++
				successMutex.Unlock()
				
				// Log ping test failure
				if sm.logger != nil {
					sm.logger.Error("Ping test failed for server %s: %v", s.Name, err)
				}
				return
			}

			// Update the server's ping
			s.Ping = ping
			if err := sm.UpdateServer(s); err != nil {
				errChan <- fmt.Errorf("failed to update server %s: %v", s.ID, err)
				successMutex.Lock()
				failCount++
				successMutex.Unlock()
				
				// Log update failure
				if sm.logger != nil {
					sm.logger.Error("Failed to update ping for server %s: %v", s.Name, err)
				}
				return
			}
			
			successMutex.Lock()
			successCount++
			successMutex.Unlock()
			
			// Log successful ping test
			if sm.logger != nil {
				sm.logger.Debug("Ping test completed for server %s: %d ms", s.Name, ping)
			}
		}(server)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Log ping test completion
	if sm.logger != nil {
		sm.logger.Info("Completed ping test: %d successful, %d failed", successCount, failCount)
	}

	// Send notification with test results
	if sm.notificationManager != nil {
		message := fmt.Sprintf("Ping test completed. %d servers successful, %d failed", successCount, failCount)
		notifType := notifications.Success
		if failCount > 0 {
			notifType = notifications.Warning
		}
		
		sm.notificationManager.AddNotification(
			"All Servers Ping Test",
			message,
			notifType,
		)
	}

	// Return first error if any
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// testPing tests the ping for a server
func (sm *ServerManager) testPing(server *core.Server) (int, error) {
	if server.Host == "" {
		return 0, fmt.Errorf("server host is empty")
	}

	// Create a dialer with a timeout
	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}

	// Measure connection time
	start := time.Now()
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port))
	if err != nil {
		return 0, fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	ping := int(time.Since(start).Milliseconds())
	return ping, nil
}

// validateServer validates a server
func (sm *ServerManager) validateServer(server *core.Server) error {
	if server == nil {
		return fmt.Errorf("server is nil")
	}

	if server.Host == "" {
		return fmt.Errorf("server host is required")
	}

	if server.Port <= 0 || server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	if server.Protocol == "" {
		return fmt.Errorf("server protocol is required")
	}

	// Validate protocol
	validProtocols := map[string]bool{
		"vmess":      true,
		"shadowsocks": true,
		"trojan":     true,
		"wireguard":  true,
	}

	if !validProtocols[server.Protocol] {
		return fmt.Errorf("invalid protocol: %s", server.Protocol)
	}

	return nil
}

// generateServerID generates a unique ID for a server
func (sm *ServerManager) generateServerID(server *core.Server) string {
	// Simple ID generation based on host, port and protocol
	return fmt.Sprintf("%s_%d_%s", server.Host, server.Port, server.Protocol)
}

// ClearCache clears the server cache
func (sm *ServerManager) ClearCache() {
	sm.cacheMux.Lock()
	sm.cache = make(map[string]*core.Server)
	sm.cacheMux.Unlock()
	
	// Log cache clear
	if sm.logger != nil {
		sm.logger.Info("Server cache cleared")
	}
}