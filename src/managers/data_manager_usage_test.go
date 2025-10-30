package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"os"
	"testing"
)

func TestDataUsageRecording(t *testing.T) {
	// Create temporary file paths
	serversFile := "./test_servers.json"
	subscriptionsFile := "./test_subscriptions.json"
	defer os.Remove(serversFile)
	defer os.Remove(subscriptionsFile)

	// Create data manager
	dataManager := NewDataManager(serversFile, subscriptionsFile)

	// Create a test server
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test Server",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolVMess,
		Enabled:  true,
	}

	// Save the server
	err := dataManager.SaveServers([]core.Server{server})
	if err != nil {
		t.Errorf("Failed to save server: %v", err)
	}

	// Record some data usage
	sent := int64(1024 * 1024) // 1 MB
	received := int64(2048 * 1024) // 2 MB
	total := sent + received

	// Record data usage multiple times to test accumulation
	dataManager.RecordDataUsage(server.ID, sent, received)
	dataManager.RecordDataUsage(server.ID, sent, received)

	// Load servers and check data usage
	servers, err := dataManager.LoadServers()
	if err != nil {
		t.Errorf("Failed to load servers: %v", err)
	}

	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
	}

	loadedServer := servers[0]
	expectedUsage := total * 2 // Because we recorded twice
	if loadedServer.DataUsed != expectedUsage {
		t.Errorf("Expected data usage %d, got %d", expectedUsage, loadedServer.DataUsed)
	}
}

func TestDataLimitEnforcement(t *testing.T) {
	// Create temporary file paths
	serversFile := "./test_servers.json"
	subscriptionsFile := "./test_subscriptions.json"
	defer os.Remove(serversFile)
	defer os.Remove(subscriptionsFile)

	// Create data manager
	dataManager := NewDataManager(serversFile, subscriptionsFile)

	// Create a test server with data limit
	dataLimit := int64(5 * 1024 * 1024) // 5 MB
	server := core.Server{
		ID:        utils.GenerateID(),
		Name:      "Test Server with Limit",
		Host:      "example.com",
		Port:      443,
		Protocol:  core.ProtocolVMess,
		DataLimit: dataLimit,
		Enabled:   true,
	}

	// Save the server
	err := dataManager.SaveServers([]core.Server{server})
	if err != nil {
		t.Errorf("Failed to save server: %v", err)
	}

	// Record data usage that exceeds the limit
	sent := int64(3 * 1024 * 1024) // 3 MB
	received := int64(3 * 1024 * 1024) // 3 MB
	total := sent + received // 6 MB, which exceeds the 5 MB limit

	dataManager.RecordDataUsage(server.ID, sent, received)

	// Load servers and check data usage
	servers, err := dataManager.LoadServers()
	if err != nil {
		t.Errorf("Failed to load servers: %v", err)
	}

	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
	}

	loadedServer := servers[0]
	if loadedServer.DataUsed != total {
		t.Errorf("Expected data usage %d, got %d", total, loadedServer.DataUsed)
	}

	// Check if data limit is exceeded
	if loadedServer.DataLimit > 0 && loadedServer.DataUsed >= loadedServer.DataLimit {
		// This is expected - the limit is exceeded
		if loadedServer.DataUsed < loadedServer.DataLimit {
			t.Errorf("Expected data usage %d to exceed limit %d", loadedServer.DataUsed, loadedServer.DataLimit)
		}
	} else {
		t.Error("Expected data limit to be exceeded")
	}
}

func TestDataUsageReset(t *testing.T) {
	// Create temporary file paths
	serversFile := "./test_servers.json"
	subscriptionsFile := "./test_subscriptions.json"
	defer os.Remove(serversFile)
	defer os.Remove(subscriptionsFile)

	// Create data manager
	dataManager := NewDataManager(serversFile, subscriptionsFile)

	// Create a test server
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test Server",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolVMess,
		Enabled:  true,
	}

	// Record some data usage
	sent := int64(1024 * 1024) // 1 MB
	received := int64(2048 * 1024) // 2 MB

	dataManager.RecordDataUsage(server.ID, sent, received)

	// Reset data usage
	dataManager.ResetDataUsage(server.ID)

	// Load servers and check data usage is reset
	servers, err := dataManager.LoadServers()
	if err != nil {
		t.Errorf("Failed to load servers: %v", err)
	}

	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
	}

	loadedServer := servers[0]
	if loadedServer.DataUsed != 0 {
		t.Errorf("Expected data usage to be 0 after reset, got %d", loadedServer.DataUsed)
	}
}

func TestGetAllDataUsage(t *testing.T) {
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
			Enabled:  true,
		},
	}

	// Save the servers
	err := dataManager.SaveServers(servers)
	if err != nil {
		t.Errorf("Failed to save servers: %v", err)
	}

	// Record data usage for both servers
	dataManager.RecordDataUsage(servers[0].ID, 1024*1024, 2048*1024) // 1MB sent, 2MB received
	dataManager.RecordDataUsage(servers[1].ID, 2048*1024, 4096*1024) // 2MB sent, 4MB received

	// Get all data usage
	allData := dataManager.GetAllData()

	// Check we have data for both servers
	if len(allData) != 2 {
		t.Errorf("Expected data for 2 servers, got %d", len(allData))
	}

	// Check data for first server
	server1Data, exists := allData[servers[0].ID]
	if !exists {
		t.Error("Expected data for server 1")
	} else {
		expectedTotal := int64(1024*1024 + 2048*1024) // 3MB
		if server1Data.TotalSent+server1Data.TotalRecv != expectedTotal {
			t.Errorf("Expected total data %d for server 1, got %d", expectedTotal, server1Data.TotalSent+server1Data.TotalRecv)
		}
	}

	// Check data for second server
	server2Data, exists := allData[servers[1].ID]
	if !exists {
		t.Error("Expected data for server 2")
	} else {
		expectedTotal := int64(2048*1024 + 4096*1024) // 6MB
		if server2Data.TotalSent+server2Data.TotalRecv != expectedTotal {
			t.Errorf("Expected total data %d for server 2, got %d", expectedTotal, server2Data.TotalSent+server2Data.TotalRecv)
		}
	}
}