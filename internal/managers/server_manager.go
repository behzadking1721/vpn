package managers

import (
	"fmt"
	"strings"
	"time"

	"vpnclient/internal/database"
	"vpnclient/src/core"
)

// ServerManager manages VPN servers
type ServerManager struct {
	store *database.ServerStoreWrapper
	subStore *database.SubscriptionStoreWrapper
}

// NewServerManager creates a new server manager
func NewServerManager(store database.Store) *ServerManager {
	return &ServerManager{
		store: database.NewServerStore(store),
		subStore: database.NewSubscriptionStore(store),
	}
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

	return sm.store.AddServer(server)
}

// GetServer retrieves a server by ID
func (sm *ServerManager) GetServer(id string) (*core.Server, error) {
	if id == "" {
		return nil, fmt.Errorf("server ID is required")
	}
	return sm.store.GetServer(id)
}

// GetAllServers retrieves all servers
func (sm *ServerManager) GetAllServers() ([]*core.Server, error) {
	return sm.store.GetAllServers()
}

// GetEnabledServers retrieves only enabled servers
func (sm *ServerManager) GetEnabledServers() ([]*core.Server, error) {
	return sm.store.GetEnabledServers()
}

// UpdateServer updates an existing server
func (sm *ServerManager) UpdateServer(server *core.Server) error {
	if err := sm.validateServer(server); err != nil {
		return err
	}
	return sm.store.UpdateServer(server)
}

// DeleteServer deletes a server by ID
func (sm *ServerManager) DeleteServer(id string) error {
	if id == "" {
		return fmt.Errorf("server ID is required")
	}
	return sm.store.DeleteServer(id)
}

// EnableServer enables a server
func (sm *ServerManager) EnableServer(id string) error {
	server, err := sm.store.GetServer(id)
	if err != nil {
		return err
	}
	server.Enabled = true
	return sm.store.UpdateServer(server)
}

// DisableServer disables a server
func (sm *ServerManager) DisableServer(id string) error {
	server, err := sm.store.GetServer(id)
	if err != nil {
		return err
	}
	server.Enabled = false
	return sm.store.UpdateServer(server)
}

// UpdatePing updates the ping value for a server
func (sm *ServerManager) UpdatePing(id string, ping int) error {
	server, err := sm.store.GetServer(id)
	if err != nil {
		return err
	}
	
	if ping < 0 {
		return fmt.Errorf("ping cannot be negative")
	}
	
	server.Ping = ping
	return sm.store.UpdateServer(server)
}

// GetFastestServer returns the enabled server with the lowest ping
func (sm *ServerManager) GetFastestServer() (*core.Server, error) {
	servers, err := sm.GetEnabledServers()
	if err != nil {
		return nil, err
	}
	
	if len(servers) == 0 {
		return nil, fmt.Errorf("no enabled servers available")
	}
	
	// Find server with lowest ping
	fastest := servers[0]
	for _, server := range servers[1:] {
		// Consider servers with ping=0 as having high latency
		if server.Ping > 0 && (fastest.Ping == 0 || server.Ping < fastest.Ping) {
			fastest = server
		}
	}
	
	return fastest, nil
}

// AddSubscription adds a new subscription
func (sm *ServerManager) AddSubscription(sub *core.Subscription) error {
	// Validate subscription
	if sub.URL == "" {
		return fmt.Errorf("subscription URL is required")
	}
	
	if sub.Name == "" {
		sub.Name = "Unnamed Subscription"
	}
	
	// Generate ID if not provided
	if sub.ID == "" {
		sub.ID = sm.generateSubscriptionID(sub)
	}
	
	// Set timestamps if not set
	now := time.Now()
	if sub.CreatedAt.IsZero() {
		sub.CreatedAt = now
	}
	sub.UpdatedAt = now
	
	return sm.subStore.AddSubscription(sub)
}

// GetSubscription retrieves a subscription by ID
func (sm *ServerManager) GetSubscription(id string) (*core.Subscription, error) {
	if id == "" {
		return nil, fmt.Errorf("subscription ID is required")
	}
	return sm.subStore.GetSubscription(id)
}

// GetAllSubscriptions retrieves all subscriptions
func (sm *ServerManager) GetAllSubscriptions() ([]*core.Subscription, error) {
	return sm.subStore.GetAllSubscriptions()
}

// UpdateSubscription updates an existing subscription
func (sm *ServerManager) UpdateSubscription(sub *core.Subscription) error {
	if sub.ID == "" {
		return fmt.Errorf("subscription ID is required")
	}
	
	sub.UpdatedAt = time.Now()
	return sm.subStore.UpdateSubscription(sub)
}

// DeleteSubscription deletes a subscription by ID
func (sm *ServerManager) DeleteSubscription(id string) error {
	if id == "" {
		return fmt.Errorf("subscription ID is required")
	}
	return sm.subStore.DeleteSubscription(id)
}

// UpdateSubscriptionServers updates servers from a subscription
func (sm *ServerManager) UpdateSubscriptionServers(id string) error {
	sub, err := sm.GetSubscription(id)
	if err != nil {
		return err
	}
	
	// TODO: Parse subscription URL and import servers
	// This would involve:
	// 1. Fetching the subscription content from sub.URL
	// 2. Parsing the format (Clash, Surfboard, etc.)
	// 3. Converting to our server format
	// 4. Adding/updating servers in the database
	
	// For now, just update the last update time
	sub.LastUpdate = time.Now()
	sub.UpdatedAt = time.Now()
	
	return sm.UpdateSubscription(sub)
}

// validateServer validates a server's fields
func (sm *ServerManager) validateServer(server *core.Server) error {
	if server.Name == "" {
		return fmt.Errorf("server name is required")
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
	
	return nil
}

// generateServerID generates a unique ID for a server based on its properties
func (sm *ServerManager) generateServerID(server *core.Server) string {
	// Simple ID generation based on server properties
	// In a real implementation, you might use UUID or a more sophisticated approach
	return fmt.Sprintf("%s-%s-%d", server.Host, server.Protocol, server.Port)
}

// generateSubscriptionID generates a unique ID for a subscription
func (sm *ServerManager) generateSubscriptionID(sub *core.Subscription) string {
	// Simple ID generation based on URL
	// In a real implementation, you might use UUID
	parts := strings.Split(sub.URL, "//")
	if len(parts) > 1 {
		return strings.ReplaceAll(parts[1], "/", "-")
	}
	return "subscription-" + time.Now().Format("20060102150405")
}