package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"errors"
	"sync"
	"time"
)

// ServerManager handles server operations
type ServerManager struct {
	servers       []core.Server
	subscriptions []core.Subscription
	mutex         sync.RWMutex
}

// NewServerManager creates a new server manager
func NewServerManager() *ServerManager {
	return &ServerManager{
		servers:       make([]core.Server, 0),
		subscriptions: make([]core.Subscription, 0),
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

// FindFastestServer finds the server with the lowest ping
func (sm *ServerManager) FindFastestServer() (core.Server, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if len(sm.servers) == 0 {
		return core.Server{}, errors.New("no servers available")
	}

	var fastestServer core.Server
	lowestPing := -1

	for _, server := range sm.servers {
		// Only consider enabled servers
		if !server.Enabled {
			continue
		}

		// If this is the first enabled server or has a lower ping
		if lowestPing == -1 || (server.Ping < lowestPing && server.Ping > 0) {
			fastestServer = server
			lowestPing = server.Ping
		}
	}

	if lowestPing == -1 {
		return core.Server{}, errors.New("no enabled servers available")
	}

	return fastestServer, nil
}

// UpdateServerPing updates the ping value for a server
func (sm *ServerManager) UpdateServerPing(serverID string, ping int) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for i, server := range sm.servers {
		if server.ID == serverID {
			sm.servers[i].Ping = ping
			sm.servers[i].LastPing = time.Now()
			return nil
		}
	}

	return errors.New("server not found")
}