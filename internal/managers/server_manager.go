package managers

import (
	"fmt"
	"strings"

	"vpnclient/internal/database"
	"vpnclient/src/core"
)

// ServerManager manages VPN servers
type ServerManager struct {
	store *database.ServerStoreWrapper
}

// NewServerManager creates a new server manager
func NewServerManager(store database.Store) *ServerManager {
	return &ServerManager{
		store: database.NewServerStore(store),
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
	return sm.store.UpdatePing(id, ping)
}

// GetServersByProtocol retrieves servers by protocol
func (sm *ServerManager) GetServersByProtocol(protocol string) ([]*core.Server, error) {
	all, err := sm.store.GetAllServers()
	if err != nil {
		return nil, err
	}

	var filtered []*core.Server
	for _, server := range all {
		if strings.EqualFold(server.Protocol, protocol) {
			filtered = append(filtered, server)
		}
	}

	return filtered, nil
}

// SearchServers searches servers by name or host
func (sm *ServerManager) SearchServers(query string) ([]*core.Server, error) {
	if query == "" {
		return sm.store.GetAllServers()
	}

	all, err := sm.store.GetAllServers()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var filtered []*core.Server

	for _, server := range all {
		name := strings.ToLower(server.Name)
		host := strings.ToLower(server.Host)
		country := strings.ToLower(server.Country)

		if strings.Contains(name, query) ||
			strings.Contains(host, query) ||
			strings.Contains(country, query) {
			filtered = append(filtered, server)
		}
	}

	return filtered, nil
}

// GetFastestServer returns the server with the lowest ping
func (sm *ServerManager) GetFastestServer() (*core.Server, error) {
	enabled, err := sm.store.GetEnabledServers()
	if err != nil {
		return nil, err
	}

	if len(enabled) == 0 {
		return nil, fmt.Errorf("no enabled servers available")
	}

	fastest := enabled[0]
	for _, server := range enabled {
		if server.Ping > 0 && (fastest.Ping == 0 || server.Ping < fastest.Ping) {
			fastest = server
		}
	}

	return fastest, nil
}

// validateServer validates server data
func (sm *ServerManager) validateServer(server *core.Server) error {
	if server.Name == "" {
		return fmt.Errorf("server name is required")
	}
	if server.Host == "" {
		return fmt.Errorf("server host is required")
	}
	if server.Port <= 0 || server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", server.Port)
	}
	if server.Protocol == "" {
		return fmt.Errorf("server protocol is required")
	}

	// Validate protocol
	validProtocols := map[string]bool{
		"vmess":       true,
		"vless":       true,
		"shadowsocks": true,
		"trojan":      true,
		"reality":     true,
		"hysteria2":   true,
		"tuic":        true,
		"ssh":         true,
	}

	if !validProtocols[strings.ToLower(server.Protocol)] {
		return fmt.Errorf("unsupported protocol: %s", server.Protocol)
	}

	return nil
}

// generateServerID generates a unique ID for a server
func (sm *ServerManager) generateServerID(server *core.Server) string {
	// Simple ID generation: protocol-host-port-hash
	base := fmt.Sprintf("%s-%s-%d", strings.ToLower(server.Protocol), server.Host, server.Port)
	
	// Add a simple hash based on name
	hash := 0
	for _, c := range server.Name {
		hash = hash*31 + int(c)
	}
	if hash < 0 {
		hash = -hash
	}

	return fmt.Sprintf("%s-%d", base, hash)
}

