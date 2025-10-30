package database

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

// Manager handles database operations
type Manager struct {
	db *DB
}

// NewManager creates a new database manager
func NewManager(dbPath string) (*Manager, error) {
	// Open database connection
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create DB wrapper
	db := &DB{sqlDB}

	// Create manager
	manager := &Manager{
		db: db,
	}

	// Run migrations
	if err := manager.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return manager, nil
}

// Close closes the database connection
func (m *Manager) Close() error {
	return m.db.Close()
}

// runMigrations runs database migrations
func (m *Manager) runMigrations() error {
	// Create migrations table if it doesn't exist
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Check which migrations have been applied
	appliedMigrations := make(map[string]bool)
	rows, err := m.db.Query("SELECT name FROM migrations")
	if err != nil {
		return fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan migration name: %w", err)
		}
		appliedMigrations[name] = true
	}

	// Define migrations
	migrations := []struct {
		Name string
		SQL  string
	}{
		{
			Name: "create_connections_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS connections (
					id TEXT PRIMARY KEY,
					server_id TEXT NOT NULL,
					server_name TEXT NOT NULL,
					start_time TIMESTAMP NOT NULL,
					end_time TIMESTAMP,
					duration INTEGER,
					data_sent INTEGER,
					data_received INTEGER,
					protocol TEXT,
					status TEXT,
					disconnect_reason TEXT
				)`,
		},
		{
			Name: "create_data_usage_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS data_usage (
					id TEXT PRIMARY KEY,
					timestamp TIMESTAMP NOT NULL,
					server_id TEXT NOT NULL,
					server_name TEXT NOT NULL,
					data_sent INTEGER,
					data_received INTEGER,
					total_sent INTEGER,
					total_received INTEGER
				)`,
		},
		{
			Name: "create_pings_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS pings (
					id TEXT PRIMARY KEY,
					timestamp TIMESTAMP NOT NULL,
					server_id TEXT NOT NULL,
					server_name TEXT NOT NULL,
					ping INTEGER
				)`,
		},
		{
			Name: "create_alerts_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS alerts (
					id TEXT PRIMARY KEY,
					type TEXT NOT NULL,
					title TEXT NOT NULL,
					message TEXT NOT NULL,
					timestamp TIMESTAMP NOT NULL,
					value REAL,
					severity TEXT,
					resolved BOOLEAN DEFAULT FALSE,
					read BOOLEAN DEFAULT FALSE,
					server_id TEXT,
					server_name TEXT
				)`,
		},
		{
			Name: "create_dashboard_settings_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS dashboard_settings (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					theme TEXT NOT NULL DEFAULT 'system',
					chart_window TEXT NOT NULL DEFAULT '24h',
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				)`,
		},
		{
			Name: "create_connections_index",
			SQL: `CREATE INDEX IF NOT EXISTS idx_connections_server_id ON connections(server_id)`,
		},
		{
			Name: "create_connections_time_index",
			SQL: `CREATE INDEX IF NOT EXISTS idx_connections_time ON connections(start_time)`,
		},
		{
			Name: "create_data_usage_index",
			SQL: `CREATE INDEX IF NOT EXISTS idx_data_usage_server_id ON data_usage(server_id)`,
		},
		{
			Name: "create_data_usage_time_index",
			SQL: `CREATE INDEX IF NOT EXISTS idx_data_usage_time ON data_usage(timestamp)`,
		},
		{
			Name: "create_pings_index",
			SQL: `CREATE INDEX IF NOT EXISTS idx_pings_server_id ON pings(server_id)`,
		},
		{
			Name: "create_pings_time_index",
			SQL: `CREATE INDEX IF NOT EXISTS idx_pings_time ON pings(timestamp)`,
		},
		{
			Name: "create_alerts_index",
			SQL: `CREATE INDEX IF NOT EXISTS idx_alerts_time ON alerts(timestamp)`,
		},
		{
			Name: "create_alerts_resolved_index",
			SQL: `CREATE INDEX IF NOT EXISTS idx_alerts_resolved ON alerts(resolved)`,
		},
		{
			Name: "create_alerts_read_index",
			SQL: `CREATE INDEX IF NOT EXISTS idx_alerts_read ON alerts(read)`,
		},
	}

	// Apply migrations
	for _, migration := range migrations {
		// Skip if already applied
		if appliedMigrations[migration.Name] {
			continue
		}

		// Apply migration
		_, err := m.db.Exec(migration.SQL)
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Name, err)
		}

		// Record that migration was applied
		_, err = m.db.Exec("INSERT INTO migrations (name) VALUES (?)", migration.Name)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
		}
	}

	return nil
}

// AddConnectionRecord adds a connection record to the database
func (m *Manager) AddConnectionRecord(record ConnectionRecord) error {
	_, err := m.db.Exec(`
		INSERT INTO connections (
			id, server_id, server_name, start_time, end_time, duration,
			data_sent, data_received, protocol, status, disconnect_reason
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, record.ID, record.ServerID, record.ServerName, record.StartTime, record.EndTime,
		record.Duration, record.DataSent, record.DataReceived, record.Protocol,
		record.Status, record.DisconnectReason)

	return err
}

// GetConnectionRecords retrieves connection records from the database
func (m *Manager) GetConnectionRecords(limit, offset int) ([]ConnectionRecord, error) {
	query := `
		SELECT id, server_id, server_name, start_time, end_time, duration,
		       data_sent, data_received, protocol, status, disconnect_reason
		FROM connections
		ORDER BY start_time DESC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
		if offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", offset)
		}
	}

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ConnectionRecord
	for rows.Next() {
		var record ConnectionRecord
		err := rows.Scan(
			&record.ID, &record.ServerID, &record.ServerName, &record.StartTime,
			&record.EndTime, &record.Duration, &record.DataSent, &record.DataReceived,
			&record.Protocol, &record.Status, &record.DisconnectReason,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// AddDataUsageRecord adds a data usage record to the database
func (m *Manager) AddDataUsageRecord(record DataUsageRecord) error {
	_, err := m.db.Exec(`
		INSERT INTO data_usage (
			id, timestamp, server_id, server_name, data_sent, data_received,
			total_sent, total_received
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, record.ID, record.Timestamp, record.ServerID, record.ServerName,
		record.DataSent, record.DataReceived, record.TotalSent, record.TotalReceived)

	return err
}

