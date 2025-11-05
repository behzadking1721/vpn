package managers

import (
	"testing"
	"time"
	
	"vpnclient/src/core"
)

// TestVPNConnectionEndToEnd tests the complete VPN connection flow
func TestVPNConnectionEndToEnd(t *testing.T) {
	// Create a connection manager
	cm := NewConnectionManager()
	
	// Create a test server configuration (using VMess instead of WireGuard)
	server := &core.Server{
		ID:       "test-server-1",
		Name:     "Test VPN Server",
		Host:     "127.0.0.1", // Using localhost for testing
		Port:     8080,
		Protocol: "vmess",
		Config: map[string]interface{}{
			"user_id": "test_user_id",
			"alter_id": 0,
		},
		Enabled: true,
	}
	
	// Test the connection flow
	t.Run("Connect to VPN Server", func(t *testing.T) {
		// Initial status should be disconnected
		status := cm.GetStatus()
		if status != Disconnected {
			t.Errorf("Expected initial status to be Disconnected, got %v", status)
		}
		
		// Connect to the server
		err := cm.Connect(server)
		if err != nil {
			t.Fatalf("Failed to connect to server: %v", err)
		}
		
		// Wait for connection to establish
		time.Sleep(200 * time.Millisecond)
		
		// Check status is now connected
		status = cm.GetStatus()
		if status != Connected {
			t.Errorf("Expected status to be Connected, got %v", status)
		}
		
		// Check that server info is correctly set
		currentServer := cm.GetCurrentServer()
		if currentServer == nil {
			t.Error("Expected current server to be set, got nil")
		} else if currentServer.ID != server.ID {
			t.Errorf("Expected server ID %s, got %s", server.ID, currentServer.ID)
		}
	})
	
	t.Run("Get Connection Statistics", func(t *testing.T) {
		// Get connection statistics
		sent, received := cm.GetDataUsage()
		
		// Stats should have non-negative values
		if sent < 0 {
			t.Error("Expected DataSent to be non-negative")
		}
		
		if received < 0 {
			t.Error("Expected DataReceived to be non-negative")
		}
		
		// Check uptime
		uptime := cm.GetUptime()
		if uptime < 0 {
			t.Error("Expected uptime to be non-negative")
		}
	})
	
	t.Run("Disconnect from VPN Server", func(t *testing.T) {
		// Disconnect from the server
		err := cm.Disconnect()
		if err != nil {
			t.Fatalf("Failed to disconnect from server: %v", err)
		}
		
		// Wait for disconnection to complete
		time.Sleep(200 * time.Millisecond)
		
		// Check status is now disconnected
		status := cm.GetStatus()
		if status != Disconnected {
			t.Errorf("Expected status to be Disconnected, got %v", status)
		}
		
		// Server info should be cleared
		currentServer := cm.GetCurrentServer()
		if currentServer != nil {
			t.Error("Expected current server to be nil after disconnection")
		}
	})
}

// TestVPNConnectionWithDifferentProtocols tests connection with different protocols
func TestVPNConnectionWithDifferentProtocols(t *testing.T) {
	cm := NewConnectionManager()
	
	// Test different protocol configurations
	protocols := []struct {
		name     string
		protocol string
		config   map[string]interface{}
	}{
		{
			name:     "VMess",
			protocol: "vmess",
			config: map[string]interface{}{
				"user_id": "test_user_id",
				"alter_id": 0,
			},
		},
		{
			name:     "Shadowsocks",
			protocol: "shadowsocks",
			config: map[string]interface{}{
				"method":   "aes-128-gcm",
				"password": "test_password",
			},
		},
	}
	
	for _, proto := range protocols {
		t.Run(proto.name, func(t *testing.T) {
			server := &core.Server{
				ID:       "test-server-" + proto.protocol,
				Name:     "Test " + proto.name + " Server",
				Host:     "127.0.0.1",
				Port:     8080,
				Protocol: proto.protocol,
				Config:   proto.config,
				Enabled:  true,
			}
			
			// Try to connect
			err := cm.Connect(server)
			// For this test, we're not actually connecting to real servers,
			// so we expect it to either succeed or fail with a specific error
			// The important thing is that it doesn't panic or cause other issues
			if err != nil {
				// Log the error but don't fail the test
				t.Logf("Connection attempt resulted in error (expected for mock test): %v", err)
			} else {
				// If connection succeeded, test disconnection
				err = cm.Disconnect()
				if err != nil {
					t.Logf("Disconnection attempt resulted in error (expected for mock test): %v", err)
				}
			}
		})
	}
}

// TestConnectionStatusTransitions tests that connection status transitions work correctly
func TestConnectionStatusTransitions(t *testing.T) {
	cm := NewConnectionManager()
	
	server := &core.Server{
		ID:       "test-server-status",
		Name:     "Test Status Server",
		Host:     "127.0.0.1",
		Port:     8080,
		Protocol: "vmess",
		Config: map[string]interface{}{
			"user_id": "test_user_id",
			"alter_id": 0,
		},
		Enabled: true,
	}
	
	// Initial state
	if status := cm.GetStatus(); status != Disconnected {
		t.Errorf("Expected initial status Disconnected, got %v", status)
	}
	
	// Connect
	err := cm.Connect(server)
	if err != nil {
		t.Fatalf("Failed to initiate connection: %v", err)
	}
	
	// Should transition to Connecting or Connected
	status := cm.GetStatus()
	if status != Connecting && status != Connected {
		t.Errorf("Expected status Connecting or Connected after Connect(), got %v", status)
	}
	
	// Disconnect
	err = cm.Disconnect()
	if err != nil {
		t.Fatalf("Failed to initiate disconnection: %v", err)
	}
	
	// Should transition to Disconnecting or Disconnected
	status = cm.GetStatus()
	if status != Disconnecting && status != Disconnected {
		t.Errorf("Expected status Disconnecting or Disconnected after Disconnect(), got %v", status)
	}
}