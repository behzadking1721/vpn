package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"os"
	"testing"
)

func TestDataManagerSaveAndLoadServers(t *testing.T) {
	// Create temporary file paths
	serversFile := "./test_servers.json"
	subscriptionsFile := "./test_subscriptions.json"
	defer os.Remove(serversFile)
	defer os.Remove(subscriptionsFile)

	// Create data manager
	dataManager := NewDataManager(serversFile, subscriptionsFile)

	// Create test servers
	servers := []core.Server{
		{
			ID:       utils.GenerateID(),
			Name:     "Test Server 1",
			Host:     "example1.com",
			Port:     443,
			Protocol: core.ProtocolVMess,
			Enabled:  true,
		},
		{
			ID:       utils.GenerateID(),
			Name:     "Test Server 2",
			Host:     "example2.com",
			Port:     8080,
			Protocol: core.ProtocolShadowsocks,
			Method:   "aes-256-gcm",
			Password: "test-password",
			Enabled:  true,
		},
	}

	// Test saving servers
	err := dataManager.SaveServers(servers)
	if err != nil {
		t.Errorf("Failed to save servers: %v", err)
	}

	// Test loading servers
	loadedServers, err := dataManager.LoadServers()
	if err != nil {
		t.Errorf("Failed to load servers: %v", err)
	}

	// Check if servers were loaded correctly
	if len(loadedServers) != len(servers) {
		t.Errorf("Expected %d servers, got %d", len(servers), len(loadedServers))
	}

	// Check server details
	for i, server := range servers {
		loadedServer := loadedServers[i]
		if loadedServer.ID != server.ID {
			t.Errorf("Expected server ID %s, got %s", server.ID, loadedServer.ID)
		}

		if loadedServer.Name != server.Name {
			t.Errorf("Expected server name %s, got %s", server.Name, loadedServer.Name)
		}

		if loadedServer.Host != server.Host {
			t.Errorf("Expected server host %s, got %s", server.Host, loadedServer.Host)
		}

		if loadedServer.Port != server.Port {
			t.Errorf("Expected server port %d, got %d", server.Port, loadedServer.Port)
		}

		if loadedServer.Protocol != server.Protocol {
			t.Errorf("Expected server protocol %s, got %s", server.Protocol, loadedServer.Protocol)
		}
	}
}

func TestDataManagerSaveAndLoadSubscriptions(t *testing.T) {
	// Create temporary file paths
	serversFile := "./test_servers.json"
	subscriptionsFile := "./test_subscriptions.json"
	defer os.Remove(serversFile)
	defer os.Remove(subscriptionsFile)

	// Create data manager
	dataManager := NewDataManager(serversFile, subscriptionsFile)

	// Create test subscriptions
	subscriptions := []core.Subscription{
		{
			ID:         utils.GenerateID(),
			Name:       "Test Subscription 1",
			URL:        "https://example1.com/sub",
			AutoUpdate: true,
		},
		{
			ID:         utils.GenerateID(),
			Name:       "Test Subscription 2",
			URL:        "https://example2.com/sub",
			AutoUpdate: false,
		},
	}

	// Test saving subscriptions
	err := dataManager.SaveSubscriptions(subscriptions)
	if err != nil {
		t.Errorf("Failed to save subscriptions: %v", err)
	}

	// Test loading subscriptions
	loadedSubscriptions, err := dataManager.LoadSubscriptions()
	if err != nil {
		t.Errorf("Failed to load subscriptions: %v", err)
	}

	// Check if subscriptions were loaded correctly
	if len(loadedSubscriptions) != len(subscriptions) {
		t.Errorf("Expected %d subscriptions, got %d", len(subscriptions), len(loadedSubscriptions))
	}

	// Check subscription details
	for i, subscription := range subscriptions {
		loadedSubscription := loadedSubscriptions[i]
		if loadedSubscription.ID != subscription.ID {
			t.Errorf("Expected subscription ID %s, got %s", subscription.ID, loadedSubscription.ID)
		}

		if loadedSubscription.Name != subscription.Name {
			t.Errorf("Expected subscription name %s, got %s", subscription.Name, loadedSubscription.Name)
		}

		if loadedSubscription.URL != subscription.URL {
			t.Errorf("Expected subscription URL %s, got %s", subscription.URL, loadedSubscription.URL)
		}

		if loadedSubscription.AutoUpdate != subscription.AutoUpdate {
			t.Errorf("Expected subscription autoUpdate %t, got %t", subscription.AutoUpdate, loadedSubscription.AutoUpdate)
		}
	}
}

func TestDataManagerLoadNonExistentFiles(t *testing.T) {
	// Create temporary file paths for non-existent files
	serversFile := "./nonexistent_servers.json"
	subscriptionsFile := "./nonexistent_subscriptions.json"

	// Create data manager
	dataManager := NewDataManager(serversFile, subscriptionsFile)

	// Test loading servers from non-existent file
	servers, err := dataManager.LoadServers()
	if err != nil {
		t.Errorf("Failed to load servers from non-existent file: %v", err)
	}

	// Should return empty slice
	if len(servers) != 0 {
		t.Errorf("Expected empty servers slice, got %d servers", len(servers))
	}

	// Test loading subscriptions from non-existent file
	subscriptions, err := dataManager.LoadSubscriptions()
	if err != nil {
		t.Errorf("Failed to load subscriptions from non-existent file: %v", err)
	}

	// Should return empty slice
	if len(subscriptions) != 0 {
		t.Errorf("Expected empty subscriptions slice, got %d subscriptions", len(subscriptions))
	}
}