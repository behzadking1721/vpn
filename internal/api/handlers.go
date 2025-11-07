package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

// listEnabledServers returns only enabled servers
func (s *Server) listEnabledServers(w http.ResponseWriter, r *http.Request) {
	servers, err := s.serverManager.GetAllServers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Filter only enabled servers
	enabledServers := make([]*core.Server, 0)
	for _, server := range servers {
		if server.Enabled {
			enabledServers = append(enabledServers, server)
		}
	}

	respondJSON(w, http.StatusOK, enabledServers)
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

// testServerPing tests the ping for a specific server
func (s *Server) testServerPing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ping, err := s.serverManager.TestServerPing(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"ping": ping})
}

// testAllServersPing tests the ping for all enabled servers
func (s *Server) testAllServersPing(w http.ResponseWriter, r *http.Request) {
	results, err := s.serverManager.TestAllServersPing()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, results)
}

// getBestServer finds and returns the best server based on comprehensive testing
func (s *Server) getBestServer(w http.ResponseWriter, r *http.Request) {
	server, err := s.serverManager.GetBestServer()
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, server)
}

// connectBest connects to the best server based on comprehensive testing
func (s *Server) connectBest(w http.ResponseWriter, r *http.Request) {
	server, err := s.serverManager.GetBestServer()
	if err != nil {
		respondError(w, http.StatusNotFound, "No servers available: "+err.Error())
		return
	}

	if err := s.connectionManager.Connect(server); err != nil {
		respondError(w, http.StatusInternalServerError, "Connection failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status":    "connected",
		"server_id": server.ID,
		"server":    server.Name,
	})
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

// connectFastest connects to the server with the fastest ping
func (s *Server) connectFastest(w http.ResponseWriter, r *http.Request) {
	servers, err := s.serverManager.GetAllServers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get servers: "+err.Error())
		return
	}

	if len(servers) == 0 {
		respondError(w, http.StatusNotFound, "No servers available")
		return
	}

	// Find the server with the lowest ping
	var fastestServer *core.Server
	for _, server := range servers {
		if !server.Enabled {
			continue
		}

		if fastestServer == nil || server.Ping < fastestServer.Ping {
			fastestServer = server
		}
	}

	if fastestServer == nil {
		respondError(w, http.StatusNotFound, "No enabled servers available")
		return
	}

	if err := s.connectionManager.Connect(fastestServer); err != nil {
		respondError(w, http.StatusInternalServerError, "Connection failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status":    "connected",
		"server_id": fastestServer.ID,
		"server":    fastestServer.Name,
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
		"connected": status == managers.Connected,
		"uptime":    uptime,
		"data_sent": sent,
		"data_recv": received,
		"server_id": serverID,
	})
}

// addSubscription adds a new subscription
func (s *Server) addSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription core.Subscription

	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Generate ID if not provided
	if subscription.ID == "" {
		subscription.ID = generateID()
	}

	// Set timestamps
	now := time.Now()
	if subscription.CreatedAt.IsZero() {
		subscription.CreatedAt = now
	}
	subscription.UpdatedAt = now

	// Parse subscription link and import servers
	if subscription.URL != "" {
		parser := managers.NewSubscriptionParser()
		servers, err := parser.Parse(subscription.URL)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Failed to parse subscription: "+err.Error())
			return
		}

		// Add all parsed servers
		for _, server := range servers {
			// Set server properties from subscription
			server.CreatedAt = now
			server.UpdatedAt = now

			if err := s.serverManager.AddServer(server); err != nil {
				// Log error but continue with other servers
				fmt.Printf("Error adding server: %v\n", err)
			}
		}

		// Update server count
		subscription.ServerCount = len(servers)
	}

	if err := s.serverManager.AddSubscription(&subscription); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, subscription)
}

// getAllSubscriptions returns all subscriptions
func (s *Server) getAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := s.serverManager.GetAllSubscriptions()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, subscriptions)
}

// getSubscription returns a single subscription by ID
func (s *Server) getSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	subscription, err := s.serverManager.GetSubscription(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, subscription)
}

// updateSubscription updates an existing subscription
func (s *Server) updateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var subscription core.Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	subscription.ID = id // Ensure ID matches
	subscription.UpdatedAt = time.Now()

	if err := s.serverManager.UpdateSubscription(&subscription); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, subscription)
}

// deleteSubscription deletes a subscription
func (s *Server) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.serverManager.DeleteSubscription(id); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// updateSubscriptionServers updates servers from a subscription
func (s *Server) updateSubscriptionServers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get subscription
	subscription, err := s.serverManager.GetSubscription(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	// Parse subscription link and import servers
	if subscription.URL != "" {
		parser := managers.NewSubscriptionParser()
		servers, err := parser.Parse(subscription.URL)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Failed to parse subscription: "+err.Error())
			return
		}

		// Add all parsed servers
		for _, server := range servers {
			// Set server properties from subscription
			server.CreatedAt = time.Now()
			server.UpdatedAt = time.Now()

			if err := s.serverManager.AddServer(server); err != nil {
				// Log error but continue with other servers
				fmt.Printf("Error adding server: %v\n", err)
			}
		}

		// Update server count
		subscription.ServerCount = len(servers)
	}

	// Update last update time
	subscription.LastUpdate = time.Now()
	subscription.UpdatedAt = time.Now()

	if err := s.serverManager.UpdateSubscription(subscription); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, subscription)
}

// generateID generates a simple ID (in a real app, use UUID)
func generateID() string {
	// This is a simple ID generator, in a real application you should use UUID
	return time.Now().Format("20060102150405")
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
