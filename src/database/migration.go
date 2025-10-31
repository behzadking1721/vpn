package database

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/alert"
	"c:/Users/behza/OneDrive/Documents/vpn/src/history"
	"c:/Users/behza/OneDrive/Documents/vpn/src/settings"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// MigrationManager handles migration from JSON storage to SQLite database
type MigrationManager struct {
	dbManager *Manager
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(dbManager *Manager) *MigrationManager {
	return &MigrationManager{
		dbManager: dbManager,
	}
}

// MigrateFromJSON migrates data from JSON files to SQLite database
func (m *MigrationManager) MigrateFromJSON(historyPath, settingsPath string) error {
	// Migrate history data
	if err := m.migrateHistoryData(historyPath); err != nil {
		return fmt.Errorf("failed to migrate history data: %w", err)
	}

	// Migrate settings data
	if err := m.migrateSettingsData(settingsPath); err != nil {
		return fmt.Errorf("failed to migrate settings data: %w", err)
	}

	// Migrate alert data (if exists in JSON format)
	// This would depend on how alerts were previously stored

	return nil
}

// migrateHistoryData migrates history data from JSON files
func (m *MigrationManager) migrateHistoryData(historyPath string) error {
	// Check if history path exists
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		// No history data to migrate
		return nil
	}

	// Create history manager to read existing data
	historyManager := history.NewHistoryManager(historyPath)

	// Migrate connection records
	connectionRecords, err := historyManager.GetConnectionRecords(0, 0)
	if err != nil {
		return fmt.Errorf("failed to get connection records: %w", err)
	}

	for _, record := range connectionRecords {
		dbRecord := ConnectionRecord{
			ID:               record.ID,
			ServerID:         record.ServerID,
			ServerName:       record.ServerName,
			StartTime:        record.StartTime,
			EndTime:          record.EndTime,
			Duration:         record.Duration,
			DataSent:         record.DataSent,
			DataReceived:     record.DataReceived,
			Protocol:         record.Protocol,
			Status:           record.Status,
			DisconnectReason: record.DisconnectReason,
		}

		if err := m.dbManager.AddConnectionRecord(dbRecord); err != nil {
			return fmt.Errorf("failed to add connection record %s: %w", record.ID, err)
		}
	}

	// Migrate data usage records
	dataUsageRecords, err := historyManager.GetDataUsageRecords("", 0, 0)
	if err != nil {
		return fmt.Errorf("failed to get data usage records: %w", err)
	}

	for _, record := range dataUsageRecords {
		dbRecord := DataUsageRecord{
			ID:            record.ID,
			Timestamp:     record.Timestamp,
			ServerID:      record.ServerID,
			ServerName:    record.ServerName,
			DataSent:      record.DataSent,
			DataReceived:  record.DataReceived,
			TotalSent:     record.TotalSent,
			TotalReceived: record.TotalReceived,
		}

		if err := m.dbManager.AddDataUsageRecord(dbRecord); err != nil {
			return fmt.Errorf("failed to add data usage record %s: %w", record.ID, err)
		}
	}

	return nil
}

// migrateSettingsData migrates settings data from JSON file
func (m *MigrationManager) migrateSettingsData(settingsPath string) error {
	// Check if settings file exists
	settingsFile := filepath.Join(settingsPath, "dashboard_settings.json")
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		// No settings file to migrate
		return nil
	}

	// Read settings file
	data, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		return fmt.Errorf("failed to read settings file: %w", err)
	}

	// Parse settings
	var jsonSettings settings.DashboardSettings
	if err := json.Unmarshal(data, &jsonSettings); err != nil {
		return fmt.Errorf("failed to parse settings file: %w", err)
	}

	// Convert to database record
	dbSettings := DashboardSettingsRecord{
		Theme:       string(jsonSettings.Theme),
		ChartWindow: string(jsonSettings.ChartWindow),
		CreatedAt:   jsonSettings.CreatedAt,
		UpdatedAt:   jsonSettings.UpdatedAt,
	}

	// Update database settings
	if err := m.dbManager.UpdateDashboardSettings(dbSettings); err != nil {
		return fmt.Errorf("failed to update dashboard settings: %w", err)
	}

	return nil
}

// migrateAlertData migrates alert data from JSON files (if stored separately)
func (m *MigrationManager) migrateAlertData(alertPath string) error {
	// This would depend on how alerts were previously stored
	// For now, we'll assume alerts are handled by the alert manager

	// In a real implementation, you might need to:
	// 1. Read alert data from JSON files
	// 2. Convert to database records
	// 3. Insert into the database

	return nil
}

// BackupJSONData creates a backup of JSON data before migration
func (m *MigrationManager) BackupJSONData(historyPath, settingsPath string) error {
	backupDir := filepath.Join(filepath.Dir(historyPath), "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Backup history data
	historyBackup := filepath.Join(backupDir, fmt.Sprintf("history_backup_%d", time.Now().Unix()))
	if err := copyDir(historyPath, historyBackup); err != nil {
		return fmt.Errorf("failed to backup history data: %w", err)
	}

	// Backup settings data
	settingsBackup := filepath.Join(backupDir, fmt.Sprintf("settings_backup_%d", time.Now().Unix()))
	if err := copyDir(settingsPath, settingsBackup); err != nil {
		return fmt.Errorf("failed to backup settings data: %w", err)
	}

	return nil
}

// copyDir copies a directory recursively
func copyDir(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read source directory entries
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	// Read source file
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// Write destination file
	return ioutil.WriteFile(dst, data, 0644)
}
