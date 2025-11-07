package managers

import (
	"fmt"
	"sync"
	"time"

	"vpnclient/internal/database"
	"vpnclient/internal/logging"
	"vpnclient/internal/notifications"
	"vpnclient/src/core"
)

// ServerManager handles server management operations
type ServerManager struct {
	store               database.Store
	mutex               sync.RWMutex
	notificationManager *notifications.NotificationManager
	logger              *logging.Logger
	serverCache         map[string]*core.Server
	cacheMutex          sync.RWMutex
	// optional subscription manager (can be set by caller)
	subscriptionManager *SubscriptionManager
}

// NewServerManager creates a new server manager
func NewServerManager(store database.Store) *ServerManager {
	return &ServerManager{
		store:       store,
		serverCache: make(map[string]*core.Server),
	}
}

// SetSubscriptionManager attaches a SubscriptionManager so ServerManager can delegate subscription ops
func (sm *ServerManager) SetSubscriptionManager(sub *SubscriptionManager) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.subscriptionManager = sub
}

// AddSubscription stores a subscription record. If a SubscriptionManager is present, delegate to it.
func (sm *ServerManager) AddSubscription(sub *core.Subscription) error {
	if sm.subscriptionManager != nil {
		// delegate by URL if provided
		if sub != nil && sub.URL != "" {
			_, err := sm.subscriptionManager.AddSubscription(sub.URL)
			return err
		}
	}
	return sm.store.AddSubscription(sub)
}

// GetAllSubscriptions returns subscriptions from the underlying store
func (sm *ServerManager) GetAllSubscriptions() ([]*core.Subscription, error) {
	// try subscription manager first
	if sm.subscriptionManager != nil {
		return sm.subscriptionManager.GetAllSubscriptions()
	}
	items, err := sm.store.GetAllSubscriptions()
	if err != nil {
		return nil, err
	}
	out := make([]*core.Subscription, 0, len(items))
	for _, it := range items {
		if s, ok := it.(*core.Subscription); ok {
			out = append(out, s)
		}
	}
	return out, nil
}

// GetSubscription fetches a subscription by id
func (sm *ServerManager) GetSubscription(id string) (*core.Subscription, error) {
	if sm.subscriptionManager != nil {
		return sm.subscriptionManager.GetSubscription(id)
	}
	it, err := sm.store.GetSubscription(id)
	if err != nil {
		return nil, err
	}
	if s, ok := it.(*core.Subscription); ok {
		return s, nil
	}
	return nil, fmt.Errorf("invalid subscription type in store")
}

// UpdateSubscription updates a subscription record
func (sm *ServerManager) UpdateSubscription(sub *core.Subscription) error {
	if sm.subscriptionManager != nil {
		return sm.subscriptionManager.UpdateSubscription(sub)
	}
	return sm.store.UpdateSubscription(sub)
}

// DeleteSubscription deletes a subscription
func (sm *ServerManager) DeleteSubscription(id string) error {
	if sm.subscriptionManager != nil {
		return sm.subscriptionManager.DeleteSubscription(id)
	}
	return sm.store.DeleteSubscription(id)
}

// UpdateSubscriptionServers delegates to subscription manager when available
func (sm *ServerManager) UpdateSubscriptionServers(id string) error {
	if sm.subscriptionManager != nil {
		return sm.subscriptionManager.UpdateSubscriptionServers(id)
	}
	return fmt.Errorf("subscription manager not configured")
}

// TestServerPing performs a simple simulated ping for a single server and returns ping in ms
func (sm *ServerManager) TestServerPing(id string) (int, error) {
	s, err := sm.GetServer(id)
	if err != nil {
		return 0, err
	}
	// simulate ping
	ping := 50 + int(time.Now().UnixNano()%100)
	s.Ping = ping
	if err := sm.UpdateServer(s); err != nil {
		return 0, err
	}
	return ping, nil
}

// SetNotificationManager sets the notification manager
func (sm *ServerManager) SetNotificationManager(notificationManager *notifications.NotificationManager) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.notificationManager = notificationManager
}

// SetLogger sets the logger
func (sm *ServerManager) SetLogger(logger *logging.Logger) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.logger = logger
}

// AddServer adds a new server to the manager
func (sm *ServerManager) AddServer(server *core.Server) error {
	// Validate server
	if err := sm.validateServer(server); err != nil {
		return err
	}

	// Add to store
	if err := sm.store.AddServer(server); err != nil {
		return fmt.Errorf("failed to add server to store: %w", err)
	}

	// Add to cache
	sm.cacheMutex.Lock()
	sm.serverCache[server.ID] = server
	sm.cacheMutex.Unlock()

	// Send notification
	if sm.notificationManager != nil {
		sm.notificationManager.AddNotification("Server Added", fmt.Sprintf("New server %s added successfully", server.Name), notifications.Success)
	}

	// Log
	if sm.logger != nil {
		sm.logger.Info("Server added: %s (%s)", server.Name, server.ID)
	}

	return nil
}

