package settings

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSettingsManager(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "dashboard_settings_test.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Close the file as we only need the path
	tempFile.Close()

	// Create settings manager config
	config := SettingsManagerConfig{
		StoragePath: tempFile.Name(),
	}

	// Create settings manager
	settingsManager := NewSettingsManager(config)

	// Test default settings
	t.Run("DefaultSettings", func(t *testing.T) {
		settings := settingsManager.GetSettings()

		if settings.Theme != ThemeSystem {
			t.Errorf("Expected default theme to be 'system', got '%s'", settings.Theme)
		}

		if settings.ChartWindow != Window24Hours {
			t.Errorf("Expected default chart window to be '24h', got '%s'", settings.ChartWindow)
		}
	})

	// Test updating settings
	t.Run("UpdateSettings", func(t *testing.T) {
		// Create new settings
		newSettings := DashboardSettings{
			Theme:       ThemeDark,
			ChartWindow: Window7Days,
		}

		// Update settings
		if err := settingsManager.UpdateSettings(newSettings); err != nil {
			t.Fatalf("Failed to update settings: %v", err)
		}

		// Get updated settings
		updatedSettings := settingsManager.GetSettings()

		if updatedSettings.Theme != ThemeDark {
			t.Errorf("Expected theme to be 'dark', got '%s'", updatedSettings.Theme)
		}

		if updatedSettings.ChartWindow != Window7Days {
			t.Errorf("Expected chart window to be '7d', got '%s'", updatedSettings.ChartWindow)
		}
	})

	// Test persistence
	t.Run("SettingsPersistence", func(t *testing.T) {
		// Create a new settings manager with the same config
		newSettingsManager := NewSettingsManager(config)

		// Get settings from the new manager
		settings := newSettingsManager.GetSettings()

		if settings.Theme != ThemeDark {
			t.Errorf("Expected persisted theme to be 'dark', got '%s'", settings.Theme)
		}

		if settings.ChartWindow != Window7Days {
			t.Errorf("Expected persisted chart window to be '7d', got '%s'", settings.ChartWindow)
		}
	})

	// Test validation
	t.Run("SettingsValidation", func(t *testing.T) {
		// Create invalid settings
		invalidSettings := DashboardSettings{
			Theme:       "invalid-theme",
			ChartWindow: "invalid-window",
		}

		// Update settings (should still work but with defaults)
		if err := settingsManager.UpdateSettings(invalidSettings); err != nil {
			t.Fatalf("Failed to update settings: %v", err)
		}

		// Get updated settings
		updatedSettings := settingsManager.GetSettings()

		// Should have defaulted to valid values
		if updatedSettings.Theme != ThemeSystem {
			t.Errorf("Expected theme to default to 'system', got '%s'", updatedSettings.Theme)
		}

		if updatedSettings.ChartWindow != Window24Hours {
			t.Errorf("Expected chart window to default to '24h', got '%s'", updatedSettings.ChartWindow)
		}
	})
}

func TestThemeConstants(t *testing.T) {
	// Test that theme constants have the expected values
	expectedThemes := map[DashboardTheme]string{
		ThemeLight:  "light",
		ThemeDark:   "dark",
		ThemeSystem: "system",
	}

	for theme, expectedValue := range expectedThemes {
		if string(theme) != expectedValue {
			t.Errorf("Expected theme %s to have value '%s', got '%s'", theme, expectedValue, string(theme))
		}
	}
}

func TestChartWindowConstants(t *testing.T) {
	// Test that chart window constants have the expected values
	expectedWindows := map[ChartWindow]string{
		Window24Hours: "24h",
		Window7Days:   "7d",
		Window30Days:  "30d",
	}

	for window, expectedValue := range expectedWindows {
		if string(window) != expectedValue {
			t.Errorf("Expected window %s to have value '%s', got '%s'", window, expectedValue, string(window))
		}
	}
}
