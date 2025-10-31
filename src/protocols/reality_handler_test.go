package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"testing"
	"time"
)

func TestRealityHandler(t *testing.T) {
	// Create Reality handler
	handler := NewRealityHandler()

	// Test GetProtocol
	protocol := handler.GetProtocol()
	if protocol != core.ProtocolReality {
		t.Errorf("Expected protocol Reality, got %s", protocol)
	}

	// Test IsConnected (should be false initially)
	connected := handler.IsConnected()
	if connected {
		t.Error("Expected handler to be disconnected initially")
	}

	// Create a test server
	server := core.Server{
		ID:          utils.GenerateID(),
		Name:        "Test Reality Server",
		Host:        "example.com",
		Port:        443,
		Protocol:    core.ProtocolReality,
		Password:    "test-uuid",
		SNI:         "example.com",
		Fingerprint: "chrome",
		PublicKey:   "public-key",
		ShortID:     "short-id",
		SpiderX:     "/",
		Enabled:     true,
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

	if details["protocol"] != core.ProtocolReality {
		t.Errorf("Expected protocol Reality in details, got %v", details["protocol"])
	}

	if details["host"] != "example.com" {
		t.Errorf("Expected host example.com in details, got %v", details["host"])
	}

	if details["sni"] != "example.com" {
		t.Errorf("Expected SNI example.com in details, got %v", details["sni"])
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

func TestRealityHandlerMissingParameters(t *testing.T) {
	// Create Reality handler
	handler := NewRealityHandler()

	// Create a test server without required parameters
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test Reality Server (Missing Parameters)",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolReality,
		Password: "test-uuid",
		Enabled:  true,
		// Missing SNI and Fingerprint
	}

	// Test Connect should fail
	err := handler.Connect(server)
	if err == nil {
		t.Error("Expected error when connecting with missing parameters")
	}

	expectedErrors := []string{"missing SNI", "missing fingerprint"}
	foundExpectedError := false
	for _, expectedError := range expectedErrors {
		if err != nil && err.Error() == expectedError+" for Reality protocol" {
			foundExpectedError = true
			break
		}
	}

	if !foundExpectedError {
		t.Errorf("Expected one of the specific errors, got: %v", err)
	}
}
