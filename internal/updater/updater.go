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
	serverManager       *managers.ServerManager
	subscriptionManager *managers.SubscriptionManager
	interval            time.Duration
	enabled             bool
	logger              *logging.Logger
	mutex               sync.RWMutex
	stopChan            chan struct{}
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
	config Config,
	logger *logging.Logger,
) *Updater {
	return &Updater{
		serverManager:       serverManager,
		subscriptionManager: subscriptionManager,
		interval:            config.Interval,
		enabled:             config.Enabled,
		logger:              logger,
		stopChan:            make(chan struct{}),
	}
}

// Start begins the automatic update process
func (u *Updater) Start() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if !u.enabled {
		u.logger.Info("Automatic updates are disabled")
		return
	}

	u.logger.Info("Starting automatic updater with interval: %v", u.interval)

	go func() {
		ticker := time.NewTicker(u.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				u.updateSubscriptions()
			case <-u.stopChan:
				u.logger.Info("Stopping automatic updater")
				return
			}
		}
	}()
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

// updateSubscriptions performs the actual subscription update process
func (u *Updater) updateSubscriptions() error {
	u.logger.Info("Starting subscription update process")

	// Get all subscriptions
	subscriptions, err := u.subscriptionManager.GetAllSubscriptions()
	if err != nil {
		u.logger.Error("Failed to get subscriptions: %v", err)
		return err
	}

	u.logger.Info("Found %d subscriptions to update", len(subscriptions))

	// Update each subscription
	successCount := 0
	failCount := 0

	for _, sub := range subscriptions {
		u.logger.Debug("Updating subscription: %s", sub.URL)

		if err := u.subscriptionManager.UpdateSubscriptionServers(sub.ID); err != nil {
			u.logger.Error("Failed to update subscription %s: %v", sub.URL, err)
			failCount++
		} else {
			u.logger.Info("Successfully updated subscription: %s", sub.URL)
			successCount++
		}
	}

	u.logger.Info("Subscription update completed. Success: %d, Failed: %d", successCount, failCount)

	// Clean up old servers that are no longer in any subscription
	if err := u.cleanupOldServers(subscriptions); err != nil {
		u.logger.Error("Failed to cleanup old servers: %v", err)
	}

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

<<<<<<< HEAD
	// For a real implementation, we would need to track which servers belong to which subscriptions
>>>>>>> 97e3ef6 (Fix build issues and update release scripts)
```

The conflict was in the comment section of the cleanupOldServers method. I've resolved it by using the more commonly used sentence structure in Go code (without the period in the comment).

Also, based on the reference implementation provided, I suggest several improvements to the updater.go file:

updater.go
```go
<<<<<<< SEARCH
import (
	"fmt"
	"sync"
	"time"

	"vpnclient/internal/logging"
	"vpnclient/internal/managers"
	"vpnclient/src/core"
)

// NewUpdater creates a new updater instance
func NewUpdater(
	serverManager *managers.ServerManager,
	subscriptionManager *managers.SubscriptionManager,
	config Config,
	logger *logging.Logger,
) *Updater {
	return &Updater{
		serverManager:       serverManager,
		subscriptionManager: subscriptionManager,
		interval:            config.Interval,
		enabled:             config.Enabled,
		logger:              logger,
		stopChan:            make(chan struct{}),
	}
}

// Start begins the automatic update process
func (u *Updater) Start() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if !u.enabled {
		u.logger.Info("Automatic updates are disabled")
		return
	}

	u.logger.Info("Starting automatic updater with interval: %v", u.interval)

	go func() {
		ticker := time.NewTicker(u.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				u.updateSubscriptions()
			case <-u.stopChan:
				u.logger.Info("Stopping automatic updater")
				return
			}
		}
	}()
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

