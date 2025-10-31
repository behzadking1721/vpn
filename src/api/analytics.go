package api

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/analytics"
	"c:/Users/behza/OneDrive/Documents/vpn/src/history"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// AnalyticsAPI handles analytics-related API endpoints
type AnalyticsAPI struct {
	analyticsManager *analytics.AnalyticsManager
	serverManager    *managers.ServerManager
	historyManager   *history.HistoryManager
}

// NewAnalyticsAPI creates a new analytics API instance
func NewAnalyticsAPI(
	analyticsMgr *analytics.AnalyticsManager,
	serverMgr *managers.ServerManager,
	historyMgr *history.HistoryManager) *AnalyticsAPI {

	return &AnalyticsAPI{
		analyticsManager: analyticsMgr,
		serverManager:    serverMgr,
		historyManager:   historyMgr,
	}
}

// RegisterRoutes registers analytics routes with the router
func (a *AnalyticsAPI) RegisterRoutes(router *mux.Router) {
	// Ping analytics endpoints
	router.HandleFunc("/api/analytics/ping", a.getPingStats).Methods("GET")

	// Data usage analytics endpoints
	router.HandleFunc("/api/analytics/usage", a.getDataUsageStats).Methods("GET")
	router.HandleFunc("/api/analytics/usage/daily", a.getDailyDataUsage).Methods("GET")
	router.HandleFunc("/api/analytics/usage/weekly", a.getWeeklyDataUsage).Methods("GET")
	router.HandleFunc("/api/analytics/usage/monthly", a.getMonthlyDataUsage).Methods("GET")

	// Time pattern analytics endpoints
	router.HandleFunc("/api/analytics/time-pattern", a.getTimePatterns).Methods("GET")

	// Report endpoints
	router.HandleFunc("/api/analytics/reports", a.generateReport).Methods("POST")
	router.HandleFunc("/api/analytics/reports/weekly", a.getWeeklyReport).Methods("GET")
	router.HandleFunc("/api/analytics/reports/monthly", a.getMonthlyReport).Methods("GET")

	// Server performance endpoints
	router.HandleFunc("/api/analytics/performance", a.getServerPerformance).Methods("GET")
}

// getPingStats returns ping statistics
func (a *AnalyticsAPI) getPingStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get window parameter (default to 7d)
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "7d"
	}

	stats, err := a.analyticsManager.CalculatePingStats(window)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

// getDataUsageStats returns data usage statistics
func (a *AnalyticsAPI) getDataUsageStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get window parameter (default to 7d)
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "7d"
	}

	stats, err := a.analyticsManager.CalculateDataUsageStats(window)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

// getDailyDataUsage returns daily data usage
func (a *AnalyticsAPI) getDailyDataUsage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get window parameter (default to 30d)
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "30d"
	}

	usage, err := a.analyticsManager.GetDailyDataUsage(window)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(usage)
}

// getWeeklyDataUsage returns weekly data usage
func (a *AnalyticsAPI) getWeeklyDataUsage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get window parameter (default to 90d)
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "90d"
	}

	usage, err := a.analyticsManager.GetWeeklyDataUsage(window)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(usage)
}

// getMonthlyDataUsage returns monthly data usage
func (a *AnalyticsAPI) getMonthlyDataUsage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get window parameter (default to 365d)
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "365d"
	}

	usage, err := a.analyticsManager.GetMonthlyDataUsage(window)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(usage)
}

// getTimePatterns returns time-based patterns
func (a *AnalyticsAPI) getTimePatterns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get window parameter (default to 7d)
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "7d"
	}

	patterns, err := a.analyticsManager.CalculateTimePatterns(window)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(patterns)
}

// generateReport generates a custom report
func (a *AnalyticsAPI) generateReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req analytics.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	report, err := a.analyticsManager.GenerateReport(req.Period)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(report)
}

// getWeeklyReport returns the current weekly report
func (a *AnalyticsAPI) getWeeklyReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	report, err := a.analyticsManager.GenerateReport(analytics.ReportPeriodWeekly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(report)
}

// getMonthlyReport returns the current monthly report
func (a *AnalyticsAPI) getMonthlyReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	report, err := a.analyticsManager.GenerateReport(analytics.ReportPeriodMonthly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(report)
}

// getServerPerformance returns server performance metrics
func (a *AnalyticsAPI) getServerPerformance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get window parameter (default to 7d)
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "7d"
	}

	performance, err := a.analyticsManager.GetServerPerformance(window)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(performance)
}
