package database

import (
	"vpnclient/src/core"
)

// ServerStore handles server database operations
// For now, JSONStore implements these methods directly
// This file is kept for future SQLite implementation

// NewServerStore creates a new server store
// The JSONStore already implements all server operations directly
func NewServerStore(store Store) *ServerStoreWrapper {
	return &ServerStoreWrapper{store: store}
}

// ServerStoreWrapper wraps Store to provide ServerStore interface
type ServerStoreWrapper struct {
	store Store
}

// AddServer adds a new server to the database
func (s *ServerStoreWrapper) AddServer(server *core.Server) error {
	return s.store.AddServer(server)
}

// GetServer retrieves a server by ID
func (s *ServerStoreWrapper) GetServer(id string) (*core.Server, error) {
	result, err := s.store.GetServer(id)
	if err != nil {
		return nil, err
	}
	return result.(*core.Server), nil
}

// GetAllServers retrieves all servers
func (s *ServerStoreWrapper) GetAllServers() ([]*core.Server, error) {
	results, err := s.store.GetAllServers()
	if err != nil {
		return nil, err
	}
	servers := make([]*core.Server, len(results))
	for i, r := range results {
		servers[i] = r.(*core.Server)
	}
	return servers, nil
}

// GetEnabledServers retrieves only enabled servers
func (s *ServerStoreWrapper) GetEnabledServers() ([]*core.Server, error) {
	results, err := s.store.GetEnabledServers()
	if err != nil {
		return nil, err
	}
	servers := make([]*core.Server, len(results))
	for i, r := range results {
		servers[i] = r.(*core.Server)
	}
	return servers, nil
}

// UpdateServer updates an existing server
func (s *ServerStoreWrapper) UpdateServer(server *core.Server) error {
	return s.store.UpdateServer(server)
}

// DeleteServer deletes a server by ID
func (s *ServerStoreWrapper) DeleteServer(id string) error {
	return s.store.DeleteServer(id)
}

// UpdatePing updates the ping value for a server
func (s *ServerStoreWrapper) UpdatePing(id string, ping int) error {
	return s.store.UpdatePing(id, ping)
}
