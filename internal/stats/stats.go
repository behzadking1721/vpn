package stats

import (
	"sync"
	"time"
)

// ConnectionStat represents statistics for a single connection
type ConnectionStat struct {
	Timestamp  time.Time `json:"timestamp"`
	DataSent   int64     `json:"data_sent"`
	DataRecv   int64     `json:"data_recv"`
	ServerID   string    `json:"server_id"`
	ServerName string    `json:"server_name"`
}

// SessionStat represents statistics for a connection session
type SessionStat struct {
	StartedAt  time.Time `json:"started_at"`
	EndedAt    time.Time `json:"ended_at"`
	DataSent   int64     `json:"data_sent"`
	DataRecv   int64     `json:"data_recv"`
	ServerID   string    `json:"server_id"`
	ServerName string    `json:"server_name"`
}

// StatsManager manages connection statistics
type StatsManager struct {
	currentConnection *ConnectionStat
	sessions          []SessionStat
	mutex             sync.RWMutex
}

// NewStatsManager creates a new stats manager
func NewStatsManager() *StatsManager {
	return &StatsManager{
		sessions: make([]SessionStat, 0),
	}
}

// StartConnection starts tracking a new connection
func (sm *StatsManager) StartConnection(serverID, serverName string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// If there's an ongoing connection, end it first
	if sm.currentConnection != nil {
		sm.endCurrentConnection()
	}

	// Start new connection
	sm.currentConnection = &ConnectionStat{
		Timestamp:  time.Now(),
		DataSent:   0,
		DataRecv:   0,
		ServerID:   serverID,
		ServerName: serverName,
	}
}

// UpdateConnection updates the current connection statistics
func (sm *StatsManager) UpdateConnection(dataSent, dataRecv int64) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if sm.currentConnection != nil {
		sm.currentConnection.DataSent += dataSent
		sm.currentConnection.DataRecv += dataRecv
	}
}

// EndConnection ends the current connection and saves it as a session
func (sm *StatsManager) EndConnection() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.endCurrentConnection()
}

// endCurrentConnection ends the current connection without locking
func (sm *StatsManager) endCurrentConnection() {
	if sm.currentConnection != nil {
		session := SessionStat{
			StartedAt:  sm.currentConnection.Timestamp,
			EndedAt:    time.Now(),
			DataSent:   sm.currentConnection.DataSent,
			DataRecv:   sm.currentConnection.DataRecv,
			ServerID:   sm.currentConnection.ServerID,
			ServerName: sm.currentConnection.ServerName,
		}
		
		sm.sessions = append(sm.sessions, session)
		sm.currentConnection = nil
	}
}

// GetCurrentConnection returns the current connection statistics
func (sm *StatsManager) GetCurrentConnection() *ConnectionStat {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if sm.currentConnection == nil {
		return nil
	}
	
	// Return a copy to prevent external modification
	stat := *sm.currentConnection
	return &stat
}

// GetSessions returns all session statistics
func (sm *StatsManager) GetSessions() []SessionStat {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Return a copy to prevent external modification
	sessions := make([]SessionStat, len(sm.sessions))
	copy(sessions, sm.sessions)
	return sessions
}

// GetSessionsByTimeRange returns sessions within a specific time range
func (sm *StatsManager) GetSessionsByTimeRange(start, end time.Time) []SessionStat {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var result []SessionStat
	for _, session := range sm.sessions {
		if session.StartedAt.After(start) && session.StartedAt.Before(end) {
			result = append(result, session)
		}
	}
	
	return result
}

// GetTotalDataUsage returns total data usage across all sessions
func (sm *StatsManager) GetTotalDataUsage() (int64, int64) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var totalSent, totalRecv int64
	
	// Add current connection if exists
	if sm.currentConnection != nil {
		totalSent += sm.currentConnection.DataSent
		totalRecv += sm.currentConnection.DataRecv
	}
	
	// Add all sessions
	for _, session := range sm.sessions {
		totalSent += session.DataSent
		totalRecv += session.DataRecv
	}
	
	return totalSent, totalRecv
}

// GetDailyDataUsage returns daily data usage for the last n days
func (sm *StatsManager) GetDailyDataUsage(days int) []ConnectionStat {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	dailyStats := make([]ConnectionStat, 0)
	
	// Create a map to store daily data
	dailyData := make(map[string]*ConnectionStat)
	
	// Process current connection
	if sm.currentConnection != nil {
		dateKey := sm.currentConnection.Timestamp.Format("2006-01-02")
		if dailyData[dateKey] == nil {
			dailyData[dateKey] = &ConnectionStat{
				Timestamp: sm.currentConnection.Timestamp,
				DataSent:  0,
				DataRecv:  0,
			}
		}
		dailyData[dateKey].DataSent += sm.currentConnection.DataSent
		dailyData[dateKey].DataRecv += sm.currentConnection.DataRecv
	}
	
	// Process all sessions
	for _, session := range sm.sessions {
		dateKey := session.StartedAt.Format("2006-01-02")
		if dailyData[dateKey] == nil {
			dailyData[dateKey] = &ConnectionStat{
				Timestamp: time.Date(session.StartedAt.Year(), session.StartedAt.Month(), session.StartedAt.Day(), 0, 0, 0, 0, time.UTC),
				DataSent:  0,
				DataRecv:  0,
			}
		}
		dailyData[dateKey].DataSent += session.DataSent
		dailyData[dateKey].DataRecv += session.DataRecv
	}
	
	// Convert map to slice
	for _, stat := range dailyData {
		dailyStats = append(dailyStats, *stat)
	}
	
	return dailyStats
}

// ClearStats clears all statistics
func (sm *StatsManager) ClearStats() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.currentConnection = nil
	sm.sessions = make([]SessionStat, 0)
}