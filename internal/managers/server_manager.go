package managers

import (
	"fmt"
	"net"
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

// TestServerPing tests the ping for a specific server
func (sm *ServerManager) TestServerPing(id string) (int, error) {
	server, err := sm.GetServer(id)
	if err != nil {
		return 0, err
	}
	
	ping, err := sm.pingServer(server)
	if err != nil {
		return 0, err
	}
	
	// Update the server's ping value
	err = sm.UpdatePing(id, ping)
	if err != nil {
		return 0, fmt.Errorf("failed to update server ping: %v", err)
	}
	
	return ping, nil
}

// TestAllServersPing tests the ping for all enabled servers
func (sm *ServerManager) TestAllServersPing() (map[string]int, error) {
	servers, err := sm.GetEnabledServers()
	if err != nil {
		return nil, err
	}
	
	results := make(map[string]int)
	
	for _, server := range servers {
		ping, err := sm.pingServer(server)
		if err != nil {
			// If ping fails, set ping to a high value (e.g., 9999)
			ping = 9999
		}
		
		results[server.ID] = ping
		
		// Update the server's ping value
		err = sm.UpdatePing(server.ID, ping)
		if err != nil {
			// Log error but continue with other servers
			fmt.Printf("Failed to update ping for server %s: %v\n", server.ID, err)
		}
	}
	
	return results, nil
}

// pingServer performs a ping test to a server
func (sm *ServerManager) pingServer(server *core.Server) (int, error) {
	// Perform a TCP connection test to measure latency
	address := fmt.Sprintf("%s:%d", server.Host, server.Port)
	start := time.Now()
	
	// Try to establish a TCP connection with a timeout
	conn, err := net.DialTimeout("tcp", address, time.Second*5)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to server: %v", err)
	}
	
	// Close connection immediately
	conn.Close()
	
	// Calculate round-trip time
	ping := int(time.Since(start).Milliseconds())
	
	return ping, nil
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

// GetBestServer performs a comprehensive test to find the best server
func (sm *ServerManager) GetBestServer() (*core.Server, error) {
	// First, test ping for all servers
	_, err := sm.TestAllServersPing()
	if err != nil {
		return nil, fmt.Errorf("failed to test server pings: %v", err)
	}
	
	// Get the top 5 fastest servers based on ping
	topServers, err := sm.getTopServersByPing(5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top servers: %v", err)
	}
	
	if len(topServers) == 0 {
		return nil, fmt.Errorf("no servers available for deep testing")
	}
	
	// Perform deep speed tests on top servers
	bestServer := topServers[0]
	bestSpeed := int64(0)
	
	for _, server := range topServers {
		speed, err := sm.testServerSpeed(server)
		if err != nil {
			// If speed test fails, continue with other servers
			fmt.Printf("Speed test failed for server %s: %v\n", server.Name, err)
			continue
		}
		
		// Select server with highest speed
		if speed > bestSpeed {
			bestSpeed = speed
			bestServer = server
		}
	}
	
	return bestServer, nil
}

// getTopServersByPing returns the top N servers with the lowest ping
func (sm *ServerManager) getTopServersByPing(n int) ([]*core.Server, error) {
	servers, err := sm.GetEnabledServers()
	if err != nil {
		return nil, err
	}
	
	// Sort servers by ping (simple bubble sort for demonstration)
	for i := 0; i < len(servers)-1; i++ {
		for j := 0; j < len(servers)-i-1; j++ {
			// Consider servers with ping=0 as having high latency
			if (servers[j].Ping == 0 && servers[j+1].Ping > 0) ||
				(servers[j].Ping > 0 && servers[j+1].Ping > 0 && servers[j].Ping > servers[j+1].Ping) {
				servers[j], servers[j+1] = servers[j+1], servers[j]
			}
		}
	}
	
	// Return top N servers (or all servers if less than N)
	if len(servers) < n {
		n = len(servers)
	}
	
	return servers[:n], nil
}

// testServerSpeed performs a speed test on a server
func (sm *ServerManager) testServerSpeed(server *core.Server) (int64, error) {
	// In a real implementation, you would:
	// 1. Connect to the server
	// 2. Download a test file
	// 3. Measure the download speed
	// 4. Return the speed in bytes per second
	
	// For demonstration purposes, we'll simulate a speed test
	// In a real implementation, you would replace this with actual speed testing logic
	
	// Simulate a speed test by returning a random value based on ping
	// Lower ping generally indicates higher potential speed
	if server.Ping <= 0 {
		return 0, fmt.Errorf("invalid ping value")
	}
	
	// Simulate speed in Mbps (for demonstration)
	// In real implementation, this would be actual measured speed
	simulatedSpeed := int64(1000000.0 / float64(server.Ping) * 10) // Simplified simulation
	
	return simulatedSpeed, nil
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