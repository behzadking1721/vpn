package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"math/rand"
	"time"
)

// VLESSHandler handles VLESS protocol connections
type VLESSHandler struct {
	BaseHandler
	server core.Server
	stopCh chan struct{}
}

// NewVLESSHandler creates a new VLESS handler
func NewVLESSHandler() *VLESSHandler {
	handler := &VLESSHandler{
		stopCh: make(chan struct{}),
	}
	handler.BaseHandler.protocol = core.ProtocolVLESS
	return handler
}

// Connect establishes a connection to the VLESS server
func (vh *VLESSHandler) Connect(server core.Server) error {
	// Store server info
	vh.server = server
	
	// In a real implementation, this would:
	// 1. Parse the VLESS configuration
	// 2. Initialize the VLESS client library
	// 3. Establish connection to the server
	
	fmt.Printf("Connecting to VLESS server: %s:%d\n", server.Host, server.Port)
	fmt.Printf("UUID: %s\n", server.Password)
	
	// Simulate connection process
	time.Sleep(1 * time.Second)
	
	// Check for required parameters
	if server.Password == "" {
		return fmt.Errorf("missing UUID")
	}
	
	// Handle TLS if enabled
	if server.TLS {
		fmt.Println("Establishing TLS connection...")
		time.Sleep(500 * time.Millisecond)
	}
	
	// Mark as connected
	vh.BaseHandler.connected = true
	fmt.Println("VLESS connection established")
	
	// Start data usage simulation in a goroutine
	go vh.simulateDataUsage()
	
	return nil
}

// Disconnect terminates the VLESS connection
func (vh *VLESSHandler) Disconnect() error {
	if !vh.BaseHandler.connected {
		return fmt.Errorf("not connected to VLESS server")
	}
	
	// Signal to stop the goroutine
	close(vh.stopCh)
	
	// In a real implementation, this would:
	// 1. Close the VLESS client connection
	// 2. Clean up resources
	
	fmt.Printf("Disconnecting from VLESS server: %s:%d\n", vh.server.Host, vh.server.Port)
	
	// Simulate disconnection process
	time.Sleep(500 * time.Millisecond)
	vh.BaseHandler.connected = false
	fmt.Println("VLESS connection terminated")
	
	return nil
}

// GetConnectionDetails returns detailed connection information
func (vh *VLESSHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !vh.BaseHandler.connected {
		return nil, fmt.Errorf("not connected to VLESS server")
	}
	
	details := map[string]interface{}{
		"protocol":  vh.BaseHandler.protocol,
		"host":      vh.server.Host,
		"port":      vh.server.Port,
		"uuid":      vh.server.Password,
		"tls":       vh.server.TLS,
		"connected": vh.BaseHandler.connected,
	}
	
	if vh.server.TLS {
		details["sni"] = vh.server.SNI
		details["fingerprint"] = vh.server.Fingerprint
	}
	
	return details, nil
}

// simulateDataUsage simulates data usage for demonstration purposes
func (vh *VLESSHandler) simulateDataUsage() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-vh.stopCh:
			return
		case <-ticker.C:
			// Simulate data transfer
			sent := rand.Int63n(1024) + 512     // 0.5KB to 1.5KB
			received := rand.Int63n(2048) + 1024 // 1KB to 3KB
			
			// Update data usage
			vh.BaseHandler.UpdateDataUsage(sent, received)
		}
	}
}

// GetDataUsage returns the amount of data sent and received
func (vh *VLESSHandler) GetDataUsage() (sent, received int64, err error) {
	if !vh.BaseHandler.connected {
		return 0, 0, fmt.Errorf("not connected to VLESS server")
	}
	
	return vh.BaseHandler.GetDataUsage()
}