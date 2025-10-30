package core

import (
	"time"
)

// ProtocolType represents the supported VPN protocols
type ProtocolType string

const (
	ProtocolVMess     ProtocolType = "vmess"
	ProtocolVLESS     ProtocolType = "vless"
	ProtocolTrojan    ProtocolType = "trojan"
	ProtocolReality   ProtocolType = "reality"
	ProtocolHysteria  ProtocolType = "hysteria2"
	ProtocolTUIC      ProtocolType = "tuic"
	ProtocolSSH       ProtocolType = "ssh"
	ProtocolShadowsocks ProtocolType = "shadowsocks"
)

// Server represents a VPN server configuration
type Server struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Host           string       `json:"host"`
	Port           int          `json:"port"`
	Protocol       ProtocolType `json:"protocol"`
	Username       string       `json:"username,omitempty"`
	Password       string       `json:"password,omitempty"`
	Encryption     string       `json:"encryption,omitempty"`
	TLS            bool         `json:"tls,omitempty"`
	SNI            string       `json:"sni,omitempty"`
	Fingerprint    string       `json:"fingerprint,omitempty"`
	UDP            bool         `json:"udp,omitempty"`
	Remark         string       `json:"remark,omitempty"`
	Ping           int          `json:"ping,omitempty"` // in milliseconds
	LastPing       time.Time    `json:"last_ping,omitempty"`
	DataUsed       int64        `json:"data_used,omitempty"`  // in bytes
	DataLimit      int64        `json:"data_limit,omitempty"` // in bytes
	Enabled        bool         `json:"enabled"`
	Method         string       `json:"method,omitempty"` // For Shadowsocks
	Obfs           string       `json:"obfs,omitempty"`   // For Shadowsocks
	ObfsParam      string       `json:"obfs_param,omitempty"` // For Shadowsocks
	ProtocolParam  string       `json:"protocol_param,omitempty"` // For Shadowsocks
}

// Subscription represents a server subscription
type Subscription struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	LastUpdate  time.Time `json:"last_update,omitempty"`
	AutoUpdate  bool      `json:"auto_update"`
	Servers     []Server  `json:"servers"`
}

// ConnectionStatus represents the current connection status
type ConnectionStatus string

const (
	StatusDisconnected  ConnectionStatus = "disconnected"
	StatusConnecting    ConnectionStatus = "connecting"
	StatusConnected     ConnectionStatus = "connected"
	StatusDisconnecting ConnectionStatus = "disconnecting"
	StatusError         ConnectionStatus = "error"
)

// TunnelMode represents the tunneling modes
type TunnelMode string

const (
	TunnelModeTCP     TunnelMode = "tcp"
	TunnelModeUDP     TunnelMode = "udp"
	TunnelModeBoth    TunnelMode = "both"
	TunnelModeAuto    TunnelMode = "auto"
)

// ConnectionInfo holds information about the current connection
type ConnectionInfo struct {
	Status        ConnectionStatus `json:"status"`
	ServerID      string           `json:"server_id,omitempty"`
	StartTime     time.Time        `json:"start_time,omitempty"`
	EndTime       time.Time        `json:"end_time,omitempty"`
	DataSent      int64            `json:"data_sent,omitempty"`  // in bytes
	DataReceived  int64            `json:"data_received,omitempty"` // in bytes
	LocalIP       string           `json:"local_ip,omitempty"`
	RemoteIP      string           `json:"remote_ip,omitempty"`
	Error         string           `json:"error,omitempty"`
	TunnelMode    TunnelMode       `json:"tunnel_mode,omitempty"`
}

// DataUsage represents data usage statistics
type DataUsage struct {
	ServerID     string    `json:"server_id"`
	DataSent     int64     `json:"data_sent"`     // in bytes
	DataReceived int64     `json:"data_received"` // in bytes
	LastUpdate   time.Time `json:"last_update"`
}

// AppConfig holds application configuration
type AppConfig struct {
	AutoConnect         bool `json:"auto_connect"`
	AutoLaunch          bool `json:"auto_launch"`
	AutoUpdate          bool `json:"auto_update"`
	StartMinimized      bool `json:"start_minimized"`
	ShowNotifications   bool `json:"show_notifications"`
	IPv6Support         bool `json:"ipv6_support"`
	SystemProxy         bool `json:"system_proxy"`
	BypassLAN           bool `json:"bypass_lan"`
	CloseToSystemTray   bool `json:"close_to_system_tray"`
	TunnelMode          TunnelMode `json:"tunnel_mode"` // Default tunneling mode
}