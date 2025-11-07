package api

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"

	"vpnclient/internal/logging"
	"vpnclient/internal/managers"
	"vpnclient/internal/notifications"
)

// Server represents the API server
type Server struct {
	router              *mux.Router
	serverManager       *managers.ServerManager
	connectionManager   *managers.ConnectionManager
	notificationManager *notifications.NotificationManager
	statsManager        *stats.StatsManager
	updater             *updater.Updater
	logger              *logging.Logger
	logFilePath         string
	addr                string
	httpServer          *http.Server
}

// NewServer creates a new API server
func NewServer(
	addr string,
	serverManager *managers.ServerManager,
	connectionManager *managers.ConnectionManager,
	notificationManager *notifications.NotificationManager,
	statsManager *stats.StatsManager,
	updater *updater.Updater,
	logger *logging.Logger,
	logFilePath string,
) *Server {
	s := &Server{
		router:              mux.NewRouter(),
		serverManager:       serverManager,
		connectionManager:   connectionManager,
		notificationManager: notificationManager,
		statsManager:        statsManager,
		updater:             updater,
		logger:              logger,
		logFilePath:         logFilePath,
		addr:                addr,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	api := s.router.PathPrefix("/api").Subrouter()

	// Server management endpoints
	api.HandleFunc("/servers", s.listServers).Methods("GET")
	api.HandleFunc("/servers/enabled", s.listEnabledServers).Methods("GET")
	api.HandleFunc("/servers", s.addServer).Methods("POST")
	api.HandleFunc("/servers/{id}", s.getServer).Methods("GET")
	api.HandleFunc("/servers/{id}", s.updateServer).Methods("PUT")
	api.HandleFunc("/servers/{id}", s.deleteServer).Methods("DELETE")
	api.HandleFunc("/servers/{id}/enable", s.enableServer).Methods("POST")
	api.HandleFunc("/servers/{id}/disable", s.disableServer).Methods("POST")
	api.HandleFunc("/servers/{id}/ping", s.updatePing).Methods("PUT")
	api.HandleFunc("/servers/{id}/test-ping", s.testServerPing).Methods("POST")
	api.HandleFunc("/servers/test-all-ping", s.testAllServersPing).Methods("POST")
	api.HandleFunc("/servers/best", s.getBestServer).Methods("GET")

	// Subscription management endpoints
	api.HandleFunc("/subscriptions", s.addSubscription).Methods("POST")
	api.HandleFunc("/subscriptions", s.getAllSubscriptions).Methods("GET")
	api.HandleFunc("/subscriptions/{id}", s.getSubscription).Methods("GET")
	api.HandleFunc("/subscriptions/{id}", s.updateSubscription).Methods("PUT")
	api.HandleFunc("/subscriptions/{id}", s.deleteSubscription).Methods("DELETE")
	api.HandleFunc("/subscriptions/{id}/update", s.updateSubscriptionServers).Methods("POST")

	// Connection management endpoints
	api.HandleFunc("/connect", s.connect).Methods("POST")
	api.HandleFunc("/connect/fastest", s.connectFastest).Methods("POST")
	api.HandleFunc("/connect/best", s.connectBest).Methods("POST")
	api.HandleFunc("/disconnect", s.disconnect).Methods("POST")
	api.HandleFunc("/status", s.getStatus).Methods("GET")
	api.HandleFunc("/stats", s.getStats).Methods("GET")

	// Statistics endpoints
	api.HandleFunc("/stats/connection", s.getConnectionStats).Methods("GET")
	api.HandleFunc("/stats/sessions", s.getSessionStats).Methods("GET")
	api.HandleFunc("/stats/summary", s.getStatsSummary).Methods("GET")
	api.HandleFunc("/stats/daily", s.getDailyStats).Methods("GET")
	api.HandleFunc("/stats/chart", s.getChartData).Methods("GET")
	api.HandleFunc("/stats/clear", s.clearStats).Methods("POST")

	// Updater endpoints
	api.HandleFunc("/updater/status", s.getUpdaterStatus).Methods("GET")
	api.HandleFunc("/updater/config", s.setUpdaterConfig).Methods("POST")
	api.HandleFunc("/updater/update", s.triggerUpdate).Methods("POST")

	// Notification management endpoints
	api.HandleFunc("/notifications", s.getNotifications).Methods("GET")
	api.HandleFunc("/notifications/unread", s.getUnreadNotifications).Methods("GET")
	api.HandleFunc("/notifications/read", s.markNotificationAsRead).Methods("POST")
	api.HandleFunc("/notifications/read-all", s.markAllNotificationsAsRead).Methods("POST")
	api.HandleFunc("/notifications/clear", s.clearNotifications).Methods("POST")
	api.HandleFunc("/notifications/clear-read", s.clearReadNotifications).Methods("POST")

	// Log management endpoints
	api.HandleFunc("/logs", s.getLogs).Methods("GET")
	api.HandleFunc("/logs/clear", s.clearLogs).Methods("POST")
	api.HandleFunc("/logs/stats", s.getLogStats).Methods("GET")

	// Health check
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")

	// CORS middleware
	s.router.Use(corsMiddleware)
}

// Start starts the API server
func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("API Server starting on %s\n", s.addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Add a small delay to ensure server starts
	time.Sleep(100 * time.Millisecond)
	fmt.Println("API Server started successfully!")
	fmt.Println("Listening on http://localhost:8080")

	return nil
}

// GetRouter returns the router for testing purposes
func (s *Server) GetRouter() *mux.Router {
	return s.router
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getLogs returns the contents of the log file
func (s *Server) getLogs(w http.ResponseWriter, r *http.Request) {
	if s.logger == nil {
		http.Error(w, "Logger not initialized", http.StatusInternalServerError)
		return
	}

	logs, err := s.logger.GetLogs()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading logs: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"logs":  logs,
		"total": len(logs),
	})
}

// clearLogs clears all log entries
func (s *Server) clearLogs(w http.ResponseWriter, r *http.Request) {
	if s.logger == nil {
		http.Error(w, "Logger not initialized", http.StatusInternalServerError)
		return
	}

	if err := s.logger.Clear(); err != nil {
		http.Error(w, fmt.Sprintf("Error clearing logs: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "success",
		"msg":    "Logs cleared successfully",
	})
}

// getLogStats returns statistics about the logs
func (s *Server) getLogStats(w http.ResponseWriter, r *http.Request) {
	if s.logger == nil {
		http.Error(w, "Logger not initialized", http.StatusInternalServerError)
		return
	}

	stats := s.logger.GetStats()
	respondJSON(w, http.StatusOK, stats)
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Close logger if it exists
	if s.logger != nil {
		s.logger.Close()
	}

	return s.httpServer.Shutdown(ctx)
}

// healthCheck handles health check requests
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}
