package api

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/settings"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// SettingsAPI handles settings-related API endpoints
type SettingsAPI struct {
	settingsManager *settings.SettingsManager
}

// NewSettingsAPI creates a new settings API instance
func NewSettingsAPI(settingsMgr *settings.SettingsManager) *SettingsAPI {
	return &SettingsAPI{
		settingsManager: settingsMgr,
	}
}

// RegisterRoutes registers settings routes with the router
func (s *SettingsAPI) RegisterRoutes(router *mux.Router) {
	// Dashboard settings endpoints
	router.HandleFunc("/api/settings/dashboard", s.getDashboardSettings).Methods("GET")
	router.HandleFunc("/api/settings/dashboard", s.updateDashboardSettings).Methods("POST")
}

// getDashboardSettings returns the current dashboard settings
func (s *SettingsAPI) getDashboardSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	settings := s.settingsManager.GetSettings()
	json.NewEncoder(w).Encode(settings)
}

// updateDashboardSettings updates the dashboard settings
func (s *SettingsAPI) updateDashboardSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newSettings settings.DashboardSettings
	if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.settingsManager.UpdateSettings(newSettings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated settings
	updatedSettings := s.settingsManager.GetSettings()
	json.NewEncoder(w).Encode(updatedSettings)
}
