package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"math/rand"
	"time"
)

// SSHHandler handles SSH protocol connections
type SSHHandler struct {
	BaseHandler
	server core.Server
	stopCh chan struct{}
}

// NewSSHHandler creates a new SSH handler
func NewSSHHandler() *SSHHandler {
	handler := &SSHHandler{
		stopCh: make(chan struct{}),
	}
	handler.BaseHandler.protocol = core.ProtocolSSH
	return handler
}

// Connect establishes a connection to the SSH server
func (sh *SSHHandler) Connect(server core.Server) error {
	// Store server info
	sh.server = server
	
	// In a real implementation, this would:
	// 1. Parse the SSH configuration
	// 2. Initialize the SSH client library
	// 3. Establish SSH tunnel to the server
	
	fmt.Printf("Connecting to SSH server: %s:%d\n", server.Host, server.Port)
	
	// Simulate connection process
	time.Sleep(1 * time.Second)
	
	// Check for authentication method
	if server.Username == "" {
		return fmt.Errorf("missing username for SSH protocol")
	}
	
	if server.Password == "" && server.PrivateKey == "" {
		return fmt.Errorf("missing password or private key for SSH authentication")
	}
	
	// Mark as connected
	sh.BaseHandler.connected = true
	fmt.Println("SSH connection established")
	
	// Start data usage simulation in a goroutine
	go sh.simulateDataUsage()
	
	return nil
}

// Disconnect terminates the SSH connection
func (sh *SSHHandler) Disconnect() error {
	if !sh.BaseHandler.connected {
		return fmt.Errorf("not connected to SSH server")
	}
	
	// Signal to stop the goroutine
	close(sh.stopCh)
	
	// In a real implementation, this would:
	// 1. Close the SSH client connection
	// 2. Clean up resources
	
	fmt.Printf("Disconnecting from SSH server: %s:%d\n", sh.server.Host, sh.server.Port)
	
	// Simulate disconnection process
	time.Sleep(500 * time.Millisecond)
	sh.BaseHandler.connected = false
	fmt.Println("SSH connection terminated")
	
	return nil
}

// GetConnectionDetails returns detailed connection information
func (sh *SSHHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !sh.BaseHandler.connected {
		return nil, fmt.Errorf("not connected to SSH server")
	}
	
	details := map[string]interface{}{
		"protocol":     sh.BaseHandler.protocol,
		"host":         sh.server.Host,
		"port":         sh.server.Port,
		"username":     sh.server.Username,
		"password_set": sh.server.Password != "",
		"private_key":  sh.server.PrivateKey != "",
		"connected":    sh.BaseHandler.connected,
	}
	
	return details, nil
}

// simulateDataUsage simulates data usage for demonstration purposes
func (sh *SSHHandler) simulateDataUsage() {
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
func (sh *SSHHandler) GetDataUsage() (sent, received int64, err error) {
	if !sh.BaseHandler.connected {
		return 0, 0, fmt.Errorf("not connected to SSH server")
	}
	
	return sh.BaseHandler.GetDataUsage()
}