package database

import (
	"fmt"
)

// DB is an interface for database operations
// For now, we use JSONStore as a temporary solution
// In the future, this can be replaced with SQLite or other database implementations
type DB interface {
	Close() error
}

// Store is a unified interface for all data store operations
// JSONStore implements this interface directly
type Store interface {
	DB
	
	// Server operations (type assertion needed to *core.Server)
	AddServer(server interface{}) error
	GetServer(id string) (interface{}, error)
	GetAllServers() ([]interface{}, error)
	GetEnabledServers() ([]interface{}, error)
	UpdateServer(server interface{}) error
	DeleteServer(id string) error
	UpdatePing(id string, ping int) error
	
	// Subscription operations (type assertion needed to *core.Subscription)
	AddSubscription(sub interface{}) error
	GetSubscription(id string) (interface{}, error)
	GetAllSubscriptions() ([]interface{}, error)
	UpdateSubscription(sub interface{}) error
	DeleteSubscription(id string) error
}

// NewDB creates a new database/store connection
// Currently uses JSONStore, but can be switched to SQLite later
func NewDB(dataDir string) (Store, error) {
	store, err := NewJSONStore(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}
	return store, nil
}

