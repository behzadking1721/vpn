package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"math/rand"
	"time"
)

// RealityHandler handles Reality protocol connections
type RealityHandler struct {
	BaseHandler
	server core.Server
	stopCh chan struct{}
}

// NewRealityHandler creates a new Reality handler
func NewRealityHandler() *RealityHandler {
	handler := &RealityHandler{
		stopCh: make(chan struct{}),
	}
	handler.BaseHandler.protocol = core.ProtocolReality
	return handler
}

// Connect establishes a connection to the Reality server
func (rh *RealityHandler) Connect(server core.Server) error {
	// Store server info
	rh.server = server

	// In a real implementation, this would:
	// 1. Parse the Reality configuration
	// 2. Initialize the Reality client library
	// 3. Establish connection to the server

	fmt.Printf("Connecting to Reality server: %s:%d\n", server.Host, server.Port)
	fmt.Printf("UUID: %s\n", server.Password)

	// Simulate connection process
	time.Sleep(1 * time.Second)

	// Check for required parameters
	if server.Password == "" {
		return fmt.Errorf("missing UUID")
	}

	// Handle Reality specific parameters
	if server.SNI == "" {
		return fmt.Errorf("missing SNI for Reality protocol")
	}

	if server.Fingerprint == "" {
		return fmt.Errorf("missing fingerprint for Reality protocol")
	}

	// Mark as connected
	rh.BaseHandler.connected = true
	fmt.Println("Reality connection established")

	// Start data usage simulation in a goroutine
	go rh.simulateDataUsage()

	return nil
}

// Disconnect terminates the Reality connection
func (rh *RealityHandler) Disconnect() error {
	if !rh.BaseHandler.connected {
		return fmt.Errorf("not connected to Reality server")
	}

	// Signal to stop the goroutine
	close(rh.stopCh)

	// In a real implementation, this would:
	// 1. Close the Reality client connection
	// 2. Clean up resources

	fmt.Printf("Disconnecting from Reality server: %s:%d\n", rh.server.Host, rh.server.Port)

	// Simulate disconnection process
	time.Sleep(500 * time.Millisecond)
	rh.BaseHandler.connected = false
	fmt.Println("Reality connection terminated")

	return nil
}

// GetConnectionDetails returns detailed connection information
func (rh *RealityHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !rh.BaseHandler.connected {
		return nil, fmt.Errorf("not connected to Reality server")
	}

	details := map[string]interface{}{
		"protocol":    rh.BaseHandler.protocol,
		"host":        rh.server.Host,
		"port":        rh.server.Port,
		"uuid":        rh.server.Password,
		"sni":         rh.server.SNI,
		"fingerprint": rh.server.Fingerprint,
		"public_key":  rh.server.PublicKey,
		"short_id":    rh.server.ShortID,
		"spider_x":    rh.server.SpiderX,
		"connected":   rh.BaseHandler.connected,
	}

	return details, nil
}

// simulateDataUsage simulates data usage for demonstration purposes
func (rh *RealityHandler) simulateDataUsage() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-rh.stopCh:
			return
		case <-ticker.C:
			// Simulate data transfer
			sent := rand.Int63n(1024) + 512      // 0.5KB to 1.5KB
			received := rand.Int63n(2048) + 1024 // 1KB to 3KB

			// Update data usage
			rh.BaseHandler.UpdateDataUsage(sent, received)
		}
	}
}

// GetDataUsage returns the amount of data sent and received
func (rh *RealityHandler) GetDataUsage() (sent, received int64, err error) {
	if !rh.BaseHandler.connected {
		return 0, 0, fmt.Errorf("not connected to Reality server")
	}

	return rh.BaseHandler.GetDataUsage()
}
