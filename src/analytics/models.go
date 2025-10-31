package analytics

import (
	"time"
)

// PingStats represents ping statistics
type PingStats struct {
	Average float64 `json:"average"`
	Max     int     `json:"max"`
	P95     int     `json:"p95"`
	Min     int     `json:"min"`
	Samples int     `json:"samples"`
}

// DataUsageStats represents data usage statistics
type DataUsageStats struct {
	TotalSent     int64 `json:"total_sent"`
	TotalReceived int64 `json:"total_received"`
	AveragePerDay int64 `json:"average_per_day"`
	MaxPerDay     int64 `json:"max_per_day"`
}

// TimePattern represents time-based patterns
type TimePattern struct {
	Hour        int `json:"hour"`
	UsageCount  int `json:"usage_count"`
	Disconnects int `json:"disconnects"`
}

// ReportPeriod represents the period of a report
type ReportPeriod string

const (
	ReportPeriodWeekly  ReportPeriod = "weekly"
	ReportPeriodMonthly ReportPeriod = "monthly"
)

// AnalyticsReport represents a complete analytics report
type AnalyticsReport struct {
	ID          string         `json:"id"`
	Period      ReportPeriod   `json:"period"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	PingStats   PingStats      `json:"ping_stats"`
	DataUsage   DataUsageStats `json:"data_usage"`
	TimePattern []TimePattern  `json:"time_pattern"`
	CreatedAt   time.Time      `json:"created_at"`
}

// PingAnalysisRequest represents a request for ping analysis
type PingAnalysisRequest struct {
	Window string `json:"window"` // e.g., "7d", "30d"
}

// DataUsageAnalysisRequest represents a request for data usage analysis
type DataUsageAnalysisRequest struct {
	Window string `json:"window"` // e.g., "7d", "30d"
}

// ReportRequest represents a request for a report
type ReportRequest struct {
	Period ReportPeriod `json:"period"`
}

// DailyDataUsage represents daily data usage
type DailyDataUsage struct {
	Date         time.Time `json:"date"`
	DataSent     int64     `json:"data_sent"`
	DataReceived int64     `json:"data_received"`
}

// WeeklyDataUsage represents weekly data usage
type WeeklyDataUsage struct {
	WeekStart    time.Time `json:"week_start"`
	WeekEnd      time.Time `json:"week_end"`
	DataSent     int64     `json:"data_sent"`
	DataReceived int64     `json:"data_received"`
}

// MonthlyDataUsage represents monthly data usage
type MonthlyDataUsage struct {
	MonthStart   time.Time `json:"month_start"`
	MonthEnd     time.Time `json:"month_end"`
	DataSent     int64     `json:"data_sent"`
	DataReceived int64     `json:"data_received"`
}

// ConnectionPattern represents connection patterns
type ConnectionPattern struct {
	Hour        int `json:"hour"`
	Connections int `json:"connections"`
	Disconnects int `json:"disconnects"`
}

// ServerPerformance represents server performance metrics
type ServerPerformance struct {
	ServerID   string    `json:"server_id"`
	ServerName string    `json:"server_name"`
	PingStats  PingStats `json:"ping_stats"`
}
