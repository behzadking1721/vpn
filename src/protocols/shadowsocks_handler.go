package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"math/rand"
	"time"
)

// ShadowsocksHandler handles Shadowsocks protocol connections
type ShadowsocksHandler struct {
	BaseHandler
	server core.Server
	stopCh chan struct{}
}

// NewShadowsocksHandler creates a new Shadowsocks handler
func NewShadowsocksHandler() *ShadowsocksHandler {
	handler := &ShadowsocksHandler{
		stopCh: make(chan struct{}),
	}
	handler.BaseHandler.protocol = core.ProtocolShadowsocks
	return handler
}

// Connect establishes a connection to the Shadowsocks server
func (sh *ShadowsocksHandler) Connect(server core.Server) error {
	// Store server info
	sh.server = server
	
	// In a real implementation, this would:
	// 1. Parse the Shadowsocks configuration
	// 2. Initialize the Shadowsocks client library
	// 3. Establish connection to the server
	
	fmt.Printf("Connecting to Shadowsocks server: %s:%d\n", server.Host, server.Port)
	fmt.Printf("Method: %s, Password: %s\n", server.Method, server.Password)
	
	// Simulate connection process
	time.Sleep(1 * time.Second)
	
	// Check for required parameters
	if server.Method == "" {
		return fmt.Errorf("missing encryption method")
	}
	
	if server.Password == "" {
		return fmt.Errorf("missing password")
	}
	
	// Mark as connected
	sh.BaseHandler.connected = true
	fmt.Println("Shadowsocks connection established")
	
	// Start data usage simulation in a goroutine
	go sh.simulateDataUsage()
	
	return nil
}

// Disconnect terminates the Shadowsocks connection
func (sh *ShadowsocksHandler) Disconnect() error {
	if !sh.BaseHandler.connected {
		return fmt.Errorf("not connected to Shadowsocks server")
	}
	
	// Signal to stop the goroutine
	close(sh.stopCh)
	
	// In a real implementation, this would:
	// 1. Close the Shadowsocks client connection
	// 2. Clean up resources
	
	fmt.Printf("Disconnecting from Shadowsocks server: %s:%d\n", sh.server.Host, sh.server.Port)
	
	// Simulate disconnection process
	time.Sleep(500 * time.Millisecond)
	sh.BaseHandler.connected = false
	fmt.Println("Shadowsocks connection terminated")
	
	return nil
}

// GetConnectionDetails returns detailed connection information
func (sh *ShadowsocksHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !sh.BaseHandler.connected {
		return nil, fmt.Errorf("not connected to Shadowsocks server")
	}
	
	details := map[string]interface{}{
		"protocol":  sh.BaseHandler.protocol,
		"host":      sh.server.Host,
		"port":      sh.server.Port,
		"method":    sh.server.Method,
		"connected": sh.BaseHandler.connected,
	}
	
	return details, nil
}

// simulateDataUsage simulates data usage for demonstration purposes
func (sh *ShadowsocksHandler) simulateDataUsage() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-sh.stopCh:
			return
		case <-ticker.C:
			// Simulate data transfer
			sent := rand.Int63n(1024) + 512     // 0.5KB to 1.5KB
			received := rand.Int63n(2048) + 1024 // 1KB to 3KB
			
			// Update data usage
			sh.BaseHandler.UpdateDataUsage(sent, received)
		}
	}
}

// GetDataUsage returns the amount of data sent and received
func (sh *ShadowsocksHandler) GetDataUsage() (sent, received int64, err error) {
	if !sh.BaseHandler.connected {
		return 0, 0, fmt.Errorf("not connected to Shadowsocks server")
	}
	
	return sh.BaseHandler.GetDataUsage()
}