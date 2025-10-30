package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"math/rand"
	"time"
)

// TrojanHandler handles Trojan protocol connections
type TrojanHandler struct {
	BaseHandler
	server core.Server
	stopCh chan struct{}
}

// NewTrojanHandler creates a new Trojan handler
func NewTrojanHandler() *TrojanHandler {
	handler := &TrojanHandler{
		stopCh: make(chan struct{}),
	}
	handler.BaseHandler.protocol = core.ProtocolTrojan
	return handler
}

// Connect establishes a connection to the Trojan server
func (th *TrojanHandler) Connect(server core.Server) error {
	// Store server info
	th.server = server
	
	// In a real implementation, this would:
	// 1. Parse the Trojan configuration
	// 2. Initialize the Trojan client library
	// 3. Establish connection to the server
	
	fmt.Printf("Connecting to Trojan server: %s:%d\n", server.Host, server.Port)
	fmt.Printf("Password: %s\n", server.Password)
	
	// Simulate connection process
	time.Sleep(1 * time.Second)
	
	// Check for required parameters
	if server.Password == "" {
		return fmt.Errorf("missing password")
	}
	
	// Mark as connected
	th.BaseHandler.connected = true
	fmt.Println("Trojan connection established")
	
	// Start data usage simulation in a goroutine
	go th.simulateDataUsage()
	
	return nil
}

// Disconnect terminates the Trojan connection
func (th *TrojanHandler) Disconnect() error {
	if !th.BaseHandler.connected {
		return fmt.Errorf("not connected to Trojan server")
	}
	
	// Signal to stop the goroutine
	close(th.stopCh)
	
	// In a real implementation, this would:
	// 1. Close the Trojan client connection
	// 2. Clean up resources
	
	fmt.Printf("Disconnecting from Trojan server: %s:%d\n", th.server.Host, th.server.Port)
	
	// Simulate disconnection process
	time.Sleep(500 * time.Millisecond)
	th.BaseHandler.connected = false
	fmt.Println("Trojan connection terminated")
	
	return nil
}

// GetConnectionDetails returns detailed connection information
func (th *TrojanHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !th.BaseHandler.connected {
		return nil, fmt.Errorf("not connected to Trojan server")
	}
	
	details := map[string]interface{}{
		"protocol":  th.BaseHandler.protocol,
		"host":      th.server.Host,
		"port":      th.server.Port,
		"password":  th.server.Password,
		"connected": th.BaseHandler.connected,
	}
	
	return details, nil
}

// simulateDataUsage simulates data usage for demonstration purposes
func (th *TrojanHandler) simulateDataUsage() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-th.stopCh:
			return
		case <-ticker.C:
			// Simulate data transfer
			sent := rand.Int63n(1024) + 512     // 0.5KB to 1.5KB
			received := rand.Int63n(2048) + 1024 // 1KB to 3KB
			
			// Update data usage
			th.BaseHandler.UpdateDataUsage(sent, received)
		}
	}
}

// GetDataUsage returns the amount of data sent and received
func (th *TrojanHandler) GetDataUsage() (sent, received int64, err error) {
	if !th.BaseHandler.connected {
		return 0, 0, fmt.Errorf("not connected to Trojan server")
	}
	
	return th.BaseHandler.GetDataUsage()
}