package api

import (
	"encoding/json"
	"net/http"
)

// getNotifications returns all notifications
func (s *Server) getNotifications(w http.ResponseWriter, r *http.Request) {
	notifs := s.notificationManager.GetNotifications()
	respondJSON(w, http.StatusOK, notifs)
}

// getUnreadNotifications returns unread notifications
func (s *Server) getUnreadNotifications(w http.ResponseWriter, r *http.Request) {
	notifs := s.notificationManager.GetUnreadNotifications()
	respondJSON(w, http.StatusOK, notifs)
}

// markNotificationAsRead marks a notification as read
func (s *Server) markNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.notificationManager.MarkAsRead(req.ID); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// markAllNotificationsAsRead marks all notifications as read
func (s *Server) markAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	s.notificationManager.MarkAllAsRead()
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// clearNotifications clears all notifications
func (s *Server) clearNotifications(w http.ResponseWriter, r *http.Request) {
	s.notificationManager.ClearNotifications()
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// clearReadNotifications clears read notifications
func (s *Server) clearReadNotifications(w http.ResponseWriter, r *http.Request) {
	s.notificationManager.ClearReadNotifications()
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
