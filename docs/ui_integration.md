# UI Integration Guide

This document explains how to integrate the frontend UI with the backend Go application.

## Overview

The VPN client application uses a client-server architecture where:
- The backend is written in Go and handles all VPN protocol operations
- The frontend is built with web technologies (HTML/CSS/JS) for cross-platform compatibility
- Communication between frontend and backend happens through a REST API or WebSocket

## Backend API

The Go backend exposes a local HTTP API that the frontend can communicate with. The API includes endpoints for:

### Server Management
- `GET /api/servers` - List all servers
- `POST /api/servers` - Add a new server
- `PUT /api/servers/{id}` - Update a server
- `DELETE /api/servers/{id}` - Remove a server
- `GET /api/servers/{id}` - Get server details

### Connection Management
- `POST /api/connect` - Connect to a server
- `POST /api/disconnect` - Disconnect from current server
- `GET /api/status` - Get connection status
- `GET /api/stats` - Get connection statistics

### Subscription Management
- `POST /api/subscriptions` - Add a subscription
- `PUT /api/subscriptions/{id}` - Update a subscription
- `DELETE /api/subscriptions/{id}` - Remove a subscription

### Configuration
- `GET /api/config` - Get application configuration
- `PUT /api/config` - Update application configuration

## Example API Implementation

Here's how the Go backend would expose these endpoints:

```go
// In a new file: src/api/server.go

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

// Example endpoint implementation
func (a *APIServer) listServers(w http.ResponseWriter, r *http.Request) {
    servers := a.serverManager.GetAllServers()
    json.NewEncoder(w).Encode(servers)
}

func (a *APIServer) connect(w http.ResponseWriter, r *http.Request) {
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

func (a *APIServer) getStatus(w http.ResponseWriter, r *http.Request) {
    status := a.connManager.GetStatus()
    json.NewEncoder(w).Encode(map[string]core.ConnectionStatus{"status": status})
}

// Start the API server
func (a *APIServer) Start(addr string) error {
    return http.ListenAndServe(addr, a.router)
}
```

## Frontend Integration

The frontend communicates with the backend using standard HTTP requests. Here's an example of how the frontend would connect to a server:

```javascript
// In the frontend JavaScript code
async function connectToServer(serverId) {
    try {
        const response = await fetch('http://localhost:8080/api/connect', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ server_id: serverId })
        });
        
        if (response.ok) {
            console.log('Connected successfully');
            updateUIConnectedState();
        } else {
            console.error('Connection failed');
        }
    } catch (error) {
        console.error('Connection error:', error);
    }
}

async function getServerList() {
    try {
        const response = await fetch('http://localhost:8080/api/servers');
        const servers = await response.json();
        updateServerListUI(servers);
    } catch (error) {
        console.error('Failed to fetch servers:', error);
    }
}

// Poll for connection status updates
setInterval(async () => {
    try {
        const response = await fetch('http://localhost:8080/api/status');
        const status = await response.json();
        updateConnectionStatusUI(status.status);
    } catch (error) {
        console.error('Failed to fetch status:', error);
    }
}, 1000);
```

## Desktop Application Packaging

For desktop platforms, the application can be packaged using Electron or Tauri:

### Electron Approach

1. Create a main Electron process that:
   - Starts the Go backend as a child process
   - Serves the frontend UI
   - Manages the application window

2. Example Electron main process:
```javascript
// main.js
const { app, BrowserWindow } = require('electron');
const { spawn } = require('child_process');
const path = require('path');

let backendProcess;
let mainWindow;

function createWindow() {
    mainWindow = new BrowserWindow({
        width: 1200,
        height: 800,
        webPreferences: {
            nodeIntegration: false
        }
    });

    // Load the UI
    mainWindow.loadFile('ui/desktop/index.html');
}

function startBackend() {
    // Start the Go backend
    backendProcess = spawn('./bin/vpn-client', ['--api-port=8080']);
    
    backendProcess.stdout.on('data', (data) => {
        console.log(`Backend: ${data}`);
    });
    
    backendProcess.stderr.on('data', (data) => {
        console.error(`Backend Error: ${data}`);
    });
}

app.whenReady().then(() => {
    startBackend();
    createWindow();
});

app.on('before-quit', () => {
    if (backendProcess) {
        backendProcess.kill();
    }
});
```

### Tauri Approach

1. Configure Tauri to bundle the Go backend
2. Use Tauri's IPC system for communication between frontend and backend

## Mobile Application Packaging

For mobile platforms, you can use Capacitor or Cordova to wrap the web UI in a native shell:

1. Add the necessary plugins for:
   - Background service management
   - System tray integration
   - VPN service control (platform-specific)

## Security Considerations

1. The local API should only listen on localhost (127.0.0.1)
2. Implement proper CORS policies
3. Consider adding authentication for the local API
4. Ensure sensitive data is not exposed through the API

## Next Steps

1. Implement the full backend API as outlined above
2. Enhance the frontend UI with real data from the backend
3. Choose a packaging solution (Electron, Tauri, Capacitor)
4. Implement platform-specific features (system tray, notifications)
5. Add comprehensive error handling and user feedback