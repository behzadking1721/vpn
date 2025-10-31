package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// HistoryManager manages the history of connections and data usage
type HistoryManager struct {
	historyFile string
	mutex       sync.RWMutex
}

// ConnectionRecord represents a connection record
type ConnectionRecord struct {
	ID               string    `json:"id"`
	ServerID         string    `json:"server_id"`
	ServerName       string    `json:"server_name"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Duration         int64     `json:"duration"` // in seconds
	DataSent         int64     `json:"data_sent"`
	DataReceived     int64     `json:"data_received"`
	Protocol         string    `json:"protocol"`
	Status           string    `json:"status"` // connected, disconnected, error
	DisconnectReason string    `json:"disconnect_reason,omitempty"`
}

// DataUsageRecord represents a data usage record
type DataUsageRecord struct {
	ID            string    `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
	ServerID      string    `json:"server_id"`
	ServerName    string    `json:"server_name"`
	DataSent      int64     `json:"data_sent"`
	DataReceived  int64     `json:"data_received"`
	TotalSent     int64     `json:"total_sent"`
	TotalReceived int64     `json:"total_received"`
}

// AlertRecord represents an alert record
type AlertRecord struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // warning, error, info
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	ServerID  string    `json:"server_id,omitempty"`
	Read      bool      `json:"read"`
}

// NewHistoryManager creates a new history manager
func NewHistoryManager(historyFile string) *HistoryManager {
	// Create directory if it doesn't exist
	dir := filepath.Dir(historyFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create history directory: %v\n", err)
	}

	return &HistoryManager{
		historyFile: historyFile,
	}
}

// AddConnectionRecord adds a connection record to history
func (h *HistoryManager) AddConnectionRecord(record ConnectionRecord) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Load existing records
	records, err := h.loadConnectionRecords()
	if err != nil {
		return err
	}

	// Add new record
	records = append(records, record)

	// Save records
	return h.saveConnectionRecords(records)
}

// GetConnectionRecords retrieves connection records
func (h *HistoryManager) GetConnectionRecords(limit int, offset int) ([]ConnectionRecord, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	records, err := h.loadConnectionRecords()
	if err != nil {
		return nil, err
	}

	// Apply limit and offset
	if limit > 0 {
		start := offset
		end := offset + limit
		if start >= len(records) {
			return []ConnectionRecord{}, nil
		}
		if end > len(records) {
			end = len(records)
		}
		records = records[start:end]
	}

	return records, nil
}

// AddDataUsageRecord adds a data usage record to history
func (h *HistoryManager) AddDataUsageRecord(record DataUsageRecord) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Load existing records
	records, err := h.loadDataUsageRecords()
	if err != nil {
		return err
	}

	// Add new record
	records = append(records, record)

	// Save records
	return h.saveDataUsageRecords(records)
}

// GetDataUsageRecords retrieves data usage records
func (h *HistoryManager) GetDataUsageRecords(serverID string, limit int, offset int) ([]DataUsageRecord, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	records, err := h.loadDataUsageRecords()
	if err != nil {
		return nil, err
	}

	// Filter by server ID if provided
	if serverID != "" {
		filtered := make([]DataUsageRecord, 0)
		for _, record := range records {
			if record.ServerID == serverID {
				filtered = append(filtered, record)
			}
		}
		records = filtered
	}

	// Apply limit and offset
	if limit > 0 {
		start := offset
		end := offset + limit
		if start >= len(records) {
			return []DataUsageRecord{}, nil
		}
		if end > len(records) {
			end = len(records)
		}
		records = records[start:end]
	}

	return records, nil
}

// AddAlertRecord adds an alert record to history
func (h *HistoryManager) AddAlertRecord(record AlertRecord) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Load existing records
	records, err := h.loadAlertRecords()
	if err != nil {
		return err
	}

	// Add new record
	records = append(records, record)

	// Save records
	return h.saveAlertRecords(records)
}

// GetAlertRecords retrieves alert records
func (h *HistoryManager) GetAlertRecords(unreadOnly bool, limit int, offset int) ([]AlertRecord, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	records, err := h.loadAlertRecords()
	if err != nil {
		return nil, err
	}

	// Filter by unread status if requested
	if unreadOnly {
		filtered := make([]AlertRecord, 0)
		for _, record := range records {
			if !record.Read {
				filtered = append(filtered, record)
			}
		}
		records = filtered
	}

	// Apply limit and offset
	if limit > 0 {
		start := offset
		end := offset + limit
		if start >= len(records) {
			return []AlertRecord{}, nil
		}
		if end > len(records) {
			end = len(records)
		}
		records = records[start:end]
	}

	return records, nil
}

// MarkAlertAsRead marks an alert as read
func (h *HistoryManager) MarkAlertAsRead(alertID string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	records, err := h.loadAlertRecords()
	if err != nil {
		return err
	}

	// Find and mark the alert as read
	for i := range records {
		if records[i].ID == alertID {
			records[i].Read = true
			break
		}
	}

	// Save records
	return h.saveAlertRecords(records)
}

// loadConnectionRecords loads connection records from file
func (h *HistoryManager) loadConnectionRecords() ([]ConnectionRecord, error) {
	filePath := h.historyFile + ".connections.json"

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []ConnectionRecord{}, nil
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var records []ConnectionRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	return records, nil
}

// saveConnectionRecords saves connection records to file
func (h *HistoryManager) saveConnectionRecords(records []ConnectionRecord) error {
	filePath := h.historyFile + ".connections.json"

	// Convert to JSON
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filePath, data, 0644)
}

// loadDataUsageRecords loads data usage records from file
func (h *HistoryManager) loadDataUsageRecords() ([]DataUsageRecord, error) {
	filePath := h.historyFile + ".datausage.json"

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []DataUsageRecord{}, nil
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var records []DataUsageRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	return records, nil
}

// saveDataUsageRecords saves data usage records to file
func (h *HistoryManager) saveDataUsageRecords(records []DataUsageRecord) error {
	filePath := h.historyFile + ".datausage.json"

	// Convert to JSON
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filePath, data, 0644)
}

// loadAlertRecords loads alert records from file
func (h *HistoryManager) loadAlertRecords() ([]AlertRecord, error) {
	filePath := h.historyFile + ".alerts.json"

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []AlertRecord{}, nil
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var records []AlertRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	return records, nil
}

// saveAlertRecords saves alert records to file
func (h *HistoryManager) saveAlertRecords(records []AlertRecord) error {
	filePath := h.historyFile + ".alerts.json"

	// Convert to JSON
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filePath, data, 0644)
}
