package api

import (
	"encoding/json"
	"net/http"
	"time"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"c:/Users/behza/OneDrive/Documents/vpn/src/history"
	"github.com/gorilla/mux"
)

// DashboardAPI handles dashboard-related API endpoints
type DashboardAPI struct {
	serverManager *managers.ServerManager
	connManager   *managers.ConnectionManager
	configManager *managers.ConfigManager
	dataManager   *managers.DataManager
	historyManager *history.HistoryManager
}

// NewDashboardAPI creates a new dashboard API instance
func NewDashboardAPI(
	serverMgr *managers.ServerManager,
	connMgr *managers.ConnectionManager,
	configMgr *managers.ConfigManager,
	dataMgr *managers.DataManager,
	historyMgr *history.HistoryManager) *DashboardAPI {
	
	return &DashboardAPI{
		serverManager:  serverMgr,
		connManager:    connMgr,
		configManager:  configMgr,
		dataManager:    dataMgr,
		historyManager: historyMgr,
	}
}

// RegisterRoutes registers dashboard routes with the router
func (d *DashboardAPI) RegisterRoutes(router *mux.Router) {
	// Dashboard endpoints
	router.HandleFunc("/api/dashboard/status", d.getConnectionStatus).Methods("GET")
	router.HandleFunc("/api/dashboard/usage", d.getDataUsage).Methods("GET")
	router.HandleFunc("/api/dashboard/servers", d.getServersStatus).Methods("GET")
	router.HandleFunc("/api/dashboard/logs", d.getRecentLogs).Methods("GET")
	router.HandleFunc("/api/dashboard/stats", d.getStatistics).Methods("GET")
	
	// History endpoints
	router.HandleFunc("/api/dashboard/history/connections", d.getConnectionHistory).Methods("GET")
	router.HandleFunc("/api/dashboard/history/datausage", d.getDataUsageHistory).Methods("GET")
	router.HandleFunc("/api/dashboard/history/alerts", d.getAlertHistory).Methods("GET")
	router.HandleFunc("/api/dashboard/history/alerts/{id}/read", d.markAlertAsRead).Methods("POST")
}

// ConnectionStatusResponse represents the connection status response
type ConnectionStatusResponse struct {
	Connected     bool              `json:"connected"`
	Server        *core.Server      `json:"server,omitempty"`
	ConnectionTime time.Time        `json:"connection_time,omitempty"`
	DataSent      int64             `json:"data_sent"`
	DataReceived  int64             `json:"data_received"`
}

// getConnectionStatus returns the current connection status
func (d *DashboardAPI) getConnectionStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	status := d.connManager.GetStatus()
	connInfo := d.connManager.GetConnectionInfo()
	
	response := ConnectionStatusResponse{
		Connected:    status == core.StatusConnected,
		DataSent:     connInfo.DataSent,
		DataReceived: connInfo.DataReceived,
	}
	
	if status == core.StatusConnected {
		currentServer := d.connManager.GetCurrentServer()
		if currentServer != nil {
			response.Server = currentServer
			response.ConnectionTime = connInfo.StartTime
		}
	}
	
	json.NewEncoder(w).Encode(response)
}

// DataUsageResponse represents the data usage response
type DataUsageResponse struct {
	CurrentSession struct {
		DataSent     int64 `json:"data_sent"`
		DataReceived int64 `json:"data_received"`
	} `json:"current_session"`
	
	TotalUsage struct {
		DataSent     int64 `json:"data_sent"`
		DataReceived int64 `json:"data_received"`
	} `json:"total_usage"`
	
	DailyUsage []struct {
		Date         string `json:"date"`
		DataSent     int64  `json:"data_sent"`
		DataReceived int64  `json:"data_received"`
	} `json:"daily_usage"`
	
	Limits []struct {
		ServerID   string `json:"server_id"`
		ServerName string `json:"server_name"`
		Limit      int64  `json:"limit"`
		Used       int64  `json:"used"`
		Remaining  int64  `json:"remaining"`
	} `json:"limits"`
}

