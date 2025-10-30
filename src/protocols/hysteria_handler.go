package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"math/rand"
	"time"
)

// HysteriaHandler handles Hysteria protocol connections
type HysteriaHandler struct {
	BaseHandler
	server core.Server
	stopCh chan struct{}
}

// NewHysteriaHandler creates a new Hysteria handler
func NewHysteriaHandler() *HysteriaHandler {
	handler := &HysteriaHandler{
		stopCh: make(chan struct{}),
	}
	handler.BaseHandler.protocol = core.ProtocolHysteria
	return handler
}

// Connect establishes a connection to the Hysteria server
func (hh *HysteriaHandler) Connect(server core.Server) error {
	// Store server info
	hh.server = server
	
	// In a real implementation, this would:
	// 1. Parse the Hysteria configuration
	// 2. Initialize the Hysteria client library
	// 3. Establish connection to the server using QUIC
	
	fmt.Printf("Connecting to Hysteria server: %s:%d\n", server.Host, server.Port)
	fmt.Printf("Password: %s\n", server.Password)
	
	// Simulate connection process
	time.Sleep(1 * time.Second)
	
	// Check for required parameters
	if server.Password == "" {
		return fmt.Errorf("missing password for Hysteria protocol")
	}
	
	// Handle Hysteria specific parameters
	if server.UpMbps <= 0 {
		fmt.Println("Warning: Upload speed not specified, using default")
	}
	
	if server.DownMbps <= 0 {
		fmt.Println("Warning: Download speed not specified, using default")
	}
	
	// Mark as connected
	hh.BaseHandler.connected = true
	fmt.Println("Hysteria connection established")
	
	// Start data usage simulation in a goroutine
	go hh.simulateDataUsage()
	
	return nil
}

// Disconnect terminates the Hysteria connection
func (hh *HysteriaHandler) Disconnect() error {
	if !hh.BaseHandler.connected {
		return fmt.Errorf("not connected to Hysteria server")
	}
	
	// Signal to stop the goroutine
	close(hh.stopCh)
	
	// In a real implementation, this would:
	// 1. Close the Hysteria client connection
	// 2. Clean up resources
	
	fmt.Printf("Disconnecting from Hysteria server: %s:%d\n", hh.server.Host, hh.server.Port)
	
	// Simulate disconnection process
	time.Sleep(500 * time.Millisecond)
	hh.BaseHandler.connected = false
	fmt.Println("Hysteria connection terminated")
	
	return nil
}

// GetConnectionDetails returns detailed connection information
func (hh *HysteriaHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !hh.BaseHandler.connected {
		return nil, fmt.Errorf("not connected to Hysteria server")
	}
	
	details := map[string]interface{}{
		"protocol":    hh.BaseHandler.protocol,
		"host":        hh.server.Host,
		"port":        hh.server.Port,
		"password":    hh.server.Password,
		"up_mbps":     hh.server.UpMbps,
		"down_mbps":   hh.server.DownMbps,
		"obfs":        hh.server.Obfs,
		"auth_str":    hh.server.AuthStr,
		"insecure":    hh.server.Insecure,
		"connected":   hh.BaseHandler.connected,
	}
	
	return details, nil
}

// simulateDataUsage simulates data usage for demonstration purposes
func (hh *HysteriaHandler) simulateDataUsage() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-hh.stopCh:
			return
		case <-ticker.C:
			// Simulate data transfer
			sent := rand.Int63n(1024) + 512     // 0.5KB to 1.5KB
			received := rand.Int63n(2048) + 1024 // 1KB to 3KB
			
			// Update data usage
			hh.BaseHandler.UpdateDataUsage(sent, received)
		}
	}
}

// GetDataUsage returns the amount of data sent and received
func (hh *HysteriaHandler) GetDataUsage() (sent, received int64, err error) {
	if !hh.BaseHandler.connected {
		return 0, 0, fmt.Errorf("not connected to Hysteria server")
	}
	
	return hh.BaseHandler.GetDataUsage()
}