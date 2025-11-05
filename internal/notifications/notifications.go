package notifications

import (
	"fmt"
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	Info    NotificationType = "info"
	Success NotificationType = "success"
	Warning NotificationType = "warning"
	Error   NotificationType = "error"
)

// Notification represents a system notification
type Notification struct {
	ID        string           `json:"id"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Type      NotificationType `json:"type"`
	Timestamp time.Time        `json:"timestamp"`
	Read      bool             `json:"read"`
}

// NotificationManager manages system notifications
type NotificationManager struct {
	notifications []Notification
	maxNotifications int
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(maxNotifications int) *NotificationManager {
	return &NotificationManager{
		notifications:    make([]Notification, 0),
		maxNotifications: maxNotifications,
	}
}

// AddNotification adds a new notification
func (nm *NotificationManager) AddNotification(title, message string, notifType NotificationType) *Notification {
	notification := Notification{
		ID:        fmt.Sprintf("notif_%d", time.Now().UnixNano()),
		Title:     title,
		Message:   message,
		Type:      notifType,
		Timestamp: time.Now(),
		Read:      false,
	}
	
	// Add to the beginning of the slice
	nm.notifications = append([]Notification{notification}, nm.notifications...)
	
	// Limit the number of notifications
	if len(nm.notifications) > nm.maxNotifications {
		nm.notifications = nm.notifications[:nm.maxNotifications]
	}
	
	return &notification
}

// GetNotifications returns all notifications
func (nm *NotificationManager) GetNotifications() []Notification {
	return nm.notifications
}

// GetUnreadNotifications returns unread notifications
func (nm *NotificationManager) GetUnreadNotifications() []Notification {
	var unread []Notification
	for _, notif := range nm.notifications {
		if !notif.Read {
			unread = append(unread, notif)
		}
	}
	return unread
}

// MarkAsRead marks a notification as read
func (nm *NotificationManager) MarkAsRead(id string) error {
	for i, notif := range nm.notifications {
		if notif.ID == id {
			nm.notifications[i].Read = true
			return nil
		}
	}
	return fmt.Errorf("notification not found")
}

// MarkAllAsRead marks all notifications as read
func (nm *NotificationManager) MarkAllAsRead() {
	for i := range nm.notifications {
		nm.notifications[i].Read = true
	}
}

// ClearNotifications clears all notifications
func (nm *NotificationManager) ClearNotifications() {
	nm.notifications = make([]Notification, 0)
}

// ClearReadNotifications clears read notifications
func (nm *NotificationManager) ClearReadNotifications() {
	var unread []Notification
	for _, notif := range nm.notifications {
		if !notif.Read {
			unread = append(unread, notif)
		}
	}
	nm.notifications = unread
}