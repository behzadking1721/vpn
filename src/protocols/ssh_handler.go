package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
)

// SSHHandler handles SSH protocol connections
type SSHHandler struct {
	BaseHandler
}

// NewSSHHandler creates a new SSH handler
func NewSSHHandler() *SSHHandler {
	handler := &SSHHandler{}
	handler.BaseHandler.protocol = core.ProtocolSSH
	return handler
}

// Connect establishes a connection to the SSH server
func (sh *SSHHandler) Connect(server core.Server) error {
	// Implementation would go here
	// This is where you'd handle the actual SSH connection logic
	sh.BaseHandler.connected = true
	return nil
}

// Disconnect terminates the SSH connection
func (sh *SSHHandler) Disconnect() error {
	// Implementation would go here
	// This is where you'd handle disconnecting from the SSH server
	sh.BaseHandler.connected = false
	return nil
}