// getDataUsage returns data usage statistics
func (d *DashboardAPI) getDataUsage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := DataUsageResponse{}
	
	// Current session data
	connInfo := d.connManager.GetConnectionInfo()
	response.CurrentSession.DataSent = connInfo.DataSent
	response.CurrentSession.DataReceived = connInfo.DataReceived
	
	// Total usage data
	allData := d.dataManager.GetAllData()
	var totalSent, totalReceived int64
	for _, data := range allData {
		totalSent += data.TotalSent
		totalReceived += data.TotalRecv
	}
	response.TotalUsage.DataSent = totalSent
	response.TotalUsage.DataReceived = totalReceived
	
	// Daily usage (simplified - in a real implementation, you would store daily data)
	response.DailyUsage = make([]struct {
		Date         string `json:"date"`
		DataSent     int64  `json:"data_sent"`
		DataReceived int64  `json:"data_received"`
	}, 0)
	
	// Limits
	servers := d.serverManager.GetAllServers()
	response.Limits = make([]struct {
		ServerID   string `json:"server_id"`
		ServerName string `json:"server_name"`
		Limit      int64  `json:"limit"`
		Used       int64  `json:"used"`
		Remaining  int64  `json:"remaining"`
	}, 0)
	
	for _, server := range servers {
		if server.DataLimit > 0 {
			limit := struct {
				ServerID   string `json:"server_id"`
				ServerName string `json:"server_name"`
				Limit      int64  `json:"limit"`
				Used       int64  `json:"used"`
				Remaining  int64  `json:"remaining"`
			}{
				ServerID:   server.ID,
				ServerName: server.Name,
				Limit:      server.DataLimit,
				Used:       server.DataUsed,
				Remaining:  server.DataLimit - server.DataUsed,
			}
			response.Limits = append(response.Limits, limit)
		}
	}
	
	json.NewEncoder(w).Encode(response)
}

// ServerStatusResponse represents a server status response
type ServerStatusResponse struct {
	Servers []struct {
		ID           string              `json:"id"`
		Name         string              `json:"name"`
		Protocol     core.ProtocolType   `json:"protocol"`
		Host         string              `json:"host"`
		Ping         int                 `json:"ping"`
		LastPing     core.Time           `json:"last_ping"`
		Enabled      bool                `json:"enabled"`
		DataLimit    int64               `json:"data_limit"`
		DataUsed     int64               `json:"data_used"`
		Status       string              `json:"status"` // "online", "offline", "limited"
	} `json:"servers"`
}

// getServersStatus returns the status of all servers
func (d *DashboardAPI) getServersStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	servers := d.serverManager.GetAllServers()
	response := ServerStatusResponse{
		Servers: make([]struct {
			ID           string              `json:"id"`
			Name         string              `json:"name"`
			Protocol     core.ProtocolType   `json:"protocol"`
			Host         string              `json:"host"`
			Ping         int                 `json:"ping"`
			LastPing     core.Time           `json:"last_ping"`
			Enabled      bool                `json:"enabled"`
			DataLimit    int64               `json:"data_limit"`
			DataUsed     int64               `json:"data_used"`
			Status       string              `json:"status"`
		}, len(servers)),
	}
	
	for i, server := range servers {
		status := "online"
		if !server.Enabled {
			status = "offline"
		} else if server.DataLimit > 0 && server.DataUsed >= server.DataLimit {
			status = "limited"
		}
		
		response.Servers[i] = struct {
			ID           string              `json:"id"`
			Name         string              `json:"name"`
			Protocol     core.ProtocolType   `json:"protocol"`
			Host         string              `json:"host"`
			Ping         int                 `json:"ping"`
			LastPing     core.Time           `json:"last_ping"`
			Enabled      bool                `json:"enabled"`
			DataLimit    int64               `json:"data_limit"`
			DataUsed     int64               `json:"data_used"`
			Status       string              `json:"status"`
		}{
			ID:        server.ID,
			Name:      server.Name,
			Protocol:  server.Protocol,
			Host:      server.Host,
			Ping:      server.Ping,
			LastPing:  server.LastPing,
			Enabled:   server.Enabled,
			DataLimit: server.DataLimit,
			DataUsed:  server.DataUsed,
			Status:    status,
		}
	}
	
	json.NewEncoder(w).Encode(response)
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
}

