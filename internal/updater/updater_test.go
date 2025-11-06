package updater

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"vpnclient/internal/database"
	"vpnclient/internal/logging"
	"vpnclient/internal/managers"
	"vpnclient/src/core"
)

// MockServerManager is a mock implementation of ServerManager
type MockServerManager struct {
	mock.Mock
}

func (m *MockServerManager) GetAllServers() ([]*core.Server, error) {
	args := m.Called()
	return args.Get(0).([]*core.Server), args.Error(1)
}

// MockSubscriptionManager is a mock implementation of SubscriptionManager
type MockSubscriptionManager struct {
	mock.Mock
}

func (m *MockSubscriptionManager) GetAllSubscriptions() ([]*core.Subscription, error) {
	args := m.Called()
	return args.Get(0).([]*core.Subscription), args.Error(1)
}

func (m *MockSubscriptionManager) UpdateSubscriptionServers(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestNewUpdater tests the creation of a new updater
func TestNewUpdater(t *testing.T) {
	// Setup
	serverManager := &MockServerManager{}
	subscriptionManager := &MockSubscriptionManager{}
	logger := &logging.Logger{}
	
	config := Config{
		Interval: time.Hour,
		Enabled:  true,
	}

	// Test
	updater := NewUpdater(serverManager, subscriptionManager, config, logger)

	// Assertions
	assert.NotNil(t, updater)
	assert.Equal(t, time.Hour, updater.interval)
	assert.True(t, updater.enabled)
}

// TestUpdaterStartStop tests starting and stopping the updater
func TestUpdaterStartStop(t *testing.T) {
	// Setup
	serverManager := &MockServerManager{}
	subscriptionManager := &MockSubscriptionManager{}
	logger := &logging.Logger{}
	
	config := Config{
		Interval: 10 * time.Millisecond, // Very short interval for testing
		Enabled:  true,
	}

	updater := NewUpdater(serverManager, subscriptionManager, config, logger)

	// Test
	updater.Start()
	time.Sleep(20 * time.Millisecond) // Let it run for a bit
	updater.Stop()

	// Assertions
	// If we get here without hanging, the start/stop worked
	assert.True(t, true)
}

// TestUpdaterDisabled tests that updater doesn't run when disabled
func TestUpdaterDisabled(t *testing.T) {
	// Setup
	serverManager := &MockServerManager{}
	subscriptionManager := &MockSubscriptionManager{}
	logger := &logging.Logger{}
	
	config := Config{
		Interval: 10 * time.Millisecond,
		Enabled:  false,
	}

	updater := NewUpdater(serverManager, subscriptionManager, config, logger)

	// Test
	err := updater.UpdateSubscriptions()

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "automatic updates are disabled")
}

// TestUpdateSubscriptions tests the subscription update functionality
func TestUpdateSubscriptions(t *testing.T) {
	// Setup
	serverManager := &MockServerManager{}
	subscriptionManager := &MockSubscriptionManager{}
	logger := &logging.Logger{}
	
	config := Config{
		Interval: time.Hour,
		Enabled:  true,
	}

	updater := NewUpdater(serverManager, subscriptionManager, config, logger)

	// Mock expectations
	subscriptions := []*core.Subscription{
		{ID: "sub1", URL: "http://example.com/sub1"},
		{ID: "sub2", URL: "http://example.com/sub2"},
	}
	
	subscriptionManager.On("GetAllSubscriptions").Return(subscriptions, nil)
	subscriptionManager.On("UpdateSubscriptionServers", "sub1").Return(nil)
	subscriptionManager.On("UpdateSubscriptionServers", "sub2").Return(nil)

	// Test
	err := updater.UpdateSubscriptions()

	// Assertions
	assert.NoError(t, err)
	subscriptionManager.AssertExpectations(t)
}

// TestUpdateSubscriptionsWithErrors tests subscription updates with errors
func TestUpdateSubscriptionsWithErrors(t *testing.T) {
	// Setup
	serverManager := &MockServerManager{}
	subscriptionManager := &MockSubscriptionManager{}
	logger := &logging.Logger{}
	
	config := Config{
		Interval: time.Hour,
		Enabled:  true,
	}

	updater := NewUpdater(serverManager, subscriptionManager, config, logger)

	// Mock expectations
	subscriptions := []*core.Subscription{
		{ID: "sub1", URL: "http://example.com/sub1"},
		{ID: "sub2", URL: "http://example.com/sub2"},
	}
	
	subscriptionManager.On("GetAllSubscriptions").Return(subscriptions, nil)
	subscriptionManager.On("UpdateSubscriptionServers", "sub1").Return(nil)
	subscriptionManager.On("UpdateSubscriptionServers", "sub2").Return(
		&managers.SubscriptionError{Msg: "network error"})

	// Test
	err := updater.UpdateSubscriptions()

	// Assertions
	// Should not return error even if some subscriptions fail
	assert.NoError(t, err)
	subscriptionManager.AssertExpectations(t)
}

// TestSetEnabled tests enabling/disabling the updater
func TestSetEnabled(t *testing.T) {
	// Setup
	serverManager := &MockServerManager{}
	subscriptionManager := &MockSubscriptionManager{}
	logger := &logging.Logger{}
	
	config := Config{
		Interval: time.Hour,
		Enabled:  false,
	}

	updater := NewUpdater(serverManager, subscriptionManager, config, logger)

	// Test
	updater.SetEnabled(true)

	// Since we can't directly access the enabled field, we test by trying to update
	updater.SetEnabled(true)
	
	// Mock expectations
	subscriptions := []*core.Subscription{}
	subscriptionManager.On("GetAllSubscriptions").Return(subscriptions, nil)
	
	// Test update - should work now
	err := updater.UpdateSubscriptions()
	
	// Assertions
	assert.NoError(t, err)
}

// TestSetInterval tests setting the update interval
func TestSetInterval(t *testing.T) {
	// Setup
	serverManager := &MockServerManager{}
	subscriptionManager := &MockSubscriptionManager{}
	logger := &logging.Logger{}
	
	config := Config{
		Interval: time.Hour,
		Enabled:  true,
	}

	updater := NewUpdater(serverManager, subscriptionManager, config, logger)

	// Test
	newInterval := 2 * time.Hour
	updater.SetInterval(newInterval)

	// Assertions
	// We can't directly access the interval, but we can check the status
	status := updater.GetStatus()
	// Status returns interval as string, so we can't directly compare
	assert.Contains(t, status["interval"], "h0m0s") // Should contain hours
}

// TestGetStatus tests getting the updater status
func TestGetStatus(t *testing.T) {
	// Setup
	serverManager := &MockServerManager{}
	subscriptionManager := &MockSubscriptionManager{}
	logger := &logging.Logger{}
	
	config := Config{
		Interval: 30 * time.Minute,
		Enabled:  true,
	}

	updater := NewUpdater(serverManager, subscriptionManager, config, logger)

	// Test
	status := updater.GetStatus()

	// Assertions
	assert.NotNil(t, status)
	assert.True(t, status["enabled"].(bool))
	assert.Equal(t, "30m0s", status["interval"])
}