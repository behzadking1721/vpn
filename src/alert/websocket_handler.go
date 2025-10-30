package alert

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/api"
)

// WebSocketAlertHandler handles alert notifications via WebSocket
type WebSocketAlertHandler struct {
	hub *api.WebSocketHub
}

// NewWebSocketAlertHandler creates a new WebSocket alert handler
func NewWebSocketAlertHandler(hub *api.WebSocketHub) *WebSocketAlertHandler {
	return &WebSocketAlertHandler{
		hub: hub,
	}
}

// HandleAlert sends alert notifications via WebSocket
func (w *WebSocketAlertHandler) HandleAlert(alert *Alert) {
	// Create WebSocket message
	message := api.WebSocketMessage{
		Type: "alert",
		Data: map[string]interface{}{
			"id":        alert.ID,
			"type":      string(alert.Severity),
			"title":     alert.Title,
			"message":   alert.Message,
			"timestamp": alert.Timestamp,
			"server_id": alert.ServerID,
		},
	}
	
	// Broadcast the alert to all connected clients
	w.hub.BroadcastMessage(message)
}