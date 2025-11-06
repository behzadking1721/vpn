package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type Server struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	UUID     string `json:"uuid"`
	Security string `json:"security"`
	Enabled  bool   `json:"enabled,omitempty"`
	Ping     int    `json:"ping,omitempty"`
	ID       string `json:"id,omitempty"`
}

type Config struct {
	Servers     []Server `json:"servers"`
	LogLevel    string   `json:"log_level"`
	AutoConnect bool     `json:"auto_connect"`
}

var (
	cfg        Config
	cfgMu      sync.RWMutex
	configPath = "config/settings.json"
)

func loadConfig() error {
	b, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}
	cfgMu.Lock()
	cfg = c
	cfgMu.Unlock()
	return nil
}

func saveConfig() error {
	cfgMu.RLock()
	b, err := json.MarshalIndent(cfg, "", "  ")
	cfgMu.RUnlock()
	if err != nil {
		return err
	}
	// ensure directory exists
	d := path.Dir(configPath)
	if d != "." {
		_ = os.MkdirAll(d, 0755)
	}
	return os.WriteFile(configPath, b, 0644)
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cfg)
}

func handlePutConfig(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var incoming Config
	if err := json.Unmarshal(body, &incoming); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid json"))
		return
	}
	cfgMu.Lock()
	cfg = incoming
	cfgMu.Unlock()
	if err := saveConfig(); err != nil {
		log.Printf("warning: could not save config: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(cfg)
}

func handleGetServers(w http.ResponseWriter, r *http.Request) {
	cfgMu.RLock()
	servers := cfg.Servers
	cfgMu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(servers)
}

func handlePostServer(w http.ResponseWriter, r *http.Request) {
	var s Server
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// ensure an ID
	if s.ID == "" {
		s.ID = s.UUID
	}
	cfgMu.Lock()
	cfg.Servers = append(cfg.Servers, s)
	cfgMu.Unlock()
	_ = saveConfig()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(s)
}

func findServerIndexByID(id string) int {
	for i, s := range cfg.Servers {
		if s.ID == id || s.UUID == id {
			return i
		}
	}
	return -1
}

func handleServerByID(w http.ResponseWriter, r *http.Request) {
	// path: /api/servers/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := parts[2]

	switch r.Method {
	case http.MethodPut:
		var s Server
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cfgMu.Lock()
		idx := findServerIndexByID(id)
		if idx == -1 {
			cfgMu.Unlock()
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// preserve ID
		s.ID = cfg.Servers[idx].ID
		cfg.Servers[idx] = s
		cfgMu.Unlock()
		_ = saveConfig()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(s)
	case http.MethodDelete:
		cfgMu.Lock()
		idx := findServerIndexByID(id)
		if idx == -1 {
			cfgMu.Unlock()
			w.WriteHeader(http.StatusNotFound)
			return
		}
		cfg.Servers = append(cfg.Servers[:idx], cfg.Servers[idx+1:]...)
		cfgMu.Unlock()
		_ = saveConfig()
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	if err := loadConfig(); err != nil {
		log.Printf("could not load config from %s: %v", configPath, err)
		// start with empty config
		cfg = Config{Servers: []Server{}, LogLevel: "info", AutoConnect: false}
	}

	http.HandleFunc("/api/config", withCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetConfig(w, r)
		case http.MethodPut:
			handlePutConfig(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/servers", withCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetServers(w, r)
		case http.MethodPost:
			handlePostServer(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/servers/", withCORS(handleServerByID))

	addr := ":8080"
	log.Printf("mock_server: listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
