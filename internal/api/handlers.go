package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"vpnclient/internal/managers"
	"vpnclient/src/core"
)

// listServers returns all servers
func (s *Server) listServers(w http.ResponseWriter, r *http.Request) {
	servers, err := s.serverManager.GetAllServers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, servers)
}

// getServer returns a single server by ID
func (s *Server) getServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	server, err := s.serverManager.GetServer(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, server)
}

// addServer adds a new server
func (s *Server) addServer(w http.ResponseWriter, r *http.Request) {
	var server core.Server

	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := s.serverManager.AddServer(&server); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, server)
}

// updateServer updates an existing server
func (s *Server) updateServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var server core.Server
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	server.ID = id // Ensure ID matches

	if err := s.serverManager.UpdateServer(&server); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, server)
}

// deleteServer deletes a server
func (s *Server) deleteServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.serverManager.DeleteServer(id); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// enableServer enables a server
func (s *Server) enableServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.serverManager.EnableServer(id); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "enabled"})
}

// disableServer disables a server
func (s *Server) disableServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.serverManager.DisableServer(id); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "disabled"})
}

// updatePing updates server ping
func (s *Server) updatePing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Ping int `json:"ping"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := s.serverManager.UpdatePing(id, req.Ping); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"ping": req.Ping})
}

// connect connects to a server
func (s *Server) connect(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerID string `json:"server_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if req.ServerID == "" {
		respondError(w, http.StatusBadRequest, "server_id is required")
		return
	}

	server, err := s.serverManager.GetServer(req.ServerID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Server not found: "+err.Error())
		return
	}

	// Convert server to interface for Connect
	// Note: ConnectionManager.Connect expects interface{}, so we pass the server pointer
	if err := s.connectionManager.Connect(server); err != nil {
		respondError(w, http.StatusInternalServerError, "Connection failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status":    "connected",
		"server_id": req.ServerID,
	})
}

// disconnect disconnects from current server
func (s *Server) disconnect(w http.ResponseWriter, r *http.Request) {
	if err := s.connectionManager.Disconnect(); err != nil {
		respondError(w, http.StatusInternalServerError, "Disconnection failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "disconnected"})
}

// getStatus returns connection status
func (s *Server) getStatus(w http.ResponseWriter, r *http.Request) {
	status := s.connectionManager.GetStatus()
	statusStr := s.connectionManager.GetStatusString()

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":      statusStr,
		"status_code": int(status),
	})
}

// getStats returns connection statistics
func (s *Server) getStats(w http.ResponseWriter, r *http.Request) {
	status := s.connectionManager.GetStatus()
	sent, received := s.connectionManager.GetDataUsage()
	uptime := s.connectionManager.GetUptime()
	currentServer := s.connectionManager.GetCurrentServer()
	
	var serverID string
	if currentServer != nil {
		serverID = currentServer.ID
	}
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"connected":  status == managers.Connected,
		"uptime":     uptime,
		"data_sent":  sent,
		"data_recv":  received,
		"server_id":  serverID,
	})
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