// GetServer retrieves a server by ID
func (sm *ServerManager) GetServer(id string) (*core.Server, error) {
	// Check cache first
	sm.cacheMutex.RLock()
	if server, exists := sm.serverCache[id]; exists {
		sm.cacheMutex.RUnlock()
		return server, nil
	}
	sm.cacheMutex.RUnlock()

	// Get from store
	item, err := sm.store.GetServer(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get server from store: %w", err)
	}

	server, ok := item.(*core.Server)
	if !ok {
		return nil, fmt.Errorf("invalid server data in store")
	}

	// Add to cache
	sm.cacheMutex.Lock()
	sm.serverCache[id] = server
	sm.cacheMutex.Unlock()

	return server, nil
}

// UpdateServer updates an existing server
func (sm *ServerManager) UpdateServer(server *core.Server) error {
	// Validate server
	if err := sm.validateServer(server); err != nil {
		return err
	}

	// Update in store
	if err := sm.store.UpdateServer(server); err != nil {
		return fmt.Errorf("failed to update server in store: %w", err)
	}

	// Update in cache
	sm.cacheMutex.Lock()
	sm.serverCache[server.ID] = server
	sm.cacheMutex.Unlock()

	// Send notification
	if sm.notificationManager != nil {
		sm.notificationManager.AddNotification("Server Updated", fmt.Sprintf("Server %s updated successfully", server.Name), notifications.Success)
	}

	// Log
	if sm.logger != nil {
		sm.logger.Info("Server updated: %s (%s)", server.Name, server.ID)
	}

	return nil
}

// DeleteServer removes a server by ID
func (sm *ServerManager) DeleteServer(id string) error {
	// Delete from store
	if err := sm.store.DeleteServer(id); err != nil {
		return fmt.Errorf("failed to delete server from store: %w", err)
	}

	// Remove from cache
	sm.cacheMutex.Lock()
	delete(sm.serverCache, id)
	sm.cacheMutex.Unlock()

	// Send notification
	if sm.notificationManager != nil {
		sm.notificationManager.AddNotification("Server Deleted", "Server deleted successfully", notifications.Success)
	}

	// Log
	if sm.logger != nil {
		sm.logger.Info("Server deleted: %s", id)
	}

	return nil
}

// GetAllServers retrieves all servers
func (sm *ServerManager) GetAllServers() ([]*core.Server, error) {
	// Get all items from store
	items, err := sm.store.GetAllServers()
	if err != nil {
		return nil, fmt.Errorf("failed to get servers from store: %w", err)
	}

	// Convert items to servers
	servers := make([]*core.Server, len(items))
	for i, item := range items {
		server, ok := item.(*core.Server)
		if !ok {
			return nil, fmt.Errorf("invalid server data in store")
		}
		servers[i] = server
	}

	return servers, nil
}

// GetEnabledServers retrieves all enabled servers
func (sm *ServerManager) GetEnabledServers() ([]*core.Server, error) {
	allServers, err := sm.GetAllServers()
	if err != nil {
		return nil, err
	}

	enabledServers := make([]*core.Server, 0)
	for _, server := range allServers {
		if server.Enabled {
			enabledServers = append(enabledServers, server)
		}
	}

	return enabledServers, nil
}

// EnableServer enables a server
func (sm *ServerManager) EnableServer(id string) error {
	server, err := sm.GetServer(id)
	if err != nil {
		return err
	}

	server.Enabled = true
	return sm.UpdateServer(server)
}

// DisableServer disables a server
func (sm *ServerManager) DisableServer(id string) error {
	server, err := sm.GetServer(id)
	if err != nil {
		return err
	}

	server.Enabled = false
	return sm.UpdateServer(server)
}

// UpdatePing updates the ping value for a server
func (sm *ServerManager) UpdatePing(id string, ping int) error {
	server, err := sm.GetServer(id)
	if err != nil {
		return err
	}

	server.Ping = ping
	return sm.UpdateServer(server)
}

// GetBestServer returns the enabled server with the lowest ping
func (sm *ServerManager) GetBestServer() (*core.Server, error) {
	servers, err := sm.GetEnabledServers()
	if err != nil {
		return nil, err
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no enabled servers available")
	}

	bestServer := servers[0]
	for _, server := range servers[1:] {
		if server.Ping < bestServer.Ping {
			bestServer = server
		}
	}

	return bestServer, nil
}

// TestAllServersPing tests the ping for all servers
func (sm *ServerManager) TestAllServersPing() error {
	servers, err := sm.GetAllServers()
	if err != nil {
		return err
	}

	for _, server := range servers {
		// In a real implementation, we would actually ping the server
		// For now, we'll just simulate a ping test
		ping := 50 + int(time.Now().UnixNano()%100) // Simulate ping between 50-150ms
		server.Ping = ping

		if err := sm.UpdateServer(server); err != nil {
			// Log error but continue with other servers
			if sm.logger != nil {
				sm.logger.Error("Failed to update ping for server %s: %v", server.Name, err)
			}
		}
	}

	return nil
}

// validateServer validates a server's fields
func (sm *ServerManager) validateServer(server *core.Server) error {
	if server.Host == "" {
		return fmt.Errorf("server host is required")
	}

	if server.Port <= 0 || server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	if server.Protocol == "" {
		return fmt.Errorf("server protocol is required")
	}

	validProtocols := map[string]bool{
		"wireguard":   true,
		"vmess":       true,
		"shadowsocks": true,
		"trojan":      true,
	}

	if !validProtocols[server.Protocol] {
		return fmt.Errorf("invalid protocol: %s", server.Protocol)
	}

	return nil
}
