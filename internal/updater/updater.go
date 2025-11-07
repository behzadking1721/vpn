package updater

import (
	"fmt"
	"sync"
	"time"

	"vpnclient/internal/logging"
	"vpnclient/internal/managers"
	"vpnclient/src/core"
)

// Updater handles automatic server updates
type Updater struct {
	serverManager      *managers.ServerManager
	interval           time.Duration
	enabled            bool
	logger             *logging.Logger
	mutex              sync.RWMutex
	stopChan           chan struct{}
}

// Config holds updater configuration
type Config struct {
	Interval time.Duration
	Enabled  bool
}

// NewUpdater creates a new updater instance
func NewUpdater(
	serverManager *managers.ServerManager,
	subscriptionManager *managers.SubscriptionManager,
	logger *logging.Logger,
	config *Config,
) *Updater {
	return &Updater{
		serverManager:      serverManager,
		subscriptionManager: subscriptionManager,
		interval:           config.Interval,
		enabled:            config.Enabled,
		logger:             logger,
		mutex:              sync.RWMutex{},
		stopChan:           make(chan struct{}),
	}
}

// Start begins the updater process
func (u *Updater) Start() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if u.enabled {
		go u.run()
		u.logger.Info("Updater started with interval: %v", u.interval)
	} else {
		u.logger.Info("Updater is disabled")
	}

	return nil
}

// Stop stops the automatic update process
func (u *Updater) Stop() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	close(u.stopChan)
	u.logger.Info("Updater stop signal sent")
}

// SetEnabled enables or disables automatic updates
func (u *Updater) SetEnabled(enabled bool) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.enabled = enabled
	u.logger.Info("Automatic updates %s", map[bool]string{true: "enabled", false: "disabled"}[enabled])
}

// SetInterval sets the update interval
func (u *Updater) SetInterval(interval time.Duration) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.interval = interval
	u.logger.Info("Update interval set to: %v", interval)
}

// UpdateSubscriptions manually triggers subscription updates
func (u *Updater) UpdateSubscriptions() error {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	if !u.enabled {
		return fmt.Errorf("automatic updates are disabled")
	}

	return u.updateSubscriptions()
}

// run executes the update process periodically
func (u *Updater) run() {
	ticker := time.NewTicker(u.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := u.updateServers(); err != nil {
				u.logger.Error("Failed to update servers: %v", err)
			}
		case <-u.stopChan:
			u.logger.Debug("Updater received stop signal")
			return
		}
	}
}

// updateServers fetches updates from all subscriptions and updates the server list
func (u *Updater) updateServers() error {
	u.logger.Debug("Starting server update process")

	// Get all subscriptions
	subscriptions, err := u.subscriptionManager.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %v", err)
	}

	u.logger.Debug("Found %d subscriptions", len(subscriptions))

	// Process each subscription
	for _, sub := range subscriptions {
		if err := u.processSubscription(sub); err != nil {
			u.logger.Error("Failed to process subscription %s: %v", sub.URL, err)
		}
	}

	// Clean up old servers that are no longer in any subscription
	if err := u.cleanupOldServers(subscriptions); err != nil {
		u.logger.Error("Failed to cleanup old servers: %v", err)
	}

	u.logger.Debug("Server update process completed")
	return nil
}

// cleanupOldServers removes servers that are not part of any active subscription
func (u *Updater) cleanupOldServers(subscriptions []*core.Subscription) error {
	u.logger.Debug("Starting cleanup of old servers")

	// Get all servers
	allServers, err := u.serverManager.GetAllServers()
	if err != nil {
		return fmt.Errorf("failed to get all servers: %v", err)
	}

	// For a real implementation, we would need to track which servers belong to which subscriptions
	// For now, we'll just log that we're doing cleanup
	u.logger.Debug("Would clean up servers not in any subscription. Total servers: %d", len(allServers))
	
	return nil
}

// GetStatus returns the current status of the updater
func (u *Updater) GetStatus() map[string]interface{} {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	return map[string]interface{}{
		"enabled":  u.enabled,
		"interval": u.interval.String(),
	}
}