package core

import "time"

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

// ConnectionInfo holds connection information
type ConnectionInfo struct {
	ID        string    `json:"id"`
	StartedAt time.Time `json:"started_at"`
	DataSent  int64     `json:"data_sent"`
	DataRecv  int64     `json:"data_recv"`
}

// Server represents a VPN server
type Server struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Host        string                 `json:"host"`
	Port        int                    `json:"port"`
	Protocol    string                 `json:"protocol"`
	Country     string                 `json:"country,omitempty"`
	Enabled     bool                   `json:"enabled"`
	Ping        int                    `json:"ping"`        // in milliseconds
	Config      map[string]interface{} `json:"config"`      // protocol-specific configuration
	Remark      string                 `json:"remark,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Subscription represents a subscription link
type Subscription struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	AutoUpdate  bool      `json:"auto_update"`
	LastUpdate  time.Time `json:"last_update,omitempty"`
	ServerCount int       `json:"server_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}