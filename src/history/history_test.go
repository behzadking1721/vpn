package history

import (
	"os"
	"testing"
	"time"
)

func TestHistoryManager(t *testing.T) {
	// Create a temporary file for testing
	tempFile := "./test_history"
	defer os.RemoveAll("./data/test_history.connections.json")
	defer os.RemoveAll("./data/test_history.datausage.json")
	defer os.RemoveAll("./data/test_history.alerts.json")

	// Create history manager
	hm := NewHistoryManager(tempFile)

	// Test adding and retrieving connection records
	t.Run("ConnectionRecords", func(t *testing.T) {
		record := ConnectionRecord{
			ID:             "test-connection-1",
			ServerID:       "server-1",
			ServerName:     "Test Server",
			StartTime:      time.Now(),
			EndTime:        time.Now().Add(10 * time.Minute),
			Duration:       600,
			DataSent:       1024,
			DataReceived:   2048,
			Protocol:       "VMess",
			Status:         "connected",
			DisconnectReason: "",
		}

		// Add record
		err := hm.AddConnectionRecord(record)
		if err != nil {
			t.Fatalf("Failed to add connection record: %v", err)
		}

		// Retrieve records
		records, err := hm.GetConnectionRecords(10, 0)
		if err != nil {
			t.Fatalf("Failed to get connection records: %v", err)
		}

		if len(records) != 1 {
			t.Errorf("Expected 1 record, got %d", len(records))
		}

		if records[0].ID != record.ID {
			t.Errorf("Expected ID %s, got %s", record.ID, records[0].ID)
		}
	})

	// Test adding and retrieving data usage records
	t.Run("DataUsageRecords", func(t *testing.T) {
		record := DataUsageRecord{
			ID:           "test-data-1",
			Timestamp:    time.Now(),
			ServerID:     "server-1",
			ServerName:   "Test Server",
			DataSent:     1024,
			DataReceived: 2048,
			TotalSent:    10240,
			TotalReceived: 20480,
		}

		// Add record
		err := hm.AddDataUsageRecord(record)
		if err != nil {
			t.Fatalf("Failed to add data usage record: %v", err)
		}

		// Retrieve records
		records, err := hm.GetDataUsageRecords("", 10, 0)
		if err != nil {
			t.Fatalf("Failed to get data usage records: %v", err)
		}

		if len(records) != 1 {
			t.Errorf("Expected 1 record, got %d", len(records))
		}

		if records[0].ID != record.ID {
			t.Errorf("Expected ID %s, got %s", record.ID, records[0].ID)
		}
	})

	// Test adding and retrieving alert records
	t.Run("AlertRecords", func(t *testing.T) {
		record := AlertRecord{
			ID:        "test-alert-1",
			Timestamp: time.Now(),
			Type:      "warning",
			Title:     "Test Alert",
			Message:   "This is a test alert",
			ServerID:  "server-1",
			Read:      false,
		}

		// Add record
		err := hm.AddAlertRecord(record)
		if err != nil {
			t.Fatalf("Failed to add alert record: %v", err)
		}

		// Retrieve records
		records, err := hm.GetAlertRecords(false, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get alert records: %v", err)
		}

		if len(records) != 1 {
			t.Errorf("Expected 1 record, got %d", len(records))
		}

		if records[0].ID != record.ID {
			t.Errorf("Expected ID %s, got %s", record.ID, records[0].ID)
		}

		// Test unread records
		unreadRecords, err := hm.GetAlertRecords(true, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get unread alert records: %v", err)
		}

		if len(unreadRecords) != 1 {
			t.Errorf("Expected 1 unread record, got %d", len(unreadRecords))
		}

		// Mark as read
		err = hm.MarkAlertAsRead(record.ID)
		if err != nil {
			t.Fatalf("Failed to mark alert as read: %v", err)
		}

		// Check unread records again
		unreadRecords, err = hm.GetAlertRecords(true, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get unread alert records: %v", err)
		}

		if len(unreadRecords) != 0 {
			t.Errorf("Expected 0 unread records, got %d", len(unreadRecords))
		}
	})
}