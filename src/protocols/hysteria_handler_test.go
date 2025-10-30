package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"testing"
	"time"
)

func TestHysteriaHandler(t *testing.T) {
	// Create Hysteria handler
	handler := NewHysteriaHandler()

	// Test GetProtocol
	protocol := handler.GetProtocol()
	if protocol != core.ProtocolHysteria {
		t.Errorf("Expected protocol Hysteria, got %s", protocol)
	}

	// Test IsConnected (should be false initially)
	connected := handler.IsConnected()
	if connected {
		t.Error("Expected handler to be disconnected initially")
	}

	// Create a test server
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test Hysteria Server",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolHysteria,
		Password: "test-password",
		UpMbps:   50,
		DownMbps: 100,
		Obfs:     "salamander",
		AuthStr:  "auth-string",
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

	if details["protocol"] != core.ProtocolHysteria {
		t.Errorf("Expected protocol Hysteria in details, got %v", details["protocol"])
	}

	if details["host"] != "example.com" {
		t.Errorf("Expected host example.com in details, got %v", details["host"])
	}

	if details["up_mbps"] != 50 {
		t.Errorf("Expected up_mbps 50 in details, got %v", details["up_mbps"])
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

func TestHysteriaHandlerMissingPassword(t *testing.T) {
	// Create Hysteria handler
	handler := NewHysteriaHandler()

	// Create a test server without password
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test Hysteria Server (Missing Password)",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolHysteria,
		Enabled:  true,
	}

	// Test Connect should fail
	err := handler.Connect(server)
	if err == nil {
		t.Error("Expected error when connecting with missing password")
	}

	if err != nil && err.Error() != "missing password for Hysteria protocol" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}