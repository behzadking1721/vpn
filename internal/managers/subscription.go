package managers

import (
	"fmt"
	"sync"
	"time"

	"vpnclient/internal/database"
	"vpnclient/internal/logging"
	"vpnclient/internal/notifications"
	"vpnclient/src/core"
)

// SubscriptionManager manages subscriptions
type SubscriptionManager struct {
	serverManager      *ServerManager
	store              *database.SubscriptionStoreWrapper
	parser             *SubscriptionParser
	mutex              sync.RWMutex
	notificationManager *notifications.NotificationManager
	logger             *logging.Logger
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager(serverManager *ServerManager, store database.Store) *SubscriptionManager {
	return &SubscriptionManager{
		serverManager: serverManager,
		store:         database.NewSubscriptionStore(store),
		parser:        NewSubscriptionParser(),
	}
}

// SetNotificationManager sets the notification manager
func (sm *SubscriptionManager) SetNotificationManager(nm *notifications.NotificationManager) {
	sm.notificationManager = nm
}

// SetLogger sets the logger
func (sm *SubscriptionManager) SetLogger(logger *logging.Logger) {
	sm.logger = logger
}

// AddSubscription adds a new subscription
func (sm *SubscriptionManager) AddSubscription(url string) (*core.Subscription, error) {
	// Log subscription addition attempt
	if sm.logger != nil {
		sm.logger.Info("Adding subscription from URL: %s", url)
	}
	
	// Parse the subscription
	servers, err := sm.parser.ParseSubscription(url)
	if err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to parse subscription from URL %s: %v", url, err)
		}
		
		// Send error notification
		if sm.notificationManager != nil {
			sm.notificationManager.AddNotification(
				"Subscription Error",
				fmt.Sprintf("Failed to parse subscription: %v", err),
				notifications.Error,
			)
		}
		return nil, fmt.Errorf("failed to parse subscription: %v", err)
	}

	// Create subscription
	sub := &core.Subscription{
		URL:         url,
		ServerCount: len(servers),
		LastUpdate:  time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Add to store
	if err := sm.store.AddSubscription(sub); err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to save subscription from URL %s: %v", url, err)
		}
		
		// Send error notification
		if sm.notificationManager != nil {
			sm.notificationManager.AddNotification(
				"Subscription Error",
				fmt.Sprintf("Failed to save subscription: %v", err),
				notifications.Error,
			)
		}
		return nil, fmt.Errorf("failed to save subscription: %v", err)
	}

	// Add servers
	successCount := 0
	failCount := 0
	for _, server := range servers {
		if err := sm.serverManager.AddServer(server); err != nil {
			// Continue with other servers even if one fails
			failCount++
			if sm.logger != nil {
				sm.logger.Warning("Failed to add server %s: %v", server.Name, err)
			}
		} else {
			successCount++
		}
	}
	
	// Log subscription addition result
	if sm.logger != nil {
		sm.logger.Info("Subscription added successfully: %s (Servers: %d successful, %d failed)", 
			url, successCount, failCount)
	}
	
	// Send success notification
	if sm.notificationManager != nil {
		sm.notificationManager.AddNotification(
			"Subscription Added",
			fmt.Sprintf("Successfully added subscription with %d servers", len(servers)),
			notifications.Success,
		)
	}

	return sub, nil
}

// GetAllSubscriptions returns all subscriptions
func (sm *SubscriptionManager) GetAllSubscriptions() ([]*core.Subscription, error) {
	subs, err := sm.store.GetAllSubscriptions()
	if err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to get all subscriptions: %v", err)
		}
		return nil, err
	}
	
	// Log successful retrieval
	if sm.logger != nil {
		sm.logger.Debug("Retrieved %d subscriptions", len(subs))
	}
	
	return subs, nil
}

// GetSubscription returns a subscription by ID
func (sm *SubscriptionManager) GetSubscription(id string) (*core.Subscription, error) {
	sub, err := sm.store.GetSubscription(id)
	if err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to get subscription %s: %v", id, err)
		}
		return nil, err
	}
	
	// Log successful retrieval
	if sm.logger != nil {
		sm.logger.Debug("Retrieved subscription: %s", sub.URL)
	}
	
	return sub, nil
}

