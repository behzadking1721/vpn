package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"testing"
	"time"
)

func TestVMessHandler(t *testing.T) {
	// Create VMess handler
	handler := NewVMessHandler()

	// Test GetProtocol
	protocol := handler.GetProtocol()
	if protocol != core.ProtocolVMess {
		t.Errorf("Expected protocol VMess, got %s", protocol)
	}

	// Test IsConnected (should be false initially)
	connected := handler.IsConnected()
	if connected {
		t.Error("Expected handler to be disconnected initially")
	}

	// Create a test server
	server := core.Server{
		ID:         utils.GenerateID(),
		Name:       "Test VMess Server",
		Host:       "example.com",
		Port:       443,
		Protocol:   core.ProtocolVMess,
		Password:   "test-uuid",
		Encryption: "auto",
		TLS:        true,
		Enabled:    true,
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

	if details["protocol"] != core.ProtocolVMess {
		t.Errorf("Expected protocol VMess in details, got %v", details["protocol"])
	}

	if details["host"] != "example.com" {
		t.Errorf("Expected host example.com in details, got %v", details["host"])
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

func TestShadowsocksHandler(t *testing.T) {
	// Create Shadowsocks handler
	handler := NewShadowsocksHandler()

	// Test GetProtocol
	protocol := handler.GetProtocol()
	if protocol != core.ProtocolShadowsocks {
		t.Errorf("Expected protocol Shadowsocks, got %s", protocol)
	}

	// Test IsConnected (should be false initially)
	connected := handler.IsConnected()
	if connected {
		t.Error("Expected handler to be disconnected initially")
	}

	// Create a test server
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     "Test Shadowsocks Server",
		Host:     "example.com",
		Port:     8388,
		Protocol: core.ProtocolShadowsocks,
		Method:   "aes-256-gcm",
		Password: "test-password",
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

	if details["protocol"] != core.ProtocolShadowsocks {
		t.Errorf("Expected protocol Shadowsocks in details, got %v", details["protocol"])
	}

	if details["host"] != "example.com" {
		t.Errorf("Expected host example.com in details, got %v", details["host"])
	}

	if details["method"] != "aes-256-gcm" {
		t.Errorf("Expected method aes-256-gcm in details, got %v", details["method"])
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

func TestProtocolFactory(t *testing.T) {
	// Create protocol factory
	factory := NewProtocolFactory()

	// Test creating VMess handler
	handler, err := factory.CreateHandler(core.ProtocolVMess)
	if err != nil {
		t.Errorf("Failed to create VMess handler: %v", err)
	}

	if handler.GetProtocol() != core.ProtocolVMess {
		t.Error("Expected VMess handler")
	}

	// Test creating Shadowsocks handler
	handler, err = factory.CreateHandler(core.ProtocolShadowsocks)
	if err != nil {
		t.Errorf("Failed to create Shadowsocks handler: %v", err)
	}

	if handler.GetProtocol() != core.ProtocolShadowsocks {
		t.Error("Expected Shadowsocks handler")
	}

	// Test creating unsupported protocol
	_, err = factory.CreateHandler("unsupported")
	if err == nil {
		t.Error("Expected error when creating unsupported protocol handler")
	}
}
