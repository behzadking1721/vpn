package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"time"
)

// VMessHandler handles VMess protocol connections
type VMessHandler struct {
	BaseHandler
	server core.Server
}

// NewVMessHandler creates a new VMess handler
func NewVMessHandler() *VMessHandler {
	handler := &VMessHandler{}
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
	
	return nil
}

// Disconnect terminates the VMess connection
func (vh *VMessHandler) Disconnect() error {
	if !vh.BaseHandler.connected {
		return fmt.Errorf("not connected to VMess server")
	}
	
	// In a real implementation, this would:
	// 1. Close the V2Ray/XRay core connection
	// 2. Clean up resources
	
	fmt.Printf("Disconnecting from VMess server: %s:%d\n", vh.server.Host, vh.server.Port)
	
	// Simulate disconnection process
	time.Sleep(500 * time.Millisecond)
	vh.BaseHandler.connected = false
	vh.server = core.Server{} // Clear server info
	
	fmt.Println("VMess connection terminated")
	
	return nil
}

// GetDataUsage returns the amount of data sent and received
func (vh *VMessHandler) GetDataUsage() (sent, received int64, err error) {
	if !vh.BaseHandler.connected {
		return 0, 0, fmt.Errorf("not connected to VMess server")
	}
	
	// In a real implementation, this would get actual data from the V2Ray/XRay core
	// For now, we'll simulate data usage
	sent = 1024 * 1024     // 1 MB
	received = 5 * 1024 * 1024 // 5 MB
	
	return sent, received, nil
}

// GetConnectionDetails returns detailed connection information
func (vh *VMessHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !vh.BaseHandler.connected {
		return nil, fmt.Errorf("not connected to VMess server")
	}
	
	details := map[string]interface{}{
		"protocol":   "VMess",
		"host":       vh.server.Host,
		"port":       vh.server.Port,
		"encryption": vh.server.Encryption,
		"tls":        vh.server.TLS,
		"connected":  vh.BaseHandler.connected,
	}
	
	if vh.server.TLS {
		details["sni"] = vh.server.SNI
		details["fingerprint"] = vh.server.Fingerprint
	}
	
	return details, nil
}