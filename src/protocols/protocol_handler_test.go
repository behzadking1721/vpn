package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"testing"
	"time"
)

func TestBaseHandlerDataUsage(t *testing.T) {
	handler := &BaseHandler{
		protocol: core.ProtocolVMess,
	}

	// Test initial data usage
	sent, received, err := handler.GetDataUsage()
	if err != nil {
		t.Errorf("Failed to get data usage: %v", err)
	}
	
	if sent != 0 || received != 0 {
		t.Errorf("Expected initial data usage to be 0, got sent=%d, received=%d", sent, received)
	}

	// Update data usage
	handler.UpdateDataUsage(1024, 2048)
	
	// Test updated data usage
	sent, received, err = handler.GetDataUsage()
	if err != nil {
		t.Errorf("Failed to get data usage: %v", err)
	}
	
	if sent != 1024 || received != 2048 {
		t.Errorf("Expected data usage to be sent=1024, received=2048, got sent=%d, received=%d", sent, received)
	}

	// Update data usage again
	handler.UpdateDataUsage(512, 1024)
	
	// Test cumulative data usage
	sent, received, err = handler.GetDataUsage()
	if err != nil {
		t.Errorf("Failed to get data usage: %v", err)
	}
	
	if sent != 1536 || received != 3072 {
		t.Errorf("Expected cumulative data usage to be sent=1536, received=3072, got sent=%d, received=%d", sent, received)
	}
}

func TestVMessHandlerDataUsage(t *testing.T) {
	handler := NewVMessHandler()
	
	// Test data usage before connection
	sent, received, err := handler.GetDataUsage()
	if err != nil {
		t.Errorf("Failed to get data usage: %v", err)
	}
	
	if sent != 0 || received != 0 {
		t.Errorf("Expected initial data usage to be 0, got sent=%d, received=%d", sent, received)
	}
	
	// Connect to a server
	server := core.Server{
		ID:       "test-server",
		Name:     "Test Server",
		Host:     "example.com",
		Port:     443,
		Protocol: core.ProtocolVMess,
		Encryption: "auto",
		TLS:      true,
		Enabled:  true,
	}
	
	err = handler.Connect(server)
	if err != nil {
		t.Errorf("Failed to connect: %v", err)
	}
	
	// Give some time for data usage simulation to start
	time.Sleep(2 * time.Second)
	
	// Test data usage after connection
	sent, received, err = handler.GetDataUsage()
	if err != nil {
		t.Errorf("Failed to get data usage: %v", err)
	}
	
	// Should have some data usage now
	if sent == 0 || received == 0 {
		t.Errorf("Expected some data usage after connection, got sent=%d, received=%d", sent, received)
	}
	
	// Disconnect
	err = handler.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect: %v", err)
	}
}