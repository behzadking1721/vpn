package api

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/alert"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// AlertAPI handles alert-related API endpoints
type AlertAPI struct {
	alertManager  *alert.AlertManager
	serverManager *managers.ServerManager
}

// NewAlertAPI creates a new alert API instance
func NewAlertAPI(
	alertMgr *alert.AlertManager,
	serverMgr *managers.ServerManager) *AlertAPI {

	return &AlertAPI{
		alertManager:  alertMgr,
		serverManager: serverMgr,
	}
}

// RegisterRoutes registers alert routes with the router
func (a *AlertAPI) RegisterRoutes(router *mux.Router) {
	// Alert rule management endpoints
	router.HandleFunc("/api/alerts/rules", a.listRules).Methods("GET")
	router.HandleFunc("/api/alerts/rules", a.createRule).Methods("POST")
	router.HandleFunc("/api/alerts/rules/{id}", a.getRule).Methods("GET")
	router.HandleFunc("/api/alerts/rules/{id}", a.updateRule).Methods("PUT")
	router.HandleFunc("/api/alerts/rules/{id}", a.deleteRule).Methods("DELETE")

	// Alert management endpoints
	router.HandleFunc("/api/alerts", a.listAlerts).Methods("GET")
	router.HandleFunc("/api/alerts/{id}/read", a.markAlertAsRead).Methods("POST")
	router.HandleFunc("/api/alerts/{id}/resolve", a.resolveAlert).Methods("POST")

	// Alert rule import/export
	router.HandleFunc("/api/alerts/rules/export", a.exportRules).Methods("GET")
	router.HandleFunc("/api/alerts/rules/import", a.importRules).Methods("POST")
}

// listRules returns all alert rules
func (a *AlertAPI) listRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rules := a.alertManager.GetAlertRules()
	json.NewEncoder(w).Encode(rules)
}

// createRule creates a new alert rule
func (a *AlertAPI) createRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rule alert.AlertRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate ID if not provided
	if rule.ID == "" {
		rule.ID = utils.GenerateID()
	}

	// Set timestamps
	now := utils.TimeNow()
	if rule.CreatedAt.IsZero() {
		rule.CreatedAt = now.Time
	}
	rule.UpdatedAt = now.Time

	if err := a.alertManager.AddAlertRule(rule); err != nil {
		if err == alert.ErrRuleAlreadyExists {
			http.Error(w, "Rule already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rule)
}

// getRule returns a specific alert rule
func (a *AlertAPI) getRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ruleID := vars["id"]

	rule, err := a.alertManager.GetAlertRule(ruleID)
	if err != nil {
		if err == alert.ErrRuleNotFound {
			http.Error(w, "Rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rule)
}

// updateRule updates an existing alert rule
func (a *AlertAPI) updateRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ruleID := vars["id"]

	var rule alert.AlertRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure the ID matches
	rule.ID = ruleID

	// Update timestamp
	rule.UpdatedAt = utils.TimeNow().Time

	if err := a.alertManager.UpdateAlertRule(rule); err != nil {
		if err == alert.ErrRuleNotFound {
			http.Error(w, "Rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rule)
}

// deleteRule removes an alert rule
func (a *AlertAPI) deleteRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ruleID := vars["id"]

	if err := a.alertManager.DeleteAlertRule(ruleID); err != nil {
		if err == alert.ErrRuleNotFound {
			http.Error(w, "Rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// listAlerts returns all alerts
func (a *AlertAPI) listAlerts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	unreadOnly := r.URL.Query().Get("unread") == "true"
	unresolvedOnly := r.URL.Query().Get("unresolved") == "true"

	alerts, err := a.alertManager.GetAlerts(unreadOnly, unresolvedOnly, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(alerts)
}

// markAlertAsRead marks an alert as read
func (a *AlertAPI) markAlertAsRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	alertID := vars["id"]

	if err := a.alertManager.MarkAlertAsRead(alertID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "marked as read"})
}

// resolveAlert resolves an alert
func (a *AlertAPI) resolveAlert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	alertID := vars["id"]

	if err := a.alertManager.ResolveAlert(alertID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "resolved"})
}

// exportRules exports alert rules
func (a *AlertAPI) exportRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := a.alertManager.ExportRules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// importRules imports alert rules
func (a *AlertAPI) importRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rules []alert.AlertRule
	if err := json.NewDecoder(r.Body).Decode(&rules); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to JSON for import
	data, err := json.Marshal(rules)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.alertManager.ImportRules(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "imported"})
}
