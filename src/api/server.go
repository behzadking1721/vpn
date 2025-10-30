package api

import (
	"encoding/json"
	"net/http"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"github.com/gorilla/mux"
)

type APIServer struct {
	serverManager *managers.ServerManager
	connManager   *managers.ConnectionManager
	configManager *managers.ConfigManager
	router        *mux.Router
}

func NewAPIServer(
	serverMgr *managers.ServerManager,
	connMgr *managers.ConnectionManager,
	configMgr *managers.ConfigManager) *APIServer {
	
	api := &APIServer{
		serverManager: serverMgr,
		connManager:   connMgr,
		configManager: configMgr,
		router:        mux.NewRouter(),
	}
	
	api.setupRoutes()
	return api
}

func (a *APIServer) setupRoutes() {
	// Server management endpoints
	a.router.HandleFunc("/api/servers", a.listServers).Methods("GET")
	a.router.HandleFunc("/api/servers", a.addServer).Methods("POST")
	a.router.HandleFunc("/api/servers/{id}", a.getServer).Methods("GET")
	a.router.HandleFunc("/api/servers/{id}", a.updateServer).Methods("PUT")
	a.router.HandleFunc("/api/servers/{id}", a.deleteServer).Methods("DELETE")
	
	// Connection management endpoints
	a.router.HandleFunc("/api/connect", a.connect).Methods("POST")
	a.router.HandleFunc("/api/disconnect", a.disconnect).Methods("POST")
	a.router.HandleFunc("/api/status", a.getStatus).Methods("GET")
	a.router.HandleFunc("/api/stats", a.getStats).Methods("GET")
	
	// Configuration endpoints
	a.router.HandleFunc("/api/config", a.getConfig).Methods("GET")
	a.router.HandleFunc("/api/config", a.updateConfig).Methods("PUT")
	
	// Serve static files (UI)
	a.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/desktop/")))
}

// Server management handlers
func (a *APIServer) listServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	servers := a.serverManager.GetAllServers()
	json.NewEncoder(w).Encode(servers)
}

func (a *APIServer) addServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var server core.Server
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := a.serverManager.AddServer(server); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(server)
}

func (a *APIServer) getServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	serverID := vars["id"]
	
	server, err := a.serverManager.GetServer(serverID)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}
	
	json.NewEncoder(w).Encode(server)
}

func (a *APIServer) updateServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	serverID := vars["id"]
	
	var server core.Server
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Ensure the ID matches
	server.ID = serverID
	
	if err := a.serverManager.UpdateServer(server); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(server)
}

func (a *APIServer) deleteServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	serverID := vars["id"]
	
	if err := a.serverManager.RemoveServer(serverID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// Connection management handlers
func (a *APIServer) connect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		ServerID string `json:"server_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	server, err := a.serverManager.GetServer(req.ServerID)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}
	
	if err := a.connManager.Connect(server); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "connected"})
}

func (a *APIServer) disconnect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := a.connManager.Disconnect(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "disconnected"})
}

func (a *APIServer) getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := a.connManager.GetStatus()
	json.NewEncoder(w).Encode(map[string]core.ConnectionStatus{"status": status})
}

func (a *APIServer) getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	info := a.connManager.GetConnectionInfo()
	json.NewEncoder(w).Encode(info)
}

// Configuration handlers
func (a *APIServer) getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	config := a.configManager.GetConfig()
	json.NewEncoder(w).Encode(config)
}

func (a *APIServer) updateConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var config core.AppConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := a.configManager.UpdateConfig(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(config)
}

// Start the API server
func (a *APIServer) Start(addr string) error {
	return http.ListenAndServe(addr, a.router)
}