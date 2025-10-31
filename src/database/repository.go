package database

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/alert"
	"c:/Users/behza/OneDrive/Documents/vpn/src/history"
	"c:/Users/behza/OneDrive/Documents/vpn/src/settings"
)

// HistoryRepository implements history.HistoryManager interface using SQLite
type HistoryRepository struct {
	dbManager *Manager
}

// NewHistoryRepository creates a new history repository
func NewHistoryRepository(dbManager *Manager) *HistoryRepository {
	return &HistoryRepository{
		dbManager: dbManager,
	}
}

// AddConnectionRecord adds a connection record
func (r *HistoryRepository) AddConnectionRecord(record history.ConnectionRecord) error {
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

	return r.dbManager.AddConnectionRecord(dbRecord)
}

// GetConnectionRecords retrieves connection records
func (r *HistoryRepository) GetConnectionRecords(limit, offset int) ([]history.ConnectionRecord, error) {
	dbRecords, err := r.dbManager.GetConnectionRecords(limit, offset)
	if err != nil {
		return nil, err
	}

	records := make([]history.ConnectionRecord, len(dbRecords))
	for i, dbRecord := range dbRecords {
		records[i] = history.ConnectionRecord{
			ID:               dbRecord.ID,
			ServerID:         dbRecord.ServerID,
			ServerName:       dbRecord.ServerName,
			StartTime:        dbRecord.StartTime,
			EndTime:          dbRecord.EndTime,
			Duration:         dbRecord.Duration,
			DataSent:         dbRecord.DataSent,
			DataReceived:     dbRecord.DataReceived,
			Protocol:         dbRecord.Protocol,
			Status:           dbRecord.Status,
			DisconnectReason: dbRecord.DisconnectReason,
		}
	}

	return records, nil
}

// AddDataUsageRecord adds a data usage record
func (r *HistoryRepository) AddDataUsageRecord(record history.DataUsageRecord) error {
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

	return r.dbManager.AddDataUsageRecord(dbRecord)
}

// GetDataUsageRecords retrieves data usage records
func (r *HistoryRepository) GetDataUsageRecords(serverID string, limit, offset int) ([]history.DataUsageRecord, error) {
	dbRecords, err := r.dbManager.GetDataUsageRecords(serverID, limit, offset)
	if err != nil {
		return nil, err
	}

	records := make([]history.DataUsageRecord, len(dbRecords))
	for i, dbRecord := range dbRecords {
		records[i] = history.DataUsageRecord{
			ID:            dbRecord.ID,
			Timestamp:     dbRecord.Timestamp,
			ServerID:      dbRecord.ServerID,
			ServerName:    dbRecord.ServerName,
			DataSent:      dbRecord.DataSent,
			DataReceived:  dbRecord.DataReceived,
			TotalSent:     dbRecord.TotalSent,
			TotalReceived: dbRecord.TotalReceived,
		}
	}

	return records, nil
}

// AlertRepository implements alert.AlertManager interface using SQLite
type AlertRepository struct {
	dbManager *Manager
}

// NewAlertRepository creates a new alert repository
func NewAlertRepository(dbManager *Manager) *AlertRepository {
	return &AlertRepository{
		dbManager: dbManager,
	}
}

// AddAlertRecord adds an alert record
func (r *AlertRepository) AddAlertRecord(record alert.Alert) error {
	dbRecord := AlertRecord{
		ID:         record.ID,
		Type:       string(record.Type),
		Title:      record.Title,
		Message:    record.Message,
		Timestamp:  record.Timestamp,
		Value:      record.Value,
		Severity:   string(record.Severity),
		Resolved:   record.Resolved,
		Read:       record.Read,
		ServerID:   record.ServerID,
		ServerName: record.ServerName,
	}

	return r.dbManager.AddAlertRecord(dbRecord)
}

// GetAlertRecords retrieves alert records
func (r *AlertRepository) GetAlertRecords(unread, unresolved bool, limit int) ([]alert.Alert, error) {
	dbRecords, err := r.dbManager.GetAlertRecords(unread, unresolved, limit)
	if err != nil {
		return nil, err
	}

	records := make([]alert.Alert, len(dbRecords))
	for i, dbRecord := range dbRecords {
		records[i] = alert.Alert{
			ID:         dbRecord.ID,
			Type:       alert.AlertRuleType(dbRecord.Type),
			Title:      dbRecord.Title,
			Message:    dbRecord.Message,
			Timestamp:  dbRecord.Timestamp,
			Value:      dbRecord.Value,
			Severity:   alert.AlertSeverity(dbRecord.Severity),
			Resolved:   dbRecord.Resolved,
			Read:       dbRecord.Read,
			ServerID:   dbRecord.ServerID,
			ServerName: dbRecord.ServerName,
		}
	}

	return records, nil
}

// UpdateAlertRecord updates an alert record
func (r *AlertRepository) UpdateAlertRecord(record alert.Alert) error {
	dbRecord := AlertRecord{
		ID:         record.ID,
		Type:       string(record.Type),
		Title:      record.Title,
		Message:    record.Message,
		Timestamp:  record.Timestamp,
		Value:      record.Value,
		Severity:   string(record.Severity),
		Resolved:   record.Resolved,
		Read:       record.Read,
		ServerID:   record.ServerID,
		ServerName: record.ServerName,
	}

	return r.dbManager.UpdateAlertRecord(dbRecord)
}

// SettingsRepository implements settings.SettingsManager interface using SQLite
type SettingsRepository struct {
	dbManager *Manager
}

// NewSettingsRepository creates a new settings repository
func NewSettingsRepository(dbManager *Manager) *SettingsRepository {
	return &SettingsRepository{
		dbManager: dbManager,
	}
}

// GetSettings retrieves dashboard settings
func (r *SettingsRepository) GetSettings() (*settings.DashboardSettings, error) {
	dbSettings, err := r.dbManager.GetDashboardSettings()
	if err != nil {
		return nil, err
	}

	settings := &settings.DashboardSettings{
		Theme:       settings.DashboardTheme(dbSettings.Theme),
		ChartWindow: settings.ChartWindow(dbSettings.ChartWindow),
		CreatedAt:   dbSettings.CreatedAt,
		UpdatedAt:   dbSettings.UpdatedAt,
	}

	return settings, nil
}

// UpdateSettings updates dashboard settings
func (r *SettingsRepository) UpdateSettings(newSettings settings.DashboardSettings) error {
	dbSettings := DashboardSettingsRecord{
		Theme:       string(newSettings.Theme),
		ChartWindow: string(newSettings.ChartWindow),
		CreatedAt:   newSettings.CreatedAt,
		UpdatedAt:   newSettings.UpdatedAt,
	}

	return r.dbManager.UpdateDashboardSettings(dbSettings)
}
