package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"testing"
	"time"
)

func TestTUICHandler(t *testing.T) {
	// Create TUIC handler
	handler := NewTUICHandler()

	// Test GetProtocol
	protocol := handler.GetProtocol()
	if protocol != core.ProtocolTUIC {
		t.Errorf("Expected protocol TUIC, got %s", protocol)
	}

	// Test IsConnected (should be false initially)
	connected := handler.IsConnected()
	if connected {
		t.Error("Expected handler to be disconnected initially")
	}

	// Create a test server
	server := core.Server{
		ID:                utils.GenerateID(),
		Name:              "Test TUIC Server",
		Host:              "example.com",
		Port:              443,
		Protocol:          core.ProtocolTUIC,
		Password:          "test-uuid",
		CongestionControl: "cubic",
		UDPRelayMode:      "native",
		ZeroRTTHandshake:  true,
		Heartbeat:         10,
		Enabled:           true,
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

	if details["protocol"] != core.ProtocolTUIC {
		t.Errorf("Expected protocol TUIC in details, got %v", details["protocol"])
	}

	if details["host"] != "example.com" {
		t.Errorf("Expected host example.com in details, got %v", details["host"])
	}

	if details["congestion_control"] != "cubic" {
		t.Errorf("Expected congestion_control cubic in details, got %v", details["congestion_control"])
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

func TestTUICHandlerMissingUUID(t *testing.T) {
	// Create TUIC handler
	handler := NewTUICHandler()

	// Create a test server without UUID
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test TUIC Server (Missing UUID)",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolTUIC,
		Enabled:  true,
	}

	// Test Connect should fail
	err := handler.Connect(server)
	if err == nil {
		t.Error("Expected error when connecting with missing UUID")
	}

	if err != nil && err.Error() != "missing UUID for TUIC protocol" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}
