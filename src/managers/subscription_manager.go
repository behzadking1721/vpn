package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

// SubscriptionManager handles subscription operations
type SubscriptionManager struct {
	serverManager *ServerManager
	mutex         sync.RWMutex
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager(serverManager *ServerManager) *SubscriptionManager {
	return &SubscriptionManager{
		serverManager: serverManager,
	}
}

// AddSubscription adds a new subscription and fetches servers
func (sm *SubscriptionManager) AddSubscription(sub core.Subscription) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Parse subscription URL and fetch servers
	servers, err := sm.parseSubscription(sub.URL)
	if err != nil {
		return errors.New("failed to parse subscription: " + err.Error())
	}

	// Update subscription with parsed servers
	sub.Servers = servers
	sub.LastUpdate = time.Now()

	// Add to server manager
	return sm.serverManager.AddSubscription(sub)
}

// UpdateSubscription updates an existing subscription
func (sm *SubscriptionManager) UpdateSubscription(subID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Get existing subscription
	sub, err := sm.serverManager.GetSubscription(subID)
	if err != nil {
		return err
	}

	// Parse subscription URL and fetch servers
	servers, err := sm.parseSubscription(sub.URL)
	if err != nil {
		return errors.New("failed to parse subscription: " + err.Error())
	}

	// Update subscription with parsed servers
	sub.Servers = servers
	sub.LastUpdate = time.Now()

	// Update in server manager
	return sm.serverManager.UpdateSubscription(sub)
}

// RemoveSubscription removes a subscription
func (sm *SubscriptionManager) RemoveSubscription(subID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	return sm.serverManager.RemoveSubscription(subID)
}

// parseSubscription parses a subscription URL and returns servers
func (sm *SubscriptionManager) parseSubscription(url string) ([]core.Server, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Make HTTP request to fetch subscription content
	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.New("failed to fetch subscription: " + err.Error())
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("subscription fetch failed with status: " + resp.Status)
	}
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read subscription content: " + err.Error())
	}
	
	// Use existing subscription parser to parse the content
	parser := NewSubscriptionParser()
	servers, err := parser.ParseSubscriptionLink(string(body))
	if err != nil {
		return nil, errors.New("failed to parse subscription content: " + err.Error())
	}
	
	return servers, nil
}

// ImportFromQRCode imports servers from a QR code
func (sm *SubscriptionManager) ImportFromQRCode(qrContent string) ([]core.Server, error) {
	// Use the existing subscription parser to parse QR code content
	// The QR code should contain a subscription link or configuration
	parser := NewSubscriptionParser()
	servers, err := parser.ParseSubscriptionLink(qrContent)
	if err != nil {
		return nil, errors.New("failed to parse QR code content: " + err.Error())
	}
	
	return servers, nil
}