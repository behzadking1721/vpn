package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"testing"
	"time"
)

func TestVLESSHandler(t *testing.T) {
	// Create VLESS handler
	handler := NewVLESSHandler()

	// Test GetProtocol
	protocol := handler.GetProtocol()
	if protocol != core.ProtocolVLESS {
		t.Errorf("Expected protocol VLESS, got %s", protocol)
	}

	// Test IsConnected (should be false initially)
	connected := handler.IsConnected()
	if connected {
		t.Error("Expected handler to be disconnected initially")
	}

	// Create a test server
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test VLESS Server",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolVLESS,
		Password: "test-uuid",
		TLS:      true,
		SNI:      "example.com",
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

	if details["protocol"] != core.ProtocolVLESS {
		t.Errorf("Expected protocol VLESS in details, got %v", details["protocol"])
	}

	if details["host"] != "example.com" {
		t.Errorf("Expected host example.com in details, got %v", details["host"])
	}

	if details["tls"] != true {
		t.Errorf("Expected TLS to be true in details, got %v", details["tls"])
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

func TestVLESSHandlerWithoutTLS(t *testing.T) {
	// Create VLESS handler
	handler := NewVLESSHandler()

	// Create a test server without TLS
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test VLESS Server (No TLS)",
		Host:     "example.com",
		Port:     80,
		Protocol: core.ProtocolVLESS,
		Password: "test-uuid",
		TLS:      false,
		Enabled:  true,
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

	if details["tls"] != false {
		t.Errorf("Expected TLS to be false in details, got %v", details["tls"])
	}

	// Test Disconnect
	err = handler.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect: %v", err)
	}
}
