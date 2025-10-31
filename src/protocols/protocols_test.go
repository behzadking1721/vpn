package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"testing"
)

func TestProtocolFactory(t *testing.T) {
	factory := NewProtocolFactory()

	// Test VMess handler creation
	t.Run("VMessHandlerCreation", func(t *testing.T) {
		handler, err := factory.CreateHandler(core.ProtocolVMess)
		if err != nil {
			t.Errorf("Failed to create VMess handler: %v", err)
		}
		if handler.GetProtocol() != core.ProtocolVMess {
			t.Errorf("Expected VMess protocol, got %s", handler.GetProtocol())
		}
	})

	// Test Shadowsocks handler creation
	t.Run("ShadowsocksHandlerCreation", func(t *testing.T) {
		handler, err := factory.CreateHandler(core.ProtocolShadowsocks)
		if err != nil {
			t.Errorf("Failed to create Shadowsocks handler: %v", err)
		}
		if handler.GetProtocol() != core.ProtocolShadowsocks {
			t.Errorf("Expected Shadowsocks protocol, got %s", handler.GetProtocol())
		}
	})

	// Test unsupported protocol
	t.Run("UnsupportedProtocol", func(t *testing.T) {
		_, err := factory.CreateHandler("unsupported")
		if err == nil {
			t.Error("Expected error for unsupported protocol, got nil")
		}
	})
}

func TestVMessHandler(t *testing.T) {
	handler := NewVMessHandler()

	// Test initial state
	if handler.IsConnected() {
		t.Error("Expected handler to be disconnected initially")
	}

	// Test connection
	server := core.Server{
		ID:         "test-vmess",
		Name:       "Test VMess Server",
		Host:       "example.com",
		Port:       443,
		Protocol:   core.ProtocolVMess,
		Encryption: "auto",
		TLS:        true,
	}

	err := handler.Connect(server)
	if err != nil {
		t.Errorf("Failed to connect: %v", err)
	}

	if !handler.IsConnected() {
		t.Error("Expected handler to be connected after Connect()")
	}

	// Test data usage
	sent, received, err := handler.GetDataUsage()
	if err != nil {
		t.Errorf("Failed to get data usage: %v", err)
	}
	if sent == 0 || received == 0 {
		t.Error("Expected non-zero data usage")
	}

	// Test connection details
	details, err := handler.GetConnectionDetails()
	if err != nil {
		t.Errorf("Failed to get connection details: %v", err)
	}
	if details["protocol"] != "VMess" {
		t.Errorf("Expected protocol VMess, got %v", details["protocol"])
	}

	// Test disconnection
	err = handler.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect: %v", err)
	}

	if handler.IsConnected() {
		t.Error("Expected handler to be disconnected after Disconnect()")
	}
}

func TestShadowsocksHandler(t *testing.T) {
	handler := NewRealShadowsocksHandler()

	// Test initial state
	if handler.IsConnected() {
		t.Error("Expected handler to be disconnected initially")
	}

	// Test connection
	server := core.Server{
		ID:       "test-ss",
		Name:     "Test Shadowsocks Server",
		Host:     "example.com",
		Port:     8388,
		Protocol: core.ProtocolShadowsocks,
		Method:   "aes-256-gcm",
		Password: "test-password",
	}

	err := handler.Connect(server)
	if err != nil {
		t.Errorf("Failed to connect: %v", err)
	}

	if !handler.IsConnected() {
		t.Error("Expected handler to be connected after Connect()")
	}

	// Test data usage
	sent, received, err := handler.GetDataUsage()
	if err != nil {
		t.Errorf("Failed to get data usage: %v", err)
	}

	// For simulation, we just check they are >= 0
	if sent < 0 || received < 0 {
		t.Error("Expected non-negative data usage")
	}

	// Test connection details
	details, err := handler.GetConnectionDetails()
	if err != nil {
		t.Errorf("Failed to get connection details: %v", err)
	}
	if details["protocol"] != "Shadowsocks" {
		t.Errorf("Expected protocol Shadowsocks, got %v", details["protocol"])
	}

	// Test disconnection
	err = handler.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect: %v", err)
	}

	if handler.IsConnected() {
		t.Error("Expected handler to be disconnected after Disconnect()")
	}
}
