package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"c:/Users/behza/OneDrive/Documents/vpn/src/history"
	"c:/Users/behza/OneDrive/Documents/vpn/src/alert"
	"c:/Users/behza/OneDrive/Documents/vpn/src/analytics"
	"c:/Users/behza/OneDrive/Documents/vpn/src/settings"
	"c:/Users/behza/OneDrive/Documents/vpn/src/database"
	"c:/Users/behza/OneDrive/Documents/vpn/src/monitoring"
	"c:/Users/behza/OneDrive/Documents/vpn/src/backup"
	"c:/Users/behza/OneDrive/Documents/vpn/src/security"
	"github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

type APIServer struct {
	serverManager    *managers.ServerManager
	connManager      *managers.ConnectionManager
	configManager    *managers.ConfigManager
	dataManager      *managers.DataManager
	historyManager   history.HistoryManager
	alertManager     *alert.AlertManager
	analyticsManager *analytics.AnalyticsManager
	settingsManager  settings.SettingsManager
	dbManager        *database.Manager
	healthManager    *monitoring.HealthManager
	backupManager    *backup.BackupManager
	encryptManager   *security.EncryptionManager
	router           *mux.Router
	dashboardAPI     *DashboardAPI
	webSocketHub     *WebSocketHub
	alertAPI         *AlertAPI
	analyticsAPI     *AnalyticsAPI
	settingsAPI      *SettingsAPI
	systemAPI        *SystemAPI
}

func NewAPIServer(
	serverMgr *managers.ServerManager,
	connMgr *managers.ConnectionManager,
	configMgr *managers.ConfigManager) *APIServer {
	
	// Create database manager
	dbPath := "./data/vpn.db"
	dbManager, err := database.NewManager(dbPath)
	if err != nil {
		// Handle error appropriately in a real implementation
		panic(err)
	}
	
	// Create history repository
	historyRepo := database.NewHistoryRepository(dbManager)
	
	// Create alert repository
	alertRepo := database.NewAlertRepository(dbManager)
	
	// Create settings repository
	settingsRepo := database.NewSettingsRepository(dbManager)
	
	// Create encryption manager (in a real implementation, you would load the password from config)
	encryptManager, _ := security.NewEncryptionManager("default-password")
	
	// Create backup manager
	backupManager := backup.NewBackupManager(dbPath, "./backups", encryptManager)
	
	// Create health manager
	healthManager := monitoring.NewHealthManager()
	
	// Add health checkers
	healthManager.AddChecker(monitoring.NewDatabaseHealthChecker(dbManager))
	healthManager.AddChecker(monitoring.NewSystemHealthChecker())
	
	api := &APIServer{
		serverManager:    serverMgr,
		connManager:      connMgr,
		configManager:    configMgr,
		dataManager:      serverMgr.GetDataManager(), // Get data manager from server manager
		historyManager:   historyRepo,
		dbManager:        dbManager,
		healthManager:    healthManager,
		backupManager:    backupManager,
		encryptManager:   encryptManager,
		router:           mux.NewRouter(),
	}
	
	// Create alert manager with default config
	alertConfig := alert.AlertManagerConfig{
		DesktopNotifications: true,
		EvaluationInterval:   10, // Check every 10 seconds
		HistoryRetention:     30, // Keep history for 30 days
	}
	
	api.alertManager = alert.NewAlertManager(
		serverMgr,
		connMgr,
		api.dataManager,
		api.historyManager,
		alertConfig,
	)
	
	// Set the alert repository for the alert manager
	api.alertManager.SetAlertRepository(alertRepo)
	
	// Create analytics manager
	api.analyticsManager = analytics.NewAnalyticsManager(
		serverMgr,
		api.historyManager,
	)
	
	// Create settings manager
	api.settingsManager = settingsRepo
	
	// Create dashboard API
	api.dashboardAPI = NewDashboardAPI(
		serverMgr, 
		connMgr, 
		configMgr, 
		api.dataManager, 
		api.historyManager,
	)
	
	// Create alert API
	api.alertAPI = NewAlertAPI(
		api.alertManager,
		serverMgr,
	)
	
	// Create analytics API
	api.analyticsAPI = NewAnalyticsAPI(
		api.analyticsManager,
		serverMgr,
		api.historyManager,
	)
	
	// Create settings API
	api.settingsAPI = NewSettingsAPI(api.settingsManager)
	
	// Create system API
	api.systemAPI = NewSystemAPI(
		api.healthManager,
		api.backupManager,
		api.encryptManager,
		dbPath,
	)
	
	// Create WebSocket hub
	api.webSocketHub = NewWebSocketHub(
		serverMgr,
		connMgr,
		api.dataManager,
		configMgr,
	)
	
	// Register WebSocket alert handler
	wsAlertHandler := alert.NewWebSocketAlertHandler(api.webSocketHub)
	api.alertManager.AddAlertHandler(wsAlertHandler)
	
	// Start the WebSocket hub in a goroutine
	go api.webSocketHub.Run()
	
	// Start the alert manager
	api.alertManager.Start()
	
	api.setupRoutes()
	return api
}

