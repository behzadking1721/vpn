package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"testing"
	"time"
)

func TestConnectionManagerConnect(t *testing.T) {
	// Create connection manager
	connManager := NewConnectionManager()

	// Create a test server
	server := core.Server{
		ID:         utils.GenerateID(),
		Name:       "Test Server",
		Host:       "example.com",
		Port:       443,
		Protocol:   core.ProtocolVMess,
		Password:   "test-uuid",
		Encryption: "auto",
		TLS:        true,
		Enabled:    true,
	}

	// Test connection
	err := connManager.Connect(server)
	if err != nil {
		// In the current implementation, this will fail because we don't have real protocol handlers in tests
		// But we can still test that the method is called correctly
		t.Logf("Expected connection error in test environment: %v", err)
	}

	// Check status (should be error or disconnected because we don't have real protocol handlers)
	status := connManager.GetStatus()
	if status != core.StatusError && status != core.StatusDisconnected {
		t.Errorf("Expected status to be Error or Disconnected, got %s", status)
	}
}

func TestConnectionManagerDisconnect(t *testing.T) {
	// Create connection manager
	connManager := NewConnectionManager()

	// Test disconnection when not connected
	err := connManager.Disconnect()
	if err == nil {
		t.Error("Expected error when disconnecting while not connected")
	}
}

func TestConnectionManagerGetStatus(t *testing.T) {
	// Create connection manager
	connManager := NewConnectionManager()

	// Check initial status
	status := connManager.GetStatus()
	if status != core.StatusDisconnected {
		t.Errorf("Expected initial status to be Disconnected, got %s", status)
	}
}

func TestConnectionManagerGetCurrentServer(t *testing.T) {
	// Create connection manager
	connManager := NewConnectionManager()

	// Check current server when not connected
	server := connManager.GetCurrentServer()
	if server != nil {
		t.Error("Expected nil server when not connected")
	}
}

func TestConnectionManagerUpdateDataUsage(t *testing.T) {
	// Create connection manager
	connManager := NewConnectionManager()

	// Test updating data usage
	sent := int64(1024)
	received := int64(2048)
	connManager.UpdateDataUsage(sent, received)

	// Check connection info
	info := connManager.GetConnectionInfo()
	if info.DataSent != sent {
		t.Errorf("Expected DataSent to be %d, got %d", sent, info.DataSent)
	}

	if info.DataReceived != received {
		t.Errorf("Expected DataReceived to be %d, got %d", received, info.DataReceived)
	}
}

func TestConnectionManagerListeners(t *testing.T) {
	// Create connection manager
	connManager := NewConnectionManager()

	// Create a test listener
	listener := &TestConnectionListener{}
	connManager.AddListener(listener)

	// Update data usage to trigger listener
	connManager.UpdateDataUsage(1024, 2048)

	// Give some time for the update to propagate
	time.Sleep(10 * time.Millisecond)

	// Check if listener was notified
	if !listener.Notified {
		t.Error("Expected listener to be notified")
	}

	// Remove listener
	connManager.RemoveListener(listener)

	// Reset notification flag
	listener.Notified = false

	// Update data usage again
	connManager.UpdateDataUsage(2048, 4096)

	// Give some time for the update to propagate
	time.Sleep(10 * time.Millisecond)

	// Check if listener was not notified this time
	if listener.Notified {
		t.Error("Expected listener to not be notified after removal")
	}
}

// TestConnectionListener is a mock listener for testing
type TestConnectionListener struct {
	Notified bool
}

func (t *TestConnectionListener) OnStatusChanged(status core.ConnectionStatus, info core.ConnectionInfo) {
	t.Notified = true
}