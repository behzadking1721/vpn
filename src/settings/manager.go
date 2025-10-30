package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// SettingsManager handles dashboard settings
type SettingsManager struct {
	config SettingsManagerConfig
	settings *DashboardSettings
	mutex sync.RWMutex
}

// NewSettingsManager creates a new settings manager
func NewSettingsManager(config SettingsManagerConfig) *SettingsManager {
	sm := &SettingsManager{
		config: config,
		settings: &DashboardSettings{
			Theme: ThemeSystem,
			ChartWindow: Window24Hours,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	
	// Load existing settings if they exist
	sm.loadSettings()
	
	return sm
}

// GetSettings returns the current dashboard settings
func (sm *SettingsManager) GetSettings() *DashboardSettings {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	// Return a copy to prevent external modification
	settingsCopy := *sm.settings
	return &settingsCopy
}

// UpdateSettings updates the dashboard settings
func (sm *SettingsManager) UpdateSettings(newSettings DashboardSettings) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// Update the settings
	newSettings.UpdatedAt = time.Now()
	sm.settings = &newSettings
	
	// Save the settings
	return sm.saveSettings()
}

// loadSettings loads settings from storage
func (sm *SettingsManager) loadSettings() error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// Check if the settings file exists
	if _, err := os.Stat(sm.config.StoragePath); os.IsNotExist(err) {
		// File doesn't exist, use default settings
		return nil
	}
	
	// Read the settings file
	data, err := ioutil.ReadFile(sm.config.StoragePath)
	if err != nil {
		return err
	}
	
	// Parse the settings
	var settings DashboardSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}
	
	// Validate the settings
	if err := sm.validateSettings(&settings); err != nil {
		// If validation fails, use default settings
		return nil
	}
	
	// Use the loaded settings
	sm.settings = &settings
	return nil
}

// saveSettings saves settings to storage
func (sm *SettingsManager) saveSettings() error {
	// Create the directory if it doesn't exist
	dir := sm.config.StoragePath[:len(sm.config.StoragePath)-len("dashboard_settings.json")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Serialize the settings
	data, err := json.MarshalIndent(sm.settings, "", "  ")
	if err != nil {
		return err
	}
	
	// Write the settings to file
	return ioutil.WriteFile(sm.config.StoragePath, data, 0644)
}

// validateSettings validates the dashboard settings
func (sm *SettingsManager) validateSettings(settings *DashboardSettings) error {
	// Validate theme
	switch settings.Theme {
	case ThemeLight, ThemeDark, ThemeSystem:
		// Valid theme
	default:
		// Invalid theme, use default
		settings.Theme = ThemeSystem
	}
	
	// Validate chart window
	switch settings.ChartWindow {
	case Window24Hours, Window7Days, Window30Days:
		// Valid window
	default:
		// Invalid window, use default
		settings.ChartWindow = Window24Hours
	}
	
	return nil
}