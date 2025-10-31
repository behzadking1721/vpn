package managers

import (
	"vpn-client/src/core"
	"encoding/json"
	"os"
	"path/filepath"
)

// ConfigManager handles application configuration
type ConfigManager struct {
	config core.AppConfig
	path   string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configPath string) *ConfigManager {
	cm := &ConfigManager{
		config: core.AppConfig{
			AutoConnect:       false,
			AutoLaunch:        false,
			AutoUpdate:        true,
			StartMinimized:    false,
			ShowNotifications: true,
			IPv6Support:       true,
			SystemProxy:       true,
			BypassLAN:         true,
			CloseToSystemTray: true,
			TunnelMode:        core.TunnelModeBoth,
		},
		path: configPath,
	}

	// Try to load existing configuration
	cm.LoadConfig()
	return cm
}

// LoadConfig loads configuration from file
func (cm *ConfigManager) LoadConfig() error {
	// Ensure config directory exists
	dir := filepath.Dir(cm.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	// Read config file
	data, err := os.ReadFile(cm.path)
	if err != nil {
		// If file doesn't exist, use defaults
		return nil
	}

	// Parse config
	err = json.Unmarshal(data, &cm.config)
	if err != nil {
		return err
	}

	return nil
}

// SaveConfig saves configuration to file
func (cm *ConfigManager) SaveConfig() error {
	// Ensure config directory exists
	dir := filepath.Dir(cm.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	// Serialize config
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(cm.path, data, 0644)
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() core.AppConfig {
	// Return a copy to prevent external modification
	config := cm.config
	return config
}

// UpdateConfig updates the configuration
func (cm *ConfigManager) UpdateConfig(newConfig core.AppConfig) error {
	cm.config = newConfig
	return cm.SaveConfig()
}

// GetTunnelMode returns the current tunnel mode
func (cm *ConfigManager) GetTunnelMode() core.TunnelMode {
	return cm.config.TunnelMode
}

// SetTunnelMode sets the tunnel mode
func (cm *ConfigManager) SetTunnelMode(mode core.TunnelMode) error {
	cm.config.TunnelMode = mode
	return cm.SaveConfig()
}
