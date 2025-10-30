package settings

import (
	"time"
)

// DashboardTheme represents the theme of the dashboard
type DashboardTheme string

const (
	// ThemeLight represents the light theme
	ThemeLight DashboardTheme = "light"
	
	// ThemeDark represents the dark theme
	ThemeDark DashboardTheme = "dark"
	
	// ThemeSystem represents the system theme (follows OS theme)
	ThemeSystem DashboardTheme = "system"
)

// ChartWindow represents the time window for charts
type ChartWindow string

const (
	// Window24Hours represents a 24-hour window
	Window24Hours ChartWindow = "24h"
	
	// Window7Days represents a 7-day window
	Window7Days ChartWindow = "7d"
	
	// Window30Days represents a 30-day window
	Window30Days ChartWindow = "30d"
)

// DashboardSettings represents the dashboard settings
type DashboardSettings struct {
	// Theme represents the dashboard theme
	Theme DashboardTheme `json:"theme"`
	
	// ChartWindow represents the time window for charts
	ChartWindow ChartWindow `json:"chart_window"`
	
	// CreatedAt represents when the settings were created
	CreatedAt time.Time `json:"created_at"`
	
	// UpdatedAt represents when the settings were last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// SettingsManagerConfig holds configuration for the settings manager
type SettingsManagerConfig struct {
	// StoragePath is the path to the settings storage file
	StoragePath string
}