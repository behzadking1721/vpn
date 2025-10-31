package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// WebSocketHub maintains the set of active clients and broadcasts messages to the clients
type WebSocketHub struct {
	// Registered clients
	clients map[*WebSocketClient]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *WebSocketClient

	// Unregister requests from clients
	unregister chan *WebSocketClient

	// Managers
	serverManager *managers.ServerManager
	connManager   *managers.ConnectionManager
	dataManager   *managers.DataManager
	configManager *managers.ConfigManager

	// Mutex for thread safety
	mutex sync.RWMutex
}

// NewWebSocketHub creates a new WebSocket hub
func NewWebSocketHub(
	serverMgr *managers.ServerManager,
	connMgr *managers.ConnectionManager,
	dataMgr *managers.DataManager,
	configMgr *managers.ConfigManager) *WebSocketHub {
	return &WebSocketHub{
		broadcast:     make(chan []byte),
		register:      make(chan *WebSocketClient),
		unregister:    make(chan *WebSocketClient),
		clients:       make(map[*WebSocketClient]bool),
		serverManager: serverMgr,
		connManager:   connMgr,
		dataManager:   dataMgr,
		configManager: configMgr,
	}
}

// Run starts the hub
func (h *WebSocketHub) Run() {
	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()

		case <-ticker.C:
			// Send dashboard updates periodically
			h.sendDashboardUpdates()
		}
	}
}

// sendDashboardUpdates sends dashboard updates to all connected clients
func (h *WebSocketHub) sendDashboardUpdates() {
	// Get connection status
	status := h.connManager.GetStatus()
	connInfo := h.connManager.GetConnectionInfo()

	connectionStatus := ConnectionStatusResponse{
		Connected:    status == core.StatusConnected,
		DataSent:     connInfo.DataSent,
		DataReceived: connInfo.DataReceived,
	}

	if status == core.StatusConnected {
		currentServer := h.connManager.GetCurrentServer()
		if currentServer != nil {
			connectionStatus.Server = currentServer
			connectionStatus.ConnectionTime = connInfo.StartTime
		}
	}

	// Prepare the message
	message := WebSocketMessage{
		Type: "dashboard_update",
		Data: connectionStatus,
	}

	// Send to all clients
	h.broadcastMessage(message)
}

// broadcastMessage broadcasts a message to all clients
func (h *WebSocketHub) broadcastMessage(msg WebSocketMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mutex.RLock()
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
	h.mutex.RUnlock()
}

// WebSocketClient is a middleman between the websocket connection and the hub
type WebSocketClient struct {
	hub *WebSocketHub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte
}

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// WebSocketAlert represents an alert message
type WebSocketAlert struct {
	ID      string `json:"id"`
	Type    string `json:"type"` // "warning", "error", "info"
	Title   string `json:"title"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in development
		// In production, you should be more restrictive
		return true
	},
}

// RegisterWebSocketRoutes registers WebSocket routes
func (h *WebSocketHub) RegisterWebSocketRoutes(router *mux.Router) {
	router.HandleFunc("/api/ws/dashboard", func(w http.ResponseWriter, r *http.Request) {
		h.serveWs(w, r)
	})
}

// serveWs handles websocket requests from the peer
func (h *WebSocketHub) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &WebSocketClient{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines
	go client.writePump()
	go client.readPump()
}

// writePump pumps messages from the hub to the websocket connection
func (c *WebSocketClient) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *WebSocketClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}

		// Handle incoming messages if needed
		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		// Process message based on type
		switch msg.Type {
		case "ping":
			// Respond to ping with pong
			response := WebSocketMessage{
				Type: "pong",
				Data: time.Now().Unix(),
			}
			data, _ := json.Marshal(response)
			c.send <- data
		}
	}
}
