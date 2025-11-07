package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"vpnclient/src/core"
)

// JSONStore is a simple JSON-based file storage (temporary solution until SQLite is available)
// It implements Store interface and provides all CRUD operations
type JSONStore struct {
	dataDir string
	servers []*core.Server
	subs    []*core.Subscription
	mutex   sync.RWMutex
}

// Ensure JSONStore implements Store interface
var _ Store = (*JSONStore)(nil)

// NewJSONStore creates a new JSON-based store
func NewJSONStore(dataDir string) (*JSONStore, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	store := &JSONStore{
		dataDir: dataDir,
		servers: make([]*core.Server, 0),
		subs:    make([]*core.Subscription, 0),
	}

	// Load existing data
	if err := store.loadData(); err != nil {
		// If load fails, start with empty data
		store.servers = make([]*core.Server, 0)
		store.subs = make([]*core.Subscription, 0)
	}

	return store, nil
}

// loadData loads data from JSON files
func (s *JSONStore) loadData() error {
	serversPath := filepath.Join(s.dataDir, "servers.json")
	subsPath := filepath.Join(s.dataDir, "subscriptions.json")

	// Load servers
	if data, err := os.ReadFile(serversPath); err == nil {
		if err := json.Unmarshal(data, &s.servers); err != nil {
			return fmt.Errorf("failed to unmarshal servers: %w", err)
		}
	}

	// Load subscriptions
	if data, err := os.ReadFile(subsPath); err == nil {
		if err := json.Unmarshal(data, &s.subs); err != nil {
			return fmt.Errorf("failed to unmarshal subscriptions: %w", err)
		}
	}

	return nil
}

// saveData saves data to JSON files
func (s *JSONStore) saveData() error {
	serversPath := filepath.Join(s.dataDir, "servers.json")
	subsPath := filepath.Join(s.dataDir, "subscriptions.json")

	// Save servers
	serversData, err := json.MarshalIndent(s.servers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal servers: %w", err)
	}
	if err := os.WriteFile(serversPath, serversData, 0644); err != nil {
		return fmt.Errorf("failed to write servers file: %w", err)
	}

	// Save subscriptions
	subsData, err := json.MarshalIndent(s.subs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal subscriptions: %w", err)
	}
	if err := os.WriteFile(subsPath, subsData, 0644); err != nil {
		return fmt.Errorf("failed to write subscriptions file: %w", err)
	}

	return nil
}

// ServerStore methods

// AddServer adds a new server
func (s *JSONStore) AddServer(server interface{}) error {
	svr, ok := server.(*core.Server)
	if !ok {
		return fmt.Errorf("invalid server type")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if svr.ID == "" {
		return fmt.Errorf("server ID is required")
	}

	// Check if server already exists
	for _, existing := range s.servers {
		if existing.ID == svr.ID {
			return fmt.Errorf("server with ID %s already exists", svr.ID)
		}
	}

	now := time.Now()
	if svr.CreatedAt.IsZero() {
		svr.CreatedAt = now
	}
	svr.UpdatedAt = now

	if svr.Config == nil {
		svr.Config = make(map[string]interface{})
	}

	s.servers = append(s.servers, svr)
	return s.saveData()
}

// GetServer retrieves a server by ID
func (s *JSONStore) GetServer(id string) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, server := range s.servers {
		if server.ID == id {
			return server, nil
		}
	}

	return nil, fmt.Errorf("server not found: %s", id)
}

// GetAllServers retrieves all servers
func (s *JSONStore) GetAllServers() ([]interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]interface{}, len(s.servers))
	for i, srv := range s.servers {
		result[i] = srv
	}
	return result, nil
}

// GetEnabledServers retrieves only enabled servers
func (s *JSONStore) GetEnabledServers() ([]interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []interface{}
	for _, server := range s.servers {
		if server.Enabled {
			result = append(result, server)
		}
	}
	return result, nil
}

// UpdateServer updates an existing server
func (s *JSONStore) UpdateServer(server interface{}) error {
	svr, ok := server.(*core.Server)
	if !ok {
		return fmt.Errorf("invalid server type")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if svr.ID == "" {
		return fmt.Errorf("server ID is required")
	}

	for i, existing := range s.servers {
		if existing.ID == svr.ID {
			svr.UpdatedAt = time.Now()
			s.servers[i] = svr
			return s.saveData()
		}
	}

	return fmt.Errorf("server not found: %s", svr.ID)
}

// DeleteServer deletes a server by ID
func (s *JSONStore) DeleteServer(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, server := range s.servers {
		if server.ID == id {
			s.servers = append(s.servers[:i], s.servers[i+1:]...)
			return s.saveData()
		}
	}

	return fmt.Errorf("server not found: %s", id)
}

// UpdatePing updates the ping value for a server
func (s *JSONStore) UpdatePing(id string, ping int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, server := range s.servers {
		if server.ID == id {
			server.Ping = ping
			server.UpdatedAt = time.Now()
			return s.saveData()
		}
	}

	return fmt.Errorf("server not found: %s", id)
}

// SubscriptionStore methods

// AddSubscription adds a new subscription
func (s *JSONStore) AddSubscription(sub interface{}) error {
	subscription, ok := sub.(*core.Subscription)
	if !ok {
		return fmt.Errorf("invalid subscription type")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if subscription.ID == "" {
		return fmt.Errorf("subscription ID is required")
	}

	// Check if subscription already exists
	for _, existing := range s.subs {
		if existing.ID == subscription.ID {
			return fmt.Errorf("subscription with ID %s already exists", subscription.ID)
		}
	}

	now := time.Now()
	if subscription.CreatedAt.IsZero() {
		subscription.CreatedAt = now
	}
	subscription.UpdatedAt = now

	s.subs = append(s.subs, subscription)
	return s.saveData()
}

// GetSubscription retrieves a subscription by ID
func (s *JSONStore) GetSubscription(id string) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, sub := range s.subs {
		if sub.ID == id {
			return sub, nil
		}
	}

	return nil, fmt.Errorf("subscription not found: %s", id)
}

// GetAllSubscriptions retrieves all subscriptions
func (s *JSONStore) GetAllSubscriptions() ([]interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]interface{}, len(s.subs))
	for i, sub := range s.subs {
		result[i] = sub
	}
	return result, nil
}

// UpdateSubscription updates an existing subscription
func (s *JSONStore) UpdateSubscription(sub interface{}) error {
	subscription, ok := sub.(*core.Subscription)
	if !ok {
		return fmt.Errorf("invalid subscription type")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if subscription.ID == "" {
		return fmt.Errorf("subscription ID is required")
	}

	for i, existing := range s.subs {
		if existing.ID == subscription.ID {
			subscription.UpdatedAt = time.Now()
			s.subs[i] = subscription
			return s.saveData()
		}
	}

	return fmt.Errorf("subscription not found: %s", subscription.ID)
}

// DeleteSubscription deletes a subscription by ID
func (s *JSONStore) DeleteSubscription(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, sub := range s.subs {
		if sub.ID == id {
			s.subs = append(s.subs[:i], s.subs[i+1:]...)
			return s.saveData()
		}
	}

	return fmt.Errorf("subscription not found: %s", id)
}

// Close closes the store (no-op for JSON store, but implements interface)
func (s *JSONStore) Close() error {
	return nil
}
