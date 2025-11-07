package api

import (
	"net/http"
	"strconv"
)

// getConnectionStats returns current connection statistics
func (s *Server) getConnectionStats(w http.ResponseWriter, r *http.Request) {
	stat := s.statsManager.GetCurrentConnection()
	if stat == nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"connected": false,
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"connected":   true,
		"timestamp":   stat.Timestamp,
		"data_sent":   stat.DataSent,
		"data_recv":   stat.DataRecv,
		"server_id":   stat.ServerID,
		"server_name": stat.ServerName,
	})
}

// getSessionStats returns session statistics
func (s *Server) getSessionStats(w http.ResponseWriter, r *http.Request) {
	sessions := s.statsManager.GetSessions()
	respondJSON(w, http.StatusOK, sessions)
}

// getStatsSummary returns a summary of statistics
func (s *Server) getStatsSummary(w http.ResponseWriter, r *http.Request) {
	// Get total data usage
	totalSent, totalRecv := s.statsManager.GetTotalDataUsage()

	// Get current connection
	current := s.statsManager.GetCurrentConnection()

	// Get recent sessions (last 10)
	allSessions := s.statsManager.GetSessions()
	recentSessions := allSessions
	if len(allSessions) > 10 {
		recentSessions = allSessions[len(allSessions)-10:]
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"total_data_sent":    totalSent,
		"total_data_recv":    totalRecv,
		"current_connection": current,
		"recent_sessions":    recentSessions,
	})
}

// getDailyStats returns daily statistics
func (s *Server) getDailyStats(w http.ResponseWriter, r *http.Request) {
	// Get days parameter (default to 7)
	days := 7
	daysParam := r.URL.Query().Get("days")
	if daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
			days = d
		}
	}

	dailyStats := s.statsManager.GetDailyDataUsage(days)
	respondJSON(w, http.StatusOK, dailyStats)
}

// clearStats clears all statistics
func (s *Server) clearStats(w http.ResponseWriter, r *http.Request) {
	s.statsManager.ClearStats()
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Statistics cleared successfully",
	})
}

// getChartData returns data formatted for charting
func (s *Server) getChartData(w http.ResponseWriter, r *http.Request) {
	chartType := r.URL.Query().Get("type")

	switch chartType {
	case "daily_usage":
		s.handleDailyUsageChart(w, r)
	case "session_duration":
		s.handleSessionDurationChart(w, r)
	case "data_comparison":
		s.handleDataComparisonChart(w, r)
	default:
		// Default to daily usage
		s.handleDailyUsageChart(w, r)
	}
}

// handleDailyUsageChart handles daily usage chart data
func (s *Server) handleDailyUsageChart(w http.ResponseWriter, r *http.Request) {
	days := 7
	daysParam := r.URL.Query().Get("days")
	if daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
			days = d
		}
	}

	dailyStats := s.statsManager.GetDailyDataUsage(days)

	// Format data for charting
	labels := make([]string, len(dailyStats))
	sentData := make([]int64, len(dailyStats))
	recvData := make([]int64, len(dailyStats))

	for i, stat := range dailyStats {
		labels[i] = stat.Timestamp.Format("2006-01-02")
		sentData[i] = stat.DataSent
		recvData[i] = stat.DataRecv
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"type":   "daily_usage",
		"labels": labels,
		"datasets": []map[string]interface{}{
			{
				"label":           "Data Sent",
				"data":            sentData,
				"backgroundColor": "rgba(54, 162, 235, 0.2)",
				"borderColor":     "rgba(54, 162, 235, 1)",
				"borderWidth":     1,
			},
			{
				"label":           "Data Received",
				"data":            recvData,
				"backgroundColor": "rgba(255, 99, 132, 0.2)",
				"borderColor":     "rgba(255, 99, 132, 1)",
				"borderWidth":     1,
			},
		},
	})
}

// handleSessionDurationChart handles session duration chart data
func (s *Server) handleSessionDurationChart(w http.ResponseWriter, r *http.Request) {
	sessions := s.statsManager.GetSessions()

	// Take last 10 sessions for the chart
	if len(sessions) > 10 {
		sessions = sessions[len(sessions)-10:]
	}

	labels := make([]string, len(sessions))
	durations := make([]float64, len(sessions))

	for i, session := range sessions {
		labels[i] = session.ServerName
		duration := session.EndedAt.Sub(session.StartedAt)
		durations[i] = duration.Seconds()
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"type":   "session_duration",
		"labels": labels,
		"datasets": []map[string]interface{}{
			{
				"label":           "Session Duration (seconds)",
				"data":            durations,
				"backgroundColor": "rgba(75, 192, 192, 0.2)",
				"borderColor":     "rgba(75, 192, 192, 1)",
				"borderWidth":     1,
			},
		},
	})
}

// handleDataComparisonChart handles data comparison chart data
func (s *Server) handleDataComparisonChart(w http.ResponseWriter, r *http.Request) {
	totalSent, totalRecv := s.statsManager.GetTotalDataUsage()

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"type":   "data_comparison",
		"labels": []string{"Data Sent", "Data Received"},
		"datasets": []map[string]interface{}{
			{
				"label": "Total Data Usage",
				"data":  []int64{totalSent, totalRecv},
				"backgroundColor": []string{
					"rgba(255, 99, 132, 0.2)",
					"rgba(54, 162, 235, 0.2)",
				},
				"borderColor": []string{
					"rgba(255, 99, 132, 1)",
					"rgba(54, 162, 235, 1)",
				},
				"borderWidth": 1,
			},
		},
	})
}