func (a *APIServer) setupRoutes() {
	// Server management endpoints
	a.router.HandleFunc("/api/servers", a.listServers).Methods("GET")
	a.router.HandleFunc("/api/servers", a.addServer).Methods("POST")
	a.router.HandleFunc("/api/servers/{id}", a.getServer).Methods("GET")
	a.router.HandleFunc("/api/servers/{id}", a.updateServer).Methods("PUT")
	a.router.HandleFunc("/api/servers/{id}", a.deleteServer).Methods("DELETE")
	
	// Subscription management endpoints
	a.router.HandleFunc("/api/subscriptions", a.listSubscriptions).Methods("GET")
	a.router.HandleFunc("/api/subscriptions", a.addSubscription).Methods("POST")
	a.router.HandleFunc("/api/subscriptions/{id}", a.getSubscription).Methods("GET")
	a.router.HandleFunc("/api/subscriptions/{id}", a.updateSubscription).Methods("PUT")
	a.router.HandleFunc("/api/subscriptions/{id}", a.deleteSubscription).Methods("DELETE")
	a.router.HandleFunc("/api/subscriptions/{id}/update", a.updateSubscriptionServers).Methods("POST")
	
	// Connection management endpoints
	a.router.HandleFunc("/api/connect", a.connect).Methods("POST")
	a.router.HandleFunc("/api/disconnect", a.disconnect).Methods("POST")
	a.router.HandleFunc("/api/status", a.getStatus).Methods("GET")
	a.router.HandleFunc("/api/stats", a.getStats).Methods("GET")
	a.router.HandleFunc("/api/connect/fastest", a.connectToFastest).Methods("POST")
	
	// Data usage management endpoints
	a.router.HandleFunc("/api/datausage", a.getDataUsage).Methods("GET")
	a.router.HandleFunc("/api/datausage/{serverId}", a.getServerDataUsage).Methods("GET")
	a.router.HandleFunc("/api/datausage/reset", a.resetDataUsage).Methods("POST")
	a.router.HandleFunc("/api/datausage/{serverId}/limit", a.setDataLimit).Methods("POST")
	
	// Server utilities endpoints
	a.router.HandleFunc("/api/servers/ping", a.pingAllServers).Methods("POST")
	a.router.HandleFunc("/api/servers/{id}/ping", a.pingServer).Methods("POST")
	
	// Configuration endpoints
	a.router.HandleFunc("/api/config", a.getConfig).Methods("GET")
	a.router.HandleFunc("/api/config", a.updateConfig).Methods("PUT")
	
	// Dashboard endpoints
	a.dashboardAPI.RegisterRoutes(a.router)
	
	// Alert endpoints
	a.alertAPI.RegisterRoutes(a.router)
	
	// Analytics endpoints
	a.analyticsAPI.RegisterRoutes(a.router)
	
	// Settings endpoints
	a.settingsAPI.RegisterRoutes(a.router)
	
	// System endpoints
	a.systemAPI.RegisterRoutes(a.router)
	
	// WebSocket endpoints
	a.webSocketHub.RegisterWebSocketRoutes(a.router)
	
	// Serve static files (UI)
	a.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/desktop/")))
}

// Server management handlers
func (a *APIServer) listServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	servers := a.serverManager.GetAllServers()
	json.NewEncoder(w).Encode(servers)
}

