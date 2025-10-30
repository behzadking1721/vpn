package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"testing"
	"time"
)

func TestSSHHandlerWithPassword(t *testing.T) {
	// Create SSH handler
	handler := NewSSHHandler()

	// Test GetProtocol
	protocol := handler.GetProtocol()
	if protocol != core.ProtocolSSH {
		t.Errorf("Expected protocol SSH, got %s", protocol)
	}

	// Test IsConnected (should be false initially)
	connected := handler.IsConnected()
	if connected {
		t.Error("Expected handler to be disconnected initially")
	}

	// Create a test server with password authentication
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test SSH Server",
		Host:     "example.com",
		Port:     22,
		Protocol: core.ProtocolSSH,
		Username: "testuser",
		Password: "testpassword",
		Enabled:  true,
	}

	// Test Connect
	err := handler.Connect(server)
	if err != nil {
		t.Errorf("Failed to connect: %v", err)
	}

	// Test IsConnected (should be true after connecting)
	connected = handler.IsConnected()
	if !connected {
		t.Error("Expected handler to be connected after Connect()")
	}

	// Give some time for the data usage simulation to run
	time.Sleep(2 * time.Second)

	// Test GetDataUsage
	sent, received, err := handler.GetDataUsage()
	if err != nil {
		t.Errorf("Failed to get data usage: %v", err)
	}

	if sent <= 0 {
		t.Error("Expected sent data to be greater than 0")
	}

	if received <= 0 {
		t.Error("Expected received data to be greater than 0")
	}

	// Test GetConnectionDetails
	details, err := handler.GetConnectionDetails()
	if err != nil {
		t.Errorf("Failed to get connection details: %v", err)
	}

	if details["protocol"] != core.ProtocolSSH {
		t.Errorf("Expected protocol SSH in details, got %v", details["protocol"])
	}

	if details["host"] != "example.com" {
		t.Errorf("Expected host example.com in details, got %v", details["host"])
	}

	if details["username"] != "testuser" {
		t.Errorf("Expected username testuser in details, got %v", details["username"])
	}

	// Test Disconnect
	err = handler.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect: %v", err)
	}

	// Test IsConnected (should be false after disconnecting)
	connected = handler.IsConnected()
	if connected {
		t.Error("Expected handler to be disconnected after Disconnect()")
	}
}

func TestSSHHandlerWithPrivateKey(t *testing.T) {
	// Create SSH handler
	handler := NewSSHHandler()

	// Create a test server with private key authentication
	server := core.Server{
		ID:        utils.GenerateID(),
		Name:      "Test SSH Server (Private Key)",
		Host:      "example.com",
		Port:      22,
		Protocol:  core.ProtocolSSH,
		Username:  "testuser",
		PrivateKey: "-----BEGIN OPENSSH PRIVATE KEY-----\ntest-key-content\n-----END OPENSSH PRIVATE KEY-----",
		Enabled:   true,
	}

	// Test Connect
	err := handler.Connect(server)
	if err != nil {
		t.Errorf("Failed to connect: %v", err)
	}

	// Test IsConnected (should be true after connecting)
	connected := handler.IsConnected()
	if !connected {
		t.Error("Expected handler to be connected after Connect()")
	}

	// Test GetConnectionDetails
	details, err := handler.GetConnectionDetails()
	if err != nil {
		t.Errorf("Failed to get connection details: %v", err)
	}

	if details["private_key"] != true {
		t.Errorf("Expected private_key to be true in details, got %v", details["private_key"])
	}

	// Test Disconnect
	err = handler.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect: %v", err)
	}
}

func TestSSHHandlerMissingAuth(t *testing.T) {
	// Create SSH handler
	handler := NewSSHHandler()

	// Create a test server without authentication
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test SSH Server (Missing Auth)",
		Host:     "example.com",
		Port:     22,
		Protocol: core.ProtocolSSH,
		Username: "testuser",
		Enabled:  true,
		// Missing Password and PrivateKey
	}

	// Test Connect should fail
	err := handler.Connect(server)
	if err == nil {
		t.Error("Expected error when connecting with missing authentication")
	}

	if err != nil && err.Error() != "missing password or private key for SSH authentication" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestSSHHandlerMissingUsername(t *testing.T) {
	// Create SSH handler
	handler := NewSSHHandler()

	// Create a test server without username
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test SSH Server (Missing Username)",
		Host:     "example.com",
		Port:     22,
		Protocol: core.ProtocolSSH,
		Password: "testpassword",
		Enabled:  true,
		// Missing Username
	}

	// Test Connect should fail
	err := handler.Connect(server)
	if err == nil {
		t.Error("Expected error when connecting with missing username")
	}

	if err != nil && err.Error() != "missing username for SSH protocol" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}