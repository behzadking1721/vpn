package alert

import (
	"errors"
	"time"
	"vpnclient/src/core"
)

// AlertRuleType represents the type of alert rule
type AlertRuleType string

const (
	// RuleTypeDataUsage represents data usage threshold rule
	RuleTypeDataUsage AlertRuleType = "data_usage"

	// RuleTypeHighPing represents high ping threshold rule
	RuleTypeHighPing AlertRuleType = "high_ping"

	// RuleTypeConnectionLoss represents connection loss rule
	RuleTypeConnectionLoss AlertRuleType = "connection_loss"

	// RuleTypeServerUnreachable represents server unreachable rule
	RuleTypeServerUnreachable AlertRuleType = "server_unreachable"
)

// AlertSeverity represents the severity level of an alert
type AlertSeverity string

const (
	// SeverityInfo represents informational alerts
	SeverityInfo AlertSeverity = "info"

	// SeverityWarning represents warning alerts
	SeverityWarning AlertSeverity = "warning"

	// SeverityError represents error alerts
	SeverityError AlertSeverity = "error"

	// SeverityCritical represents critical alerts
	SeverityCritical AlertSeverity = "critical"
)

// AlertRule defines a rule for generating alerts
type AlertRule struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        AlertRuleType `json:"type"`
	Enabled     bool          `json:"enabled"`

	// Rule specific parameters
	Threshold float64 `json:"threshold,omitempty"` // Percentage for data usage, milliseconds for ping
	Duration  int     `json:"duration,omitempty"`  // Duration in seconds for connection loss
	ServerID  string  `json:"server_id,omitempty"` // Specific server ID or empty for all servers

	// Notification settings
	NotifyDesktop bool      `json:"notify_desktop"`
	NotifyUI      bool      `json:"notify_ui"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Alert represents an alert instance
type Alert struct {
	ID         string        `json:"id"`
	RuleID     string        `json:"rule_id"`
	Type       AlertRuleType `json:"type"`
	Severity   AlertSeverity `json:"severity"`
	Title      string        `json:"title"`
	Message    string        `json:"message"`
	ServerID   string        `json:"server_id,omitempty"`
	ServerName string        `json:"server_name,omitempty"`
	Value      float64       `json:"value,omitempty"`     // Actual value that triggered the alert
	Threshold  float64       `json:"threshold,omitempty"` // Threshold that was exceeded
	Timestamp  time.Time     `json:"timestamp"`
	Read       bool          `json:"read"`
	Resolved   bool          `json:"resolved"`
	ResolvedAt *time.Time    `json:"resolved_at,omitempty"`
}

// AlertEvaluationContext contains context for evaluating alert rules
type AlertEvaluationContext struct {
	// Connection related data
	ConnectionStatus core.ConnectionStatus
	ConnectionInfo   core.ConnectionInfo
	CurrentServer    *core.Server

	// Data usage data
	DataSent     int64
	DataReceived int64

	// Server data
	Servers []core.Server

	// Time context
	Now time.Time
}

// AlertManagerConfig holds configuration for the alert manager
type AlertManagerConfig struct {
	// Enable or disable desktop notifications
	DesktopNotifications bool

	// Evaluation interval in seconds
	EvaluationInterval int

	// History retention in days
	HistoryRetention int
}

// Alert handler interface
type AlertHandler interface {
	HandleAlert(alert *Alert)
}

// Error definitions
var (
	ErrRuleNotFound      = errors.New("alert rule not found")
	ErrRuleAlreadyExists = errors.New("alert rule already exists")
)