func (a *APIServer) addServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var server core.Server
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Generate ID if not provided
	if server.ID == "" {
		server.ID = utils.GenerateID()
	}
	
	if err := a.serverManager.AddServer(server); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(server)
}

func (a *APIServer) getServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	serverID := vars["id"]
	
	server, err := a.serverManager.GetServer(serverID)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}
	
	json.NewEncoder(w).Encode(server)
}

func (a *APIServer) updateServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	serverID := vars["id"]
	
	var server core.Server
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Ensure the ID matches
	server.ID = serverID
	
	if err := a.serverManager.UpdateServer(server); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(server)
}

func (a *APIServer) deleteServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	serverID := vars["id"]
	
	if err := a.serverManager.RemoveServer(serverID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// Subscription management handlers
func (a *APIServer) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	subscriptions := a.serverManager.GetAllSubscriptions()
	json.NewEncoder(w).Encode(subscriptions)
}

func (a *APIServer) addSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var subscription core.Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Generate ID if not provided
	if subscription.ID == "" {
		subscription.ID = utils.GenerateID()
	}
	
	if err := a.serverManager.AddSubscription(subscription); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

func (a *APIServer) getSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	subID := vars["id"]
	
	subscription, err := a.serverManager.GetSubscription(subID)
	if err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}
	
	json.NewEncoder(w).Encode(subscription)
}

func (a *APIServer) updateSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	subID := vars["id"]
	
	var subscription core.Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Ensure the ID matches
	subscription.ID = subID
	
	if err := a.serverManager.UpdateSubscription(subscription); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(subscription)
}

func (a *APIServer) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	subID := vars["id"]
	
	if err := a.serverManager.RemoveSubscription(subID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func (a *APIServer) updateSubscriptionServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	subID := vars["id"]
	
	// In a real implementation, this would fetch and parse the subscription URL
	// For now, we'll just return a success message
	subscription, err := a.serverManager.GetSubscription(subID)
	if err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}
	
	// Simulate updating servers from subscription
	// In a real implementation, you would parse the subscription URL and update servers
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "updated",
		"subscription": subscription,
	})
}

