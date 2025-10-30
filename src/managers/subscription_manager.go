package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"errors"
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
	// This is where you would implement the actual subscription parsing logic
	// depending on the format (e.g., base64 encoded vmess URLs, SIP008 JSON, etc.)
	
	// For now, we'll return an empty slice and simulate parsing
	// In a real implementation, you would:
	// 1. Fetch the content from the URL
	// 2. Decode/parse the content based on the subscription format
	// 3. Convert to our Server format
	
	// Simulate some servers for demonstration
	servers := []core.Server{
		{
			ID:       utils.GenerateID(),
			Name:     "Sample Server 1",
			Host:     "server1.example.com",
			Port:     443,
			Protocol: core.ProtocolVMess,
			Encryption: "auto",
			TLS:      true,
			Remark:   "Sample server for demonstration",
			Enabled:  true,
		},
		{
			ID:       utils.GenerateID(),
			Name:     "Sample Server 2",
			Host:     "server2.example.com",
			Port:     80,
			Protocol: core.ProtocolShadowsocks,
			Method:   "aes-256-gcm",
			Password: "sample-password",
			Remark:   "Shadowsocks server",
			Enabled:  true,
		},
	}
	
	return servers, nil
}

// ImportFromQRCode imports servers from a QR code
func (sm *SubscriptionManager) ImportFromQRCode(qrContent string) ([]core.Server, error) {
	// This is where you would implement QR code parsing
	// For now, we'll just simulate it
	
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "QR Imported Server",
		Host:     "qr-server.example.com",
		Port:     443,
		Protocol: core.ProtocolVMess,
		Encryption: "auto",
		TLS:      true,
		Remark:   "Imported via QR code",
		Enabled:  true,
	}
	
	return []core.Server{server}, nil
}