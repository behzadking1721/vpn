package core

import (
	"time"
)

// ConnectionStatus represents the connection status
type ConnectionStatus int

const (
	// Disconnected represents a disconnected state
	Disconnected ConnectionStatus = iota
	// Connecting represents a connecting state
	Connecting
	// Connected represents a connected state
	Connected
	// Disconnecting represents a disconnecting state
	Disconnecting
	// Error represents an error state
	Error
)

// Server represents a VPN server
type Server struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Hostname string    `json:"hostname"`
	IP       string    `json:"ip"`
	Port     int       `json:"port"`
	Country  string    `json:"country"`
	City     string    `json:"city"`
	Ping     int       `json:"ping"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// ConnectionInfo holds information about the current connection
type ConnectionInfo struct {
	ServerID      string        `json:"server_id"`
	StartTime     time.Time     `json:"start_time"`
	DataSent      int64         `json:"data_sent"`
	DataReceived  int64         `json:"data_received"`
	Status        ConnectionStatus `json:"status"`
	LastError     string        `json:"last_error,omitempty"`
	EncryptionKey string        `json:"encryption_key,omitempty"`
}