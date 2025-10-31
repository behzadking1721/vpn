package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
	"vpn-client/src/backup"
	"vpn-client/src/monitoring"
	"vpn-client/src/security"
)

// SystemAPI handles system-related API endpoints
type SystemAPI struct {
	healthManager  *monitoring.HealthManager
	backupManager  *backup.BackupManager
	encryptManager *security.EncryptionManager
	dbPath         string
}

// NewSystemAPI creates a new system API instance
func NewSystemAPI(
	healthMgr *monitoring.HealthManager,
	backupMgr *backup.BackupManager,
	encryptMgr *security.EncryptionManager,
	dbPath string) *SystemAPI {

	return &SystemAPI{
		healthManager:  healthMgr,
		backupManager:  backupMgr,
		encryptManager: encryptMgr,
		dbPath:         dbPath,
	}
}

// RegisterRoutes registers system routes with the router
func (s *SystemAPI) RegisterRoutes(router *mux.Router) {
	// System health endpoint
	router.HandleFunc("/api/system/health", s.getSystemHealth).Methods("GET")

	// Backup endpoints
	router.HandleFunc("/api/system/backup", s.createBackup).Methods("POST")
	router.HandleFunc("/api/system/backups", s.listBackups).Methods("GET")
	router.HandleFunc("/api/system/backups/{name}", s.deleteBackup).Methods("DELETE")
	router.HandleFunc("/api/system/restore/{name}", s.restoreBackup).Methods("POST")

	// Encryption endpoints
	router.HandleFunc("/api/system/encrypt", s.encryptData).Methods("POST")
	router.HandleFunc("/api/system/decrypt", s.decryptData).Methods("POST")
}

// getSystemHealth returns the system health status
func (s *SystemAPI) getSystemHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Perform health check
	result := s.healthManager.Check()

	// Set status code based on overall health
	if result.Status == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(result)
}

// createBackup creates a new backup
func (s *SystemAPI) createBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request
	var req struct {
		Encrypt bool `json:"encrypt"`
	}

	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Create backup
	backupInfo, err := s.backupManager.CreateBackup(req.Encrypt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create backup: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(backupInfo)
}

// listBackups lists all backups
func (s *SystemAPI) listBackups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// List backups
	backups, err := s.backupManager.ListBackups()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list backups: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(backups)
}

// deleteBackup deletes a backup
func (s *SystemAPI) deleteBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get backup name from URL parameters
	vars := mux.Vars(r)
	name := vars["name"]

	// Delete backup
	if err := s.backupManager.DeleteBackup(name); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete backup: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// restoreBackup restores a backup
func (s *SystemAPI) restoreBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get backup name from URL parameters
	vars := mux.Vars(r)
	name := vars["name"]

	// Restore backup
	if err := s.backupManager.RestoreBackup(name); err != nil {
		http.Error(w, fmt.Sprintf("Failed to restore backup: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "restored"})
}

// encryptData encrypts data
func (s *SystemAPI) encryptData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request
	var req struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if encryption manager is available
	if s.encryptManager == nil {
		http.Error(w, "Encryption is not configured", http.StatusServiceUnavailable)
		return
	}

	// Encrypt data
	encrypted, err := s.encryptManager.EncryptString(req.Data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encrypt data: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"data": encrypted})
}

// decryptData decrypts data
func (s *SystemAPI) decryptData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request
	var req struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if encryption manager is available
	if s.encryptManager == nil {
		http.Error(w, "Encryption is not configured", http.StatusServiceUnavailable)
		return
	}

	// Decrypt data
	decrypted, err := s.encryptManager.DecryptString(req.Data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decrypt data: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"data": decrypted})
}

// downloadBackup allows downloading a backup file
func (s *SystemAPI) downloadBackup(w http.ResponseWriter, r *http.Request) {
	// Get backup name from URL parameters
	vars := mux.Vars(r)
	name := vars["name"]

	// Construct backup file path
	backupPath := filepath.Join(s.backupManager.backupDir, name)

	// Check if file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		http.Error(w, "Backup not found", http.StatusNotFound)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))

	// Serve the file
	http.ServeFile(w, r, backupPath)
}
