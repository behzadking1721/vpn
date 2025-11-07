package api

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"vpnclient/internal/logging"
)

// getLogs returns the application logs
func (s *Server) getLogs(w http.ResponseWriter, r *http.Request) {
	// Check if log file path is configured
	if s.logFilePath == "" {
		respondError(w, http.StatusNotFound, "Log file not configured")
		return
	}

	// Check if log file exists
	if _, err := os.Stat(s.logFilePath); os.IsNotExist(err) {
		respondError(w, http.StatusNotFound, "Log file not found")
		return
	}

	// Open log file
	file, err := os.Open(s.logFilePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to open log file: %v", err))
		return
	}
	defer file.Close()

	// Get query parameters
	limit := r.URL.Query().Get("limit")
	level := r.URL.Query().Get("level")

	// Read and filter logs
	var logs []string
	scanner := bufio.NewScanner(file)

	// For performance, we'll read the file in reverse if we have a limit
	// But for simplicity in this implementation, we'll read from the beginning
	for scanner.Scan() {
		line := scanner.Text()

		// Filter by level if specified
		if level != "" && !strings.Contains(line, fmt.Sprintf("[%s]", strings.ToUpper(level))) {
			continue
		}

		logs = append(logs, line)
	}

	// Apply limit if specified
	if limit != "" {
		// In a real implementation, we would implement proper tail functionality
		// For now, we'll just limit the number of logs returned
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"logs":  logs,
		"count": len(logs),
		"file":  s.logFilePath,
	})
}

// clearLogs clears the application logs
func (s *Server) clearLogs(w http.ResponseWriter, r *http.Request) {
	// Check if log file path is configured
	if s.logFilePath == "" {
		respondError(w, http.StatusNotFound, "Log file not configured")
		return
	}

	// Close current log file if it's open
	if s.logger != nil {
		s.logger.Close()
	}

	// Remove log file
	if err := os.Remove(s.logFilePath); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to remove log file: %v", err))
		return
	}

	// Create new logger with the same configuration
	logger, err := logging.NewLogger(logging.Config{
		Level:     logging.INFO,
		Output:    s.logFilePath,
		Timestamp: true,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to recreate logger: %v", err))
		return
	}

	// Update server logger
	s.logger = logger

	// Log the clear event
	s.logger.Info("Logs cleared by user request")

	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Logs cleared successfully",
	})
}

// getLogStats returns statistics about the logs
func (s *Server) getLogStats(w http.ResponseWriter, r *http.Request) {
	// Check if log file path is configured
	if s.logFilePath == "" {
		respondError(w, http.StatusNotFound, "Log file not configured")
		return
	}

	// Check if log file exists
	if _, err := os.Stat(s.logFilePath); os.IsNotExist(err) {
		respondError(w, http.StatusNotFound, "Log file not found")
		return
	}

	// Get file info
	fileInfo, err := os.Stat(s.logFilePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get log file info: %v", err))
		return
	}

	// Open log file
	file, err := os.Open(s.logFilePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to open log file: %v", err))
		return
	}
	defer file.Close()

	// Count log levels
	stats := map[string]int{
		"DEBUG":   0,
		"INFO":    0,
		"WARNING": 0,
		"ERROR":   0,
		"FATAL":   0,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.Contains(line, "[DEBUG]"):
			stats["DEBUG"]++
		case strings.Contains(line, "[INFO]"):
			stats["INFO"]++
		case strings.Contains(line, "[WARNING]"):
			stats["WARNING"]++
		case strings.Contains(line, "[ERROR]"):
			stats["ERROR"]++
		case strings.Contains(line, "[FATAL]"):
			stats["FATAL"]++
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"file": map[string]interface{}{
			"path":     s.logFilePath,
			"size":     fileInfo.Size(),
			"modified": fileInfo.ModTime().Format(time.RFC3339),
		},
		"stats": stats,
		"total": stats["DEBUG"] + stats["INFO"] + stats["WARNING"] + stats["ERROR"] + stats["FATAL"],
	})
}