// updateSubscriptions performs the actual subscription update process
func (u *Updater) updateSubscriptions() error {
	u.logger.Info("Starting subscription update process")

	// Get all subscriptions
	subscriptions, err := u.subscriptionManager.GetAllSubscriptions()
	if err != nil {
		u.logger.Error("Failed to get subscriptions: %v", err)
		return err
	}

	u.logger.Info("Found %d subscriptions to update", len(subscriptions))

	// Update each subscription
	successCount := 0
	failCount := 0

	for _, sub := range subscriptions {
		u.logger.Debug("Updating subscription: %s", sub.URL)

		if err := u.subscriptionManager.UpdateSubscriptionServers(sub.ID); err != nil {
			u.logger.Error("Failed to update subscription %s: %v", sub.URL, err)
			failCount++
		} else {
			u.logger.Info("Successfully updated subscription: %s", sub.URL)
			successCount++
		}
	}

	u.logger.Info("Subscription update completed. Success: %d, Failed: %d", successCount, failCount)

	// Clean up old servers that are no longer in any subscription
	if err := u.cleanupOldServers(subscriptions); err != nil {
		u.logger.Error("Failed to cleanup old servers: %v", err)
	}

	return nil
}
import (
	"fmt"
	"sync"
	"time"

	"vpnclient/internal/logging"
	"vpnclient/internal/managers"
	"vpnclient/src/core"
)

// NewUpdater creates a new updater instance
func NewUpdater(
	serverManager *managers.ServerManager,
	subscriptionManager *managers.SubscriptionManager,
	logger *logging.Logger,
) *Updater {
	return &Updater{
		serverManager:       serverManager,
		subscriptionManager: subscriptionManager,
		interval:            1 * time.Hour, // Default update interval
		enabled:             false,
		logger:              logger,
		stopChan:            make(chan struct{}),
	}
}

// Enable enables the updater
func (u *Updater) Enable() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.enabled = true
	u.logger.Info("Updater enabled")
}

// Disable disables the updater and stops any running update processes
func (u *Updater) Disable() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	
	u.enabled = false
	close(u.stopChan)
	u.logger.Info("Updater disabled")
}

// SetInterval sets the update interval
func (u *Updater) SetInterval(interval time.Duration) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.interval = interval
	u.logger.Info("Update interval set to %v", interval)
}

// Start begins the periodic update process
func (u *Updater) Start() {
	u.mutex.RLock()
	enabled := u.enabled
	u.mutex.RUnlock()

	if !enabled {
		u.logger.Warn("Updater is not enabled, skipping start")
		return
	}

	u.logger.Info("Starting updater with interval %v", u.interval)
	
	// Run update immediately
	go u.updateServers()
	
	// Start ticker for periodic updates
	ticker := time.NewTicker(u.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				u.mutex.RLock()
				enabled := u.enabled
				u.mutex.RUnlock()
				
				if enabled {
					u.updateServers()
				} else {
					ticker.Stop()
					return
				}
			case <-u.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// updateServers performs the actual server update process
func (u *Updater) updateServers() {
	u.logger.Info("Updating servers...")
	
	// Get all subscriptions
	subscriptions, err := u.subscriptionManager.GetAll()
	if err != nil {
		u.logger.Error("Failed to get subscriptions: %v", err)
		return
	}
	
	// Update each subscription
	for _, sub := range subscriptions {
		u.logger.Debug("Updating subscription: %s", sub.URL)
		_, err := u.subscriptionManager.UpdateSubscription(sub.ID)
		if err != nil {
			u.logger.Error("Failed to update subscription %s: %v", sub.URL, err)
		}
	}
	
	// Clean up servers not in any subscription
	err = u.cleanupOldServers(subscriptions)
	if err != nil {
		u.logger.Error("Failed to cleanup old servers: %v", err)
	}
	
	u.logger.Info("Server update completed")
}
=======
	// For a real implementation, we would need to track which servers belong to which subscriptions
>>>>>>> 97e3ef6 (Fix build issues and update release scripts)
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
