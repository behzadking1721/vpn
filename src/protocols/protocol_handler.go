package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"errors"
	"sync"
	"time"
)

// ProtocolHandler defines the interface for handling different VPN protocols
type ProtocolHandler interface {
	// Connect establishes a connection to the server
	Connect(server core.Server) error
	
	// Disconnect terminates the connection
	Disconnect() error
	
	// IsConnected checks if the connection is active
	IsConnected() bool
	
	// GetProtocol returns the protocol type this handler supports
	GetProtocol() core.ProtocolType
	
	// GetDataUsage returns the amount of data sent and received
	GetDataUsage() (sent, received int64, err error)
	
	// GetConnectionDetails returns detailed connection information
	GetConnectionDetails() (map[string]interface{}, error)
	
	// UpdateDataUsage updates the data usage statistics
	UpdateDataUsage(sent, received int64)
}

// ProtocolFactory creates protocol handlers
type ProtocolFactory struct{}

// NewProtocolFactory creates a new protocol factory
func NewProtocolFactory() *ProtocolFactory {
	return &ProtocolFactory{}
}

// CreateHandler creates a protocol handler for the specified protocol type
func (pf *ProtocolFactory) CreateHandler(protocolType core.ProtocolType) (ProtocolHandler, error) {
	switch protocolType {
	case core.ProtocolVMess:
		return NewVMessHandler(), nil
	case core.ProtocolVLESS:
		return NewVLESSHandler(), nil
	case core.ProtocolTrojan:
		return NewTrojanHandler(), nil
	case core.ProtocolReality:
		return NewRealityHandler(), nil
	case core.ProtocolHysteria:
		return NewHysteriaHandler(), nil
	case core.ProtocolTUIC:
		return NewTUICHandler(), nil
	case core.ProtocolSSH:
		return NewSSHHandler(), nil
	case core.ProtocolShadowsocks:
		return NewShadowsocksHandler(), nil
	default:
		return nil, errors.New("unsupported protocol type")
	}
}

// BaseHandler provides a base implementation for protocol handlers
type BaseHandler struct {
	connected   bool
	protocol    core.ProtocolType
	dataSent    int64
	dataReceived int64
	mutex       sync.RWMutex
	lastUpdate  time.Time
}

// IsConnected checks if the connection is active
func (bh *BaseHandler) IsConnected() bool {
	return bh.connected
}

// GetProtocol returns the protocol type
func (bh *BaseHandler) GetProtocol() core.ProtocolType {
	return bh.protocol
}

// GetDataUsage returns the amount of data sent and received
func (bh *BaseHandler) GetDataUsage() (sent, received int64, err error) {
	bh.mutex.RLock()
	defer bh.mutex.RUnlock()
	
	return bh.dataSent, bh.dataReceived, nil
}

// UpdateDataUsage updates the data usage statistics
func (bh *BaseHandler) UpdateDataUsage(sent, received int64) {
	bh.mutex.Lock()
	defer bh.mutex.Unlock()
	
	bh.dataSent += sent
	bh.dataReceived += received
	bh.lastUpdate = time.Now()
}

// GetConnectionDetails returns detailed connection information
func (bh *BaseHandler) GetConnectionDetails() (map[string]interface{}, error) {
	// This would be implemented by each specific protocol handler
	return nil, errors.New("not implemented")
}