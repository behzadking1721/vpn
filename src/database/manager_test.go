package database

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestDatabaseManager(t *testing.T) {
	// Create a temporary database file for testing
	tempFile, err := ioutil.TempFile("", "vpn_test.db")
	if err != nil {
		t.Fatalf("Failed to create temporary database file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Close the file as we only need the path
	tempFile.Close()

	// Create database manager
	manager, err := NewManager(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}
	defer manager.Close()

	// Test connection records
	t.Run("ConnectionRecords", func(t *testing.T) {
		// Add a connection record
		record := ConnectionRecord{
			ID:               utils.GenerateID(),
			ServerID:         "server1",
			ServerName:       "Test Server",
			StartTime:        time.Now(),
			EndTime:          time.Now().Add(10 * time.Minute),
			Duration:         600,
			DataSent:         1024 * 1024,
			DataReceived:     2048 * 1024,
			Protocol:         "VMess",
			Status:           "connected",
			DisconnectReason: "",
		}

		if err := manager.AddConnectionRecord(record); err != nil {
			t.Fatalf("Failed to add connection record: %v", err)
		}

		// Retrieve connection records
		records, err := manager.GetConnectionRecords(10, 0)
		if err != nil {
			t.Fatalf("Failed to get connection records: %v", err)
		}

		if len(records) != 1 {
			t.Errorf("Expected 1 connection record, got %d", len(records))
		}

		if records[0].ID != record.ID {
			t.Errorf("Expected record ID %s, got %s", record.ID, records[0].ID)
		}

		if records[0].ServerName != record.ServerName {
			t.Errorf("Expected server name %s, got %s", record.ServerName, records[0].ServerName)
		}
	})

	// Test data usage records
	t.Run("DataUsageRecords", func(t *testing.T) {
		// Add a data usage record
		record := DataUsageRecord{
			ID:            utils.GenerateID(),
			Timestamp:     time.Now(),
			ServerID:      "server1",
			ServerName:    "Test Server",
			DataSent:      1024 * 1024,
			DataReceived:  2048 * 1024,
			TotalSent:     1024 * 1024 * 10,
			TotalReceived: 2048 * 1024 * 10,
		}

		if err := manager.AddDataUsageRecord(record); err != nil {
			t.Fatalf("Failed to add data usage record: %v", err)
		}

		// Retrieve data usage records
		records, err := manager.GetDataUsageRecords("server1", 10, 0)
		if err != nil {
			t.Fatalf("Failed to get data usage records: %v", err)
		}

		if len(records) != 1 {
			t.Errorf("Expected 1 data usage record, got %d", len(records))
		}

		if records[0].ID != record.ID {
			t.Errorf("Expected record ID %s, got %s", record.ID, records[0].ID)
		}

		if records[0].ServerName != record.ServerName {
			t.Errorf("Expected server name %s, got %s", record.ServerName, records[0].ServerName)
		}
	})

	// Test ping records
	t.Run("PingRecords", func(t *testing.T) {
		// Add a ping record
		record := PingRecord{
			ID:         utils.GenerateID(),
			Timestamp:  time.Now(),
			ServerID:   "server1",
			ServerName: "Test Server",
			Ping:       50,
		}

		if err := manager.AddPingRecord(record); err != nil {
			t.Fatalf("Failed to add ping record: %v", err)
		}

		// Retrieve ping records
		records, err := manager.GetPingRecords("server1", 10)
		if err != nil {
			t.Fatalf("Failed to get ping records: %v", err)
		}

		if len(records) != 1 {
			t.Errorf("Expected 1 ping record, got %d", len(records))
		}

		if records[0].ID != record.ID {
			t.Errorf("Expected record ID %s, got %s", record.ID, records[0].ID)
		}

		if records[0].ServerName != record.ServerName {
			t.Errorf("Expected server name %s, got %s", record.ServerName, records[0].ServerName)
		}

		if records[0].Ping != record.Ping {
			t.Errorf("Expected ping %d, got %d", record.Ping, records[0].Ping)
		}
	})

	// Test alert records
	t.Run("AlertRecords", func(t *testing.T) {
		// Add an alert record
		record := AlertRecord{
			ID:         utils.GenerateID(),
			Type:       "data_usage",
			Title:      "High Data Usage",
			Message:    "Data usage exceeded 80% threshold",
			Timestamp:  time.Now(),
			Value:      85.5,
			Severity:   "warning",
			Resolved:   false,
			Read:       false,
			ServerID:   "server1",
			ServerName: "Test Server",
		}

		if err := manager.AddAlertRecord(record); err != nil {
			t.Fatalf("Failed to add alert record: %v", err)
		}

		// Retrieve alert records
		records, err := manager.GetAlertRecords(false, false, 10)
		if err != nil {
			t.Fatalf("Failed to get alert records: %v", err)
		}

		if len(records) != 1 {
			t.Errorf("Expected 1 alert record, got %d", len(records))
		}

		if records[0].ID != record.ID {
			t.Errorf("Expected record ID %s, got %s", record.ID, records[0].ID)
		}

		if records[0].Title != record.Title {
			t.Errorf("Expected title %s, got %s", record.Title, records[0].Title)
		}

		if records[0].Severity != record.Severity {
			t.Errorf("Expected severity %s, got %s", record.Severity, records[0].Severity)
		}

		// Test retrieving unread alerts
		unreadRecords, err := manager.GetAlertRecords(true, false, 10)
		if err != nil {
			t.Fatalf("Failed to get unread alert records: %v", err)
		}

		if len(unreadRecords) != 1 {
			t.Errorf("Expected 1 unread alert record, got %d", len(unreadRecords))
		}

		// Test retrieving unresolved alerts
		unresolvedRecords, err := manager.GetAlertRecords(false, true, 10)
		if err != nil {
			t.Fatalf("Failed to get unresolved alert records: %v", err)
		}

		if len(unresolvedRecords) != 1 {
			t.Errorf("Expected 1 unresolved alert record, got %d", len(unresolvedRecords))
		}

		// Update alert record
		record.Resolved = true
		record.Read = true

		if err := manager.UpdateAlertRecord(record); err != nil {
			t.Fatalf("Failed to update alert record: %v", err)
		}

		// Retrieve updated alert records
		updatedRecords, err := manager.GetAlertRecords(true, true, 10)
		if err != nil {
			t.Fatalf("Failed to get updated alert records: %v", err)
		}

		if len(updatedRecords) != 0 {
			t.Errorf("Expected 0 records with unread=true and unresolved=true, got %d", len(updatedRecords))
		}
	})

	// Test dashboard settings
	t.Run("DashboardSettings", func(t *testing.T) {
		// Get default settings
		settings, err := manager.GetDashboardSettings()
		if err != nil {
			t.Fatalf("Failed to get dashboard settings: %v", err)
		}

		if settings.Theme != "system" {
			t.Errorf("Expected default theme 'system', got '%s'", settings.Theme)
		}

		if settings.ChartWindow != "24h" {
			t.Errorf("Expected default chart window '24h', got '%s'", settings.ChartWindow)
		}

		// Update settings
		settings.Theme = "dark"
		settings.ChartWindow = "7d"
		settings.UpdatedAt = time.Now()

		if err := manager.UpdateDashboardSettings(*settings); err != nil {
			t.Fatalf("Failed to update dashboard settings: %v", err)
		}

		// Retrieve updated settings
		updatedSettings, err := manager.GetDashboardSettings()
		if err != nil {
			t.Fatalf("Failed to get updated dashboard settings: %v", err)
		}

		if updatedSettings.Theme != "dark" {
			t.Errorf("Expected updated theme 'dark', got '%s'", updatedSettings.Theme)
		}

		if updatedSettings.ChartWindow != "7d" {
			t.Errorf("Expected updated chart window '7d', got '%s'", updatedSettings.ChartWindow)
		}
	})
}

func TestMigrationManager(t *testing.T) {
	// Create temporary directories for testing
	tempDir, err := ioutil.TempDir("", "vpn_migration_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	historyPath := tempDir + "/history"
	settingsPath := tempDir + "/settings"

	// Create a temporary database file for testing
	dbFile, err := ioutil.TempFile("", "vpn_test.db")
	if err != nil {
		t.Fatalf("Failed to create temporary database file: %v", err)
	}
	defer os.Remove(dbFile.Name())

	// Close the file as we only need the path
	dbFile.Close()

	// Create database manager
	dbManager, err := NewManager(dbFile.Name())
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}
	defer dbManager.Close()

	// Create migration manager
	migrationManager := NewMigrationManager(dbManager)

	// Test backup functionality
	t.Run("BackupJSONData", func(t *testing.T) {
		// Create test directories and files
		if err := os.MkdirAll(historyPath, 0755); err != nil {
			t.Fatalf("Failed to create history directory: %v", err)
		}

		if err := os.MkdirAll(settingsPath, 0755); err != nil {
			t.Fatalf("Failed to create settings directory: %v", err)
		}

		// Create a test file in history directory
		historyFile := historyPath + "/test_history.json"
		if err := ioutil.WriteFile(historyFile, []byte("{}"), 0644); err != nil {
			t.Fatalf("Failed to create test history file: %v", err)
		}

		// Create a test file in settings directory
		settingsFile := settingsPath + "/test_settings.json"
		if err := ioutil.WriteFile(settingsFile, []byte("{}"), 0644); err != nil {
			t.Fatalf("Failed to create test settings file: %v", err)
		}

		// Perform backup
		if err := migrationManager.BackupJSONData(historyPath, settingsPath); err != nil {
			t.Fatalf("Failed to backup JSON data: %v", err)
		}

		// Check if backup was created
		backupDir := tempDir + "/backup"
		if _, err := os.Stat(backupDir); os.IsNotExist(err) {
			t.Error("Backup directory was not created")
		}
	})
}