// GetDataUsageRecords retrieves data usage records from the database
func (m *Manager) GetDataUsageRecords(serverID string, limit, offset int) ([]DataUsageRecord, error) {
	query := `
		SELECT id, timestamp, server_id, server_name, data_sent, data_received,
		       total_sent, total_received
		FROM data_usage
	`

	args := []interface{}{}
	if serverID != "" {
		query += " WHERE server_id = ?"
		args = append(args, serverID)
	}

	query += " ORDER BY timestamp DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
		if offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", offset)
		}
	}

	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []DataUsageRecord
	for rows.Next() {
		var record DataUsageRecord
		err := rows.Scan(
			&record.ID, &record.Timestamp, &record.ServerID, &record.ServerName,
			&record.DataSent, &record.DataReceived, &record.TotalSent, &record.TotalReceived,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// AddPingRecord adds a ping record to the database
func (m *Manager) AddPingRecord(record PingRecord) error {
	_, err := m.db.Exec(`
		INSERT INTO pings (id, timestamp, server_id, server_name, ping)
		VALUES (?, ?, ?, ?, ?)
	`, record.ID, record.Timestamp, record.ServerID, record.ServerName, record.Ping)

	return err
}

// GetPingRecords retrieves ping records from the database
func (m *Manager) GetPingRecords(serverID string, limit int) ([]PingRecord, error) {
	query := `
		SELECT id, timestamp, server_id, server_name, ping
		FROM pings
	`

	args := []interface{}{}
	if serverID != "" {
		query += " WHERE server_id = ?"
		args = append(args, serverID)
	}

	query += " ORDER BY timestamp DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []PingRecord
	for rows.Next() {
		var record PingRecord
		err := rows.Scan(&record.ID, &record.Timestamp, &record.ServerID, &record.ServerName, &record.Ping)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// AddAlertRecord adds an alert record to the database
func (m *Manager) AddAlertRecord(record AlertRecord) error {
	_, err := m.db.Exec(`
		INSERT INTO alerts (
			id, type, title, message, timestamp, value, severity, resolved, read, server_id, server_name
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, record.ID, record.Type, record.Title, record.Message, record.Timestamp,
		record.Value, record.Severity, record.Resolved, record.Read, record.ServerID, record.ServerName)

	return err
}

// GetAlertRecords retrieves alert records from the database
func (m *Manager) GetAlertRecords(unread, unresolved bool, limit int) ([]AlertRecord, error) {
	query := `
		SELECT id, type, title, message, timestamp, value, severity, resolved, read, server_id, server_name
		FROM alerts
		WHERE 1=1
	`

	args := []interface{}{}
	if unread {
		query += " AND read = ?"
		args = append(args, false)
	}

	if unresolved {
		query += " AND resolved = ?"
		args = append(args, false)
	}

	query += " ORDER BY timestamp DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []AlertRecord
	for rows.Next() {
		var record AlertRecord
		err := rows.Scan(
			&record.ID, &record.Type, &record.Title, &record.Message, &record.Timestamp,
			&record.Value, &record.Severity, &record.Resolved, &record.Read, &record.ServerID, &record.ServerName,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// UpdateAlertRecord updates an alert record in the database
func (m *Manager) UpdateAlertRecord(record AlertRecord) error {
	_, err := m.db.Exec(`
		UPDATE alerts
		SET type = ?, title = ?, message = ?, timestamp = ?, value = ?, severity = ?, resolved = ?, read = ?, server_id = ?, server_name = ?
		WHERE id = ?
	`, record.Type, record.Title, record.Message, record.Timestamp,
		record.Value, record.Severity, record.Resolved, record.Read, record.ServerID, record.ServerName, record.ID)

	return err
}

// GetDashboardSettings retrieves dashboard settings from the database
func (m *Manager) GetDashboardSettings() (*DashboardSettingsRecord, error) {
	var settings DashboardSettingsRecord
	err := m.db.QueryRow(`
		SELECT id, theme, chart_window, created_at, updated_at
		FROM dashboard_settings
		ORDER BY id DESC
		LIMIT 1
	`).Scan(&settings.ID, &settings.Theme, &settings.ChartWindow, &settings.CreatedAt, &settings.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// Return default settings if none exist
			return &DashboardSettingsRecord{
				Theme:       "system",
				ChartWindow: "24h",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, nil
		}
		return nil, err
	}

	return &settings, nil
}

// UpdateDashboardSettings updates dashboard settings in the database
func (m *Manager) UpdateDashboardSettings(settings DashboardSettingsRecord) error {
	// Check if settings already exist
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM dashboard_settings").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// Insert new settings
		_, err = m.db.Exec(`
			INSERT INTO dashboard_settings (theme, chart_window, created_at, updated_at)
			VALUES (?, ?, ?, ?)
		`, settings.Theme, settings.ChartWindow, settings.CreatedAt, settings.UpdatedAt)
	} else {
		// Update existing settings
		_, err = m.db.Exec(`
			UPDATE dashboard_settings
			SET theme = ?, chart_window = ?, updated_at = ?
			WHERE id = ?
		`, settings.Theme, settings.ChartWindow, settings.UpdatedAt, settings.ID)
	}

	return err
}