// Connection management handlers
func (a *APIServer) connect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		ServerID string `json:"server_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	server, err := a.serverManager.GetServer(req.ServerID)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}
	
	// Check data limit if set
	if server.DataLimit > 0 && server.DataUsed >= server.DataLimit {
		http.Error(w, "Data limit exceeded for this server", http.StatusForbidden)
		return
	}
	
	if err := a.connManager.Connect(server); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Add connection record to history
	connRecord := history.ConnectionRecord{
		ID:         utils.GenerateID(),
		ServerID:   server.ID,
		ServerName: server.Name,
		StartTime:  core.TimeNow().Time,
		Protocol:   string(server.Protocol),
		Status:     "connected",
	}
	
	// In a real implementation, you would save this record when the connection is established
	// and update it when the connection is disconnected
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "connected"})
}

func (a *APIServer) disconnect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := a.connManager.Disconnect(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "disconnected"})
}

func (a *APIServer) getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := a.connManager.GetStatus()
	json.NewEncoder(w).Encode(map[string]core.ConnectionStatus{"status": status})
}

func (a *APIServer) getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	info := a.connManager.GetConnectionInfo()
	json.NewEncoder(w).Encode(info)
}

func (a *APIServer) connectToFastest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Find the server with the lowest ping
	servers := a.serverManager.GetAllServers()
	if len(servers) == 0 {
		http.Error(w, "No servers available", http.StatusNotFound)
		return
	}
	
	// Filter enabled servers
	enabledServers := make([]core.Server, 0)
	for _, server := range servers {
		// Check data limit if set
		if server.DataLimit > 0 && server.DataUsed >= server.DataLimit {
			continue // Skip servers that have exceeded their data limit
		}
		
		if server.Enabled {
			enabledServers = append(enabledServers, server)
		}
	}
	
	if len(enabledServers) == 0 {
		http.Error(w, "No enabled servers available or all servers have exceeded data limits", http.StatusNotFound)
		return
	}
	
	// Find server with lowest ping
	fastestServer := enabledServers[0]
	for _, server := range enabledServers {
		if server.Ping < fastestServer.Ping {
			fastestServer = server
		}
	}
	
	// Connect to the fastest server
	if err := a.connManager.Connect(fastestServer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "connected",
		"server": fastestServer,
	})
}

// Data usage management handlers
func (a *APIServer) getDataUsage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get all servers and their data usage
	servers := a.serverManager.GetAllServers()
	usageData := make([]map[string]interface{}, 0)
	
	totalSent := int64(0)
	totalReceived := int64(0)
	
	for _, server := range servers {
		serverUsage := map[string]interface{}{
			"server_id":      server.ID,
			"server_name":    server.Name,
			"data_sent":      server.DataUsed,
			"data_limit":     server.DataLimit,
			"data_remaining": server.DataLimit - server.DataUsed,
			"limit_exceeded": server.DataLimit > 0 && server.DataUsed >= server.DataLimit,
		}
		
		usageData = append(usageData, serverUsage)
		
		totalSent += server.DataUsed
		totalReceived += server.DataUsed // In a real implementation, you might have separate sent/received values
	}
	
	response := map[string]interface{}{
		"servers":      usageData,
		"total_sent":   totalSent,
		"total_received": totalReceived,
	}
	
	json.NewEncoder(w).Encode(response)
}

func (a *APIServer) getServerDataUsage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	serverID := vars["serverId"]
	
	server, err := a.serverManager.GetServer(serverID)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}
	
	usageData := map[string]interface{}{
		"server_id":      server.ID,
		"server_name":    server.Name,
		"data_sent":      server.DataUsed,
		"data_limit":     server.DataLimit,
		"data_remaining": server.DataLimit - server.DataUsed,
		"limit_exceeded": server.DataLimit > 0 && server.DataUsed >= server.DataLimit,
	}
	
	json.NewEncoder(w).Encode(usageData)
}

func (a *APIServer) resetDataUsage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Reset data usage for all servers
	servers := a.serverManager.GetAllServers()
	for _, server := range servers {
		server.DataUsed = 0
		a.serverManager.UpdateServer(server)
	}
	
	json.NewEncoder(w).Encode(map[string]string{"status": "data usage reset for all servers"})
}

func (a *APIServer) setDataLimit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	serverID := vars["serverId"]
	
	var req struct {
		DataLimit int64 `json:"data_limit"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	server, err := a.serverManager.GetServer(serverID)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}
	
	// Update data limit
	server.DataLimit = req.DataLimit
	a.serverManager.UpdateServer(server)
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"server_id":   server.ID,
		"data_limit":  server.DataLimit,
		"data_used":   server.DataUsed,
		"data_remaining": server.DataLimit - server.DataUsed,
	})
}

// Server utilities handlers
func (a *APIServer) pingAllServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	servers := a.serverManager.GetAllServers()
	results := make([]map[string]interface{}, 0)
	
	for _, server := range servers {
		// Simulate ping
		pingResult := utils.PingServer(server.Host)
		
		// Update server ping
		server.Ping = pingResult
		server.LastPing = core.TimeNow()
		a.serverManager.UpdateServer(server)
		
		results = append(results, map[string]interface{}{
			"server_id": server.ID,
			"server_name": server.Name,
			"ping": pingResult,
		})
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "completed",
		"results": results,
	})
}

func (a *APIServer) pingServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	serverID := vars["id"]
	
	server, err := a.serverManager.GetServer(serverID)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}
	
	// Simulate ping
	pingResult := utils.PingServer(server.Host)
	
	// Update server ping
	server.Ping = pingResult
	server.LastPing = core.TimeNow()
	a.serverManager.UpdateServer(server)
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"server_id": server.ID,
		"server_name": server.Name,
		"ping": pingResult,
	})
}

// Configuration handlers
func (a *APIServer) getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	config := a.configManager.GetConfig()
	json.NewEncoder(w).Encode(config)
}

func (a *APIServer) updateConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var config core.AppConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := a.configManager.UpdateConfig(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(config)
}

// Start the API server
func (a *APIServer) Start(addr string) error {
	return http.ListenAndServe(addr, a.router)
}

// Close closes the API server and its resources
func (a *APIServer) Close() error {
	if a.dbManager != nil {
		return a.dbManager.Close()
	}
	return nil
}