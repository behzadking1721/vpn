package managers

import (
	"c:/Users/behza/OneDrive/Documents/vvpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"os"
	"testing"
)

func TestServerManagerAddServer(t *testing.T) {
	// Create a temporary data manager for testing
	dataManager := NewDataManager("./test_servers.json", "./test_subscriptions.json")
	defer os.Remove("./test_servers.json")
	defer os.Remove("./test_subscriptions.json")

	// Create server manager with test data manager
	serverManager := NewServerManagerWithDataManager(dataManager)

	// Create a test server
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test Server",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolVMess,
		Enabled:  true,
	}

	// Test adding server
	err := serverManager.AddServer(server)
	if err != nil {
		t.Errorf("Failed to add server: %v", err)
	}

	// Check if server was added
	servers := serverManager.GetAllServers()
	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
	}

	// Try to add the same server again (should fail)
	err = serverManager.AddServer(server)
	if err == nil {
		t.Error("Expected error when adding server with duplicate ID")
	}
}

func TestServerManagerRemoveServer(t *testing.T) {
	// Create a temporary data manager for testing
	dataManager := NewDataManager("./test_servers.json", "./test_subscriptions.json")
	defer os.Remove("./test_servers.json")
	defer os.Remove("./test_subscriptions.json")

	// Create server manager with test data manager
	serverManager := NewServerManagerWithDataManager(dataManager)

	// Create a test server
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test Server",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolVMess,
		Enabled:  true,
	}

	// Add server
	err := serverManager.AddServer(server)
	if err != nil {
		t.Errorf("Failed to add server: %v", err)
	}

	// Test removing server
	err = serverManager.RemoveServer(server.ID)
	if err != nil {
		t.Errorf("Failed to remove server: %v", err)
	}

	// Check if server was removed
	servers := serverManager.GetAllServers()
	if len(servers) != 0 {
		t.Errorf("Expected 0 servers, got %d", len(servers))
	}

	// Try to remove non-existent server (should fail)
	err = serverManager.RemoveServer("non-existent-id")
	if err == nil {
		t.Error("Expected error when removing non-existent server")
	}
}

func TestServerManagerUpdateServer(t *testing.T) {
	// Create a temporary data manager for testing
	dataManager := NewDataManager("./test_servers.json", "./test_subscriptions.json")
	defer os.Remove("./test_servers.json")
	defer os.Remove("./test_subscriptions.json")

	// Create server manager with test data manager
	serverManager := NewServerManagerWithDataManager(dataManager)

	// Create a test server
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test Server",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolVMess,
		Enabled:  true,
	}

	// Add server
	err := serverManager.AddServer(server)
	if err != nil {
		t.Errorf("Failed to add server: %v", err)
	}

	// Update server
	updatedServer := server
	updatedServer.Name = "Updated Test Server"
	updatedServer.Port = 8080

	err = serverManager.UpdateServer(updatedServer)
	if err != nil {
		t.Errorf("Failed to update server: %v", err)
	}

	// Check if server was updated
	retrievedServer, err := serverManager.GetServer(server.ID)
	if err != nil {
		t.Errorf("Failed to get server: %v", err)
	}

	if retrievedServer.Name != "Updated Test Server" {
		t.Errorf("Expected server name to be 'Updated Test Server', got '%s'", retrievedServer.Name)
	}

	if retrievedServer.Port != 8080 {
		t.Errorf("Expected server port to be 8080, got %d", retrievedServer.Port)
	}

	// Try to update non-existent server (should fail)
	nonExistentServer := core.Server{
		ID:   "non-existent-id",
		Name: "Non-existent Server",
	}

	err = serverManager.UpdateServer(nonExistentServer)
	if err == nil {
		t.Error("Expected error when updating non-existent server")
	}
}

func TestServerManagerGetServer(t *testing.T) {
	// Create a temporary data manager for testing
	dataManager := NewDataManager("./test_servers.json", "./test_subscriptions.json")
	defer os.Remove("./test_servers.json")
	defer os.Remove("./test_subscriptions.json")

	// Create server manager with test data manager
	serverManager := NewServerManagerWithDataManager(dataManager)

	// Create a test server
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test Server",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolVMess,
		Enabled:  true,
	}

	// Add server
	err := serverManager.AddServer(server)
	if err != nil {
		t.Errorf("Failed to add server: %v", err)
	}

	// Test getting server
	retrievedServer, err := serverManager.GetServer(server.ID)
	if err != nil {
		t.Errorf("Failed to get server: %v", err)
	}

	if retrievedServer.ID != server.ID {
		t.Errorf("Expected server ID to be %s, got %s", server.ID, retrievedServer.ID)
	}

	// Try to get non-existent server (should fail)
	_, err = serverManager.GetServer("non-existent-id")
	if err == nil {
		t.Error("Expected error when getting non-existent server")
	}
}

func TestServerManagerGetAllServers(t *testing.T) {
	// Create a temporary data manager for testing
	dataManager := NewDataManager("./test_servers.json", "./test_subscriptions.json")
	defer os.Remove("./test_servers.json")
	defer os.Remove("./test_subscriptions.json")

	// Create server manager with test data manager
	serverManager := NewServerManagerWithDataManager(dataManager)

	// Initially should have 0 servers
	servers := serverManager.GetAllServers()
	if len(servers) != 0 {
		t.Errorf("Expected 0 servers initially, got %d", len(servers))
	}

	// Add a few servers
	for i := 0; i < 3; i++ {
		server := core.Server{
			ID:       utils.GenerateID(),
			Name:     "Test Server " + string(rune(i+'0')),
			Host:     "example" + string(rune(i+'0')) + ".com",
			Port:     443 + i,
			Protocol: core.ProtocolVMess,
			Enabled:  true,
		}

		err := serverManager.AddServer(server)
		if err != nil {
			t.Errorf("Failed to add server: %v", err)
		}
	}

	// Check if all servers were added
	servers = serverManager.GetAllServers()
	if len(servers) != 3 {
		t.Errorf("Expected 3 servers, got %d", len(servers))
	}
}