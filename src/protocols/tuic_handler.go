package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"math/rand"
	"time"
)

// TUICHandler handles TUIC protocol connections
type TUICHandler struct {
	BaseHandler
	server core.Server
	stopCh chan struct{}
}

// NewTUICHandler creates a new TUIC handler
func NewTUICHandler() *TUICHandler {
	handler := &TUICHandler{
		stopCh: make(chan struct{}),
	}
	handler.BaseHandler.protocol = core.ProtocolTUIC
	return handler
}

// Connect establishes a connection to the TUIC server
func (th *TUICHandler) Connect(server core.Server) error {
	// Store server info
	th.server = server
	
	// In a real implementation, this would:
	// 1. Parse the TUIC configuration
	// 2. Initialize the TUIC client library
	// 3. Establish connection to the server using QUIC
	
	fmt.Printf("Connecting to TUIC server: %s:%d\n", server.Host, server.Port)
	fmt.Printf("UUID: %s\n", server.Password)
	
	// Simulate connection process
	time.Sleep(1 * time.Second)
	
	// Check for required parameters
	if server.Password == "" {
		return fmt.Errorf("missing UUID for TUIC protocol")
	}
	
	// Handle TUIC specific parameters
	if server.CongestionControl == "" {
		fmt.Println("Warning: Congestion control not specified, using default")
	}
	
	// Mark as connected
	th.BaseHandler.connected = true
	fmt.Println("TUIC connection established")
	
	// Start data usage simulation in a goroutine
	go th.simulateDataUsage()
	
	return nil
}

// Disconnect terminates the TUIC connection
func (th *TUICHandler) Disconnect() error {
	if !th.BaseHandler.connected {
		return fmt.Errorf("not connected to TUIC server")
	}
	
	// Signal to stop the goroutine
	close(th.stopCh)
	
	// In a real implementation, this would:
	// 1. Close the TUIC client connection
	// 2. Clean up resources
	
	fmt.Printf("Disconnecting from TUIC server: %s:%d\n", th.server.Host, th.server.Port)
	
	// Simulate disconnection process
	time.Sleep(500 * time.Millisecond)
	th.BaseHandler.connected = false
	fmt.Println("TUIC connection terminated")
	
	return nil
}

// GetConnectionDetails returns detailed connection information
func (th *TUICHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !th.BaseHandler.connected {
		return nil, fmt.Errorf("not connected to TUIC server")
	}
	
	details := map[string]interface{}{
		"protocol":           th.BaseHandler.protocol,
		"host":               th.server.Host,
		"port":               th.server.Port,
		"uuid":               th.server.Password,
		"congestion_control": th.server.CongestionControl,
		"udp_relay_mode":     th.server.UDPRelayMode,
		"zero_rtt_handshake": th.server.ZeroRTTHandshake,
		"heartbeat":          th.server.Heartbeat,
		"connected":          th.BaseHandler.connected,
	}
	
	return details, nil
}

// simulateDataUsage simulates data usage for demonstration purposes
func (th *TUICHandler) simulateDataUsage() {
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
func (th *TUICHandler) GetDataUsage() (sent, received int64, err error) {
	if !th.BaseHandler.connected {
		return 0, 0, fmt.Errorf("not connected to TUIC server")
	}
	
	return th.BaseHandler.GetDataUsage()
}