// RecentLogsResponse represents the recent logs response
type RecentLogsResponse struct {
	Logs []LogEntry `json:"logs"`
}

// getRecentLogs returns recent log entries
func (d *DashboardAPI) getRecentLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// In a real implementation, you would read from the log file
	// For now, we'll return a sample response
	response := RecentLogsResponse{
		Logs: []LogEntry{
			{
				Timestamp: time.Now().Add(-1 * time.Minute),
				Level:     "INFO",
				Message:   "Connected to server Japan-1",
				Source:    "ConnectionManager",
			},
			{
				Timestamp: time.Now().Add(-5 * time.Minute),
				Level:     "DEBUG",
				Message:   "Server ping completed: 22ms",
				Source:    "ServerManager",
			},
			{
				Timestamp: time.Now().Add(-10 * time.Minute),
				Level:     "WARN",
				Message:   "High data usage detected on server US-1",
				Source:    "DataManager",
			},
		},
	}
	
	json.NewEncoder(w).Encode(response)
}

// StatisticsResponse represents the statistics response
type StatisticsResponse struct {
	TotalConnections   int     `json:"total_connections"`
	TotalDataSent      int64   `json:"total_data_sent"`
	TotalDataReceived  int64   `json:"total_data_received"`
	AvgConnectionTime  string  `json:"avg_connection_time"`
	ActiveServers      int     `json:"active_servers"`
	TotalServers       int     `json:"total_servers"`
}

// getStatistics returns general statistics
func (d *DashboardAPI) getStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	servers := d.serverManager.GetAllServers()
	var activeServers int
	for _, server := range servers {
		if server.Enabled {
			activeServers++
		}
	}
	
	// In a real implementation, you would track connection history
	// For now, we'll return sample data
	response := StatisticsResponse{
		TotalConnections:   42,
		TotalDataSent:      1024 * 1024 * 1024 * 5, // 5 GB
		TotalDataReceived:  1024 * 1024 * 1024 * 15, // 15 GB
		AvgConnectionTime:  "2h 15m",
		ActiveServers:      activeServers,
		TotalServers:       len(servers),
	}
	
	json.NewEncoder(w).Encode(response)
}

// getConnectionHistory returns connection history
func (d *DashboardAPI) getConnectionHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get limit and offset from query parameters
	limit := 50
	offset := 0
	
	// In a real implementation, you would parse limit and offset from query params
	
	records, err := d.historyManager.GetConnectionRecords(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(records)
}

// getDataUsageHistory returns data usage history
func (d *DashboardAPI) getDataUsageHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get limit and offset from query parameters
	limit := 50
	offset := 0
	serverID := r.URL.Query().Get("server_id")
	
	// In a real implementation, you would parse limit and offset from query params
	
	records, err := d.historyManager.GetDataUsageRecords(serverID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(records)
}

// getAlertHistory returns alert history
func (d *DashboardAPI) getAlertHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get limit and offset from query parameters
	limit := 50
	offset := 0
	unreadOnly := r.URL.Query().Get("unread") == "true"
	
	// In a real implementation, you would parse limit and offset from query params
	
	records, err := d.historyManager.GetAlertRecords(unreadOnly, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(records)
}

// markAlertAsRead marks an alert as read
func (d *DashboardAPI) markAlertAsRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	alertID := vars["id"]
	
	err := d.historyManager.MarkAlertAsRead(alertID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}