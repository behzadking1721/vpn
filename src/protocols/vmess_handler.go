package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"math/rand"
	"time"
)

// VMessHandler handles VMess protocol connections
type VMessHandler struct {
	BaseHandler
	server core.Server
	stopCh chan struct{}
}

// NewVMessHandler creates a new VMess handler
func NewVMessHandler() *VMessHandler {
	handler := &VMessHandler{
		stopCh: make(chan struct{}),
	}
	handler.BaseHandler.protocol = core.ProtocolVMess
	return handler
}

// Connect establishes a connection to the VMess server
func (vh *VMessHandler) Connect(server core.Server) error {
	// Store server info
	vh.server = server
	
	// In a real implementation, this would:
	// 1. Parse the VMess URL format
	// 2. Configure the V2Ray/XRay core
	// 3. Establish the connection
	
	fmt.Printf("Connecting to VMess server: %s:%d\n", server.Host, server.Port)
	
	// Simulate connection process
	time.Sleep(1 * time.Second)
	
	// Check for TLS requirement
	if server.TLS {
		fmt.Println("Establishing TLS connection...")
		time.Sleep(500 * time.Millisecond)
	}
	
	// Mark as connected
	vh.BaseHandler.connected = true
	fmt.Println("VMess connection established")
	
	// Start data usage simulation in a goroutine
	go vh.simulateDataUsage()
	
	return nil
}

// Disconnect terminates the VMess connection
func (vh *VMessHandler) Disconnect() error {
	if !vh.BaseHandler.connected {
		return fmt.Errorf("not connected")
	}
	
	// Signal to stop the goroutine
	close(vh.stopCh)
	
	// In a real implementation, this would terminate the actual connection
	fmt.Println("Disconnecting from VMess server...")
	time.Sleep(500 * time.Millisecond)
	
	vh.BaseHandler.connected = false
	fmt.Println("VMess connection terminated")
	
	return nil
}

// GetConnectionDetails returns detailed connection information
func (vh *VMessHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !vh.BaseHandler.connected {
		return nil, errors.New("not connected")
	}
	
	details := map[string]interface{}{
		"protocol":   vh.BaseHandler.protocol,
		"host":       vh.server.Host,
		"port":       vh.server.Port,
		"encryption": vh.server.Encryption,
		"tls":        vh.server.TLS,
	}
	
	if vh.server.SNI != "" {
		details["sni"] = vh.server.SNI
	}
	
	return details, nil
}

// simulateDataUsage simulates data usage for demonstration purposes
func (vh *VMessHandler) simulateDataUsage() {
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
func (vh *VMessHandler) GetDataUsage() (sent, received int64, err error) {
	if !vh.BaseHandler.connected {
		return 0, 0, fmt.Errorf("not connected")
	}
	
	return vh.BaseHandler.GetDataUsage()
}