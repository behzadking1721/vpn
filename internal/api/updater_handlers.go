package api

import (
	"net/http"
	"time"

	"vpnclient/internal/updater"
)

// getUpdaterStatus returns the status of the automatic updater
func (s *Server) getUpdaterStatus(w http.ResponseWriter, r *http.Request) {
	status := s.updater.GetStatus()
	respondJSON(w, http.StatusOK, status)
}

// setUpdaterConfig updates the updater configuration
func (s *Server) setUpdaterConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request
	var config struct {
		Enabled  *bool   `json:"enabled,omitempty"`
		Interval *string `json:"interval,omitempty"`
	}

	if err := parseJSONBody(r, &config); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Update configuration
	if config.Enabled != nil {
		s.updater.SetEnabled(*config.Enabled)
	}

	if config.Interval != nil {
		if duration, err := time.ParseDuration(*config.Interval); err == nil {
			s.updater.SetInterval(duration)
		} else {
			respondError(w, http.StatusBadRequest, "Invalid interval format: "+err.Error())
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "success",
		"message": "Updater configuration updated successfully",
	})
}

// triggerUpdate manually triggers a subscription update
func (s *Server) triggerUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := s.updater.UpdateSubscriptions(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update subscriptions: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "success",
		"message": "Subscription update triggered successfully",
	})
}