// UpdateSubscription updates a subscription
func (sm *SubscriptionManager) UpdateSubscription(sub *core.Subscription) error {
	sub.UpdatedAt = time.Now()
	
	// Log subscription update
	if sm.logger != nil {
		sm.logger.Info("Updating subscription: %s", sub.URL)
	}
	
	err := sm.store.UpdateSubscription(sub)
	if err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to update subscription %s: %v", sub.URL, err)
		}
		return err
	}
	
	return nil
}

// DeleteSubscription deletes a subscription by ID
func (sm *SubscriptionManager) DeleteSubscription(id string) error {
	// Get subscription for notification
	sub, err := sm.GetSubscription(id)
	if err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to get subscription %s for deletion: %v", id, err)
		}
		return err
	}
	
	if err := sm.store.DeleteSubscription(id); err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to delete subscription %s: %v", id, err)
		}
		
		// Send error notification
		if sm.notificationManager != nil {
			sm.notificationManager.AddNotification(
				"Subscription Error",
				fmt.Sprintf("Failed to delete subscription: %v", err),
				notifications.Error,
			)
		}
		return err
	}
	
	// Log successful deletion
	if sm.logger != nil {
		sm.logger.Info("Subscription deleted successfully: %s", sub.URL)
	}
	
	// Send success notification
	if sm.notificationManager != nil {
		sm.notificationManager.AddNotification(
			"Subscription Deleted",
			fmt.Sprintf("Successfully deleted subscription %s", sub.URL),
			notifications.Success,
		)
	}
	
	return nil
}

// UpdateSubscriptionServers updates servers from a subscription
func (sm *SubscriptionManager) UpdateSubscriptionServers(id string) error {
	// Get subscription
	sub, err := sm.GetSubscription(id)
	if err != nil {
		return err
	}
	
	// Log subscription update attempt
	if sm.logger != nil {
		sm.logger.Info("Updating servers from subscription: %s", sub.URL)
	}

	// Parse the subscription
	servers, err := sm.parser.ParseSubscription(sub.URL)
	if err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to parse subscription %s: %v", sub.URL, err)
		}
		
		// Send error notification
		if sm.notificationManager != nil {
			sm.notificationManager.AddNotification(
				"Subscription Update Error",
				fmt.Sprintf("Failed to parse subscription %s: %v", sub.URL, err),
				notifications.Error,
			)
		}
		return fmt.Errorf("failed to parse subscription: %v", err)
	}

	// Update subscription
	sub.ServerCount = len(servers)
	sub.LastUpdate = time.Now()
	sub.UpdatedAt = time.Now()
	
	if err := sm.UpdateSubscription(sub); err != nil {
		// Log error
		if sm.logger != nil {
			sm.logger.Error("Failed to update subscription %s: %v", sub.URL, err)
		}
		
		// Send error notification
		if sm.notificationManager != nil {
			sm.notificationManager.AddNotification(
				"Subscription Update Error",
				fmt.Sprintf("Failed to update subscription %s: %v", sub.URL, err),
				notifications.Error,
			)
		}
		return err
	}

	// Add/update servers
	updatedCount := 0
	addedCount := 0
	errorCount := 0
	
	for _, server := range servers {
		// Try to get existing server
		existingServer, err := sm.serverManager.GetServer(server.ID)
		if err != nil {
			// Server doesn't exist, add it
			if err := sm.serverManager.AddServer(server); err != nil {
				errorCount++
				if sm.logger != nil {
					sm.logger.Warning("Failed to add server %s: %v", server.Name, err)
				}
			} else {
				addedCount++
			}
		} else {
			// Server exists, update it
			existingServer.Name = server.Name
			existingServer.Host = server.Host
			existingServer.Port = server.Port
			existingServer.Protocol = server.Protocol
			existingServer.Config = server.Config
			
			if err := sm.serverManager.UpdateServer(existingServer); err != nil {
				errorCount++
				if sm.logger != nil {
					sm.logger.Warning("Failed to update server %s: %v", server.Name, err)
				}
			} else {
				updatedCount++
			}
		}
	}
	
	// Log subscription update result
	if sm.logger != nil {
		sm.logger.Info("Subscription update completed: %s (Added: %d, Updated: %d, Errors: %d)", 
			sub.URL, addedCount, updatedCount, errorCount)
	}
	
	// Send success notification
	if sm.notificationManager != nil {
		sm.notificationManager.AddNotification(
			"Subscription Updated",
			fmt.Sprintf("Successfully updated subscription %s with %d servers", sub.URL, len(servers)),
			notifications.Success,
		)
	}

	return nil
}