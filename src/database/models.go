package database

import (
	"database/sql"
	"time"
)

// DB represents the database connection
type DB struct {
	*sql.DB
}

// ConnectionRecord represents a connection record in the database
type ConnectionRecord struct {
	ID               string    `json:"id" db:"id"`
	ServerID         string    `json:"server_id" db:"server_id"`
	ServerName       string    `json:"server_name" db:"server_name"`
	StartTime        time.Time `json:"start_time" db:"start_time"`
	EndTime          time.Time `json:"end_time" db:"end_time"`
	Duration         int       `json:"duration" db:"duration"`
	DataSent         int64     `json:"data_sent" db:"data_sent"`
	DataReceived     int64     `json:"data_received" db:"data_received"`
	Protocol         string    `json:"protocol" db:"protocol"`
	Status           string    `json:"status" db:"status"`
	DisconnectReason string    `json:"disconnect_reason" db:"disconnect_reason"`
}

// DataUsageRecord represents a data usage record in the database
type DataUsageRecord struct {
	ID            string    `json:"id" db:"id"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	ServerID      string    `json:"server_id" db:"server_id"`
	ServerName    string    `json:"server_name" db:"server_name"`
	DataSent      int64     `json:"data_sent" db:"data_sent"`
	DataReceived  int64     `json:"data_received" db:"data_received"`
	TotalSent     int64     `json:"total_sent" db:"total_sent"`
	TotalReceived int64     `json:"total_received" db:"total_received"`
}

// PingRecord represents a ping record in the database
type PingRecord struct {
	ID        string    `json:"id" db:"id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	ServerID  string    `json:"server_id" db:"server_id"`
	ServerName string   `json:"server_name" db:"server_name"`
	Ping      int       `json:"ping" db:"ping"`
}

// AlertRecord represents an alert record in the database
type AlertRecord struct {
	ID          string    `json:"id" db:"id"`
	Type        string    `json:"type" db:"type"`
	Title       string    `json:"title" db:"title"`
	Message     string    `json:"message" db:"message"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	Value       float64   `json:"value" db:"value"`
	Severity    string    `json:"severity" db:"severity"`
	Resolved    bool      `json:"resolved" db:"resolved"`
	Read        bool      `json:"read" db:"read"`
	ServerID    string    `json:"server_id" db:"server_id"`
	ServerName  string    `json:"server_name" db:"server_name"`
}

// DashboardSettingsRecord represents dashboard settings in the database
type DashboardSettingsRecord struct {
	ID          int       `json:"id" db:"id"`
	Theme       string    `json:"theme" db:"theme"`
	ChartWindow string    `json:"chart_window" db:"chart_window"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Migration represents a database migration
type Migration struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	AppliedAt time.Time `json:"applied_at" db:"applied_at"`
}