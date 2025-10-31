package managers

import (
	"vpn-client/src/core"
	"errors"
	"sync"
	"time"
)

// ServerManager handles server operations
type ServerManager struct {
	servers       []core.Server
	subscriptions []core.Subscription
	dataManager   *DataManager
	mutex         sync.RWMutex
}

// NewServerManager creates a new server manager
func NewServerManager() *ServerManager {
	// Create data manager with default file paths
	dataManager := NewDataManager("./data/servers.json", "./data/subscriptions.json")

	// Create server manager
	sm := &ServerManager{
		servers:       make([]core.Server, 0),
		subscriptions: make([]core.Subscription, 0),
		dataManager:   dataManager,
	}

	// Load existing data
	sm.loadExistingData()

	return sm
}

// NewServerManagerWithDataManager creates a new server manager with a specific data manager
func NewServerManagerWithDataManager(dataManager *DataManager) *ServerManager {
	sm := &ServerManager{
		servers:       make([]core.Server, 0),
		subscriptions: make([]core.Subscription, 0),
		dataManager:   dataManager,
	}

	// Load existing data
	sm.loadExistingData()

	return sm
}

// loadExistingData loads existing servers and subscriptions from files
func (sm *ServerManager) loadExistingData() {
	// Load servers
	servers, err := sm.dataManager.LoadServers()
	if err != nil {
		// If there's an error, we continue with empty lists
		// In a production environment, you might want to log this error
	} else {
		sm.servers = servers
	}

	// Load subscriptions
	subscriptions, err := sm.dataManager.LoadSubscriptions()
	if err != nil {
		// If there's an error, we continue with empty lists
		// In a production environment, you might want to log this error
	} else {
		sm.subscriptions = subscriptions
	}
}

// AddServer adds a new server
func (sm *ServerManager) AddServer(server core.Server) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if server with same ID already exists
	for _, s := range sm.servers {
		if s.ID == server.ID {
			return errors.New("server with this ID already exists")
		}
	}

	sm.servers = append(sm.servers, server)

	// Save to file
	sm.saveServers()

	return nil
}

// RemoveServer removes a server by ID
func (sm *ServerManager) RemoveServer(serverID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for i, server := range sm.servers {
		if server.ID == serverID {
			// Remove the server
			sm.servers = append(sm.servers[:i], sm.servers[i+1:]...)

			// Save to file
			sm.saveServers()

			return nil
		}
	}

	return errors.New("server not found")
}

// UpdateServer updates an existing server
func (sm *ServerManager) UpdateServer(server core.Server) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for i, s := range sm.servers {
		if s.ID == server.ID {
			sm.servers[i] = server

			// Save to file
			sm.saveServers()

			return nil
		}
	}

	return errors.New("server not found")
}

// GetServer returns a server by ID
func (sm *ServerManager) GetServer(serverID string) (core.Server, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for _, server := range sm.servers {
		if server.ID == serverID {
			return server, nil
		}
	}

	return core.Server{}, errors.New("server not found")
}

// GetAllServers returns all servers
func (sm *ServerManager) GetAllServers() []core.Server {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Return a copy to prevent external modification
	servers := make([]core.Server, len(sm.servers))
	copy(servers, sm.servers)
	return servers
}

// AddSubscription adds a new subscription
func (sm *ServerManager) AddSubscription(sub core.Subscription) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if subscription with same ID already exists
	for _, s := range sm.subscriptions {
		if s.ID == sub.ID {
			return errors.New("subscription with this ID already exists")
		}
	}

	sm.subscriptions = append(sm.subscriptions, sub)

	// Save to file
	sm.saveSubscriptions()

	return nil
}

// RemoveSubscription removes a subscription by ID
func (sm *ServerManager) RemoveSubscription(subID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for i, sub := range sm.subscriptions {
		if sub.ID == subID {
			// Remove the subscription
			sm.subscriptions = append(sm.subscriptions[:i], sm.subscriptions[i+1:]...)

			// Save to file
			sm.saveSubscriptions()

			return nil
		}
	}

	return errors.New("subscription not found")
}

// UpdateSubscription updates an existing subscription
func (sm *ServerManager) UpdateSubscription(sub core.Subscription) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for i, s := range sm.subscriptions {
		if s.ID == sub.ID {
			sm.subscriptions[i] = sub

			// Save to file
			sm.saveSubscriptions()

			return nil
		}
	}

	return errors.New("subscription not found")
}

// GetSubscription returns a subscription by ID
func (sm *ServerManager) GetSubscription(subID string) (core.Subscription, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for _, sub := range sm.subscriptions {
		if sub.ID == subID {
			return sub, nil
		}
	}

	return core.Subscription{}, errors.New("subscription not found")
}

// GetAllSubscriptions returns all subscriptions
func (sm *ServerManager) GetAllSubscriptions() []core.Subscription {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Return a copy to prevent external modification
	subs := make([]core.Subscription, len(sm.subscriptions))
	copy(subs, sm.subscriptions)
	return subs
}

// saveServers saves servers to file
func (sm *ServerManager) saveServers() {
	// Make a copy of servers to avoid holding the lock during I/O
	servers := make([]core.Server, len(sm.servers))
	copy(servers, sm.servers)

	// Save in a separate goroutine to avoid blocking
	go func() {
		err := sm.dataManager.SaveServers(servers)
		if err != nil {
			// In a production environment, you might want to log this error
		}
	}()
}

// saveSubscriptions saves subscriptions to file
func (sm *ServerManager) saveSubscriptions() {
	// Make a copy of subscriptions to avoid holding the lock during I/O
	subs := make([]core.Subscription, len(sm.subscriptions))
	copy(subs, sm.subscriptions)

	// Save in a separate goroutine to avoid blocking
	go func() {
		err := sm.dataManager.SaveSubscriptions(subs)
		if err != nil {
			// In a production environment, you might want to log this error
		}
	}()
}
