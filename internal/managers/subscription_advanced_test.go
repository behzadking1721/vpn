package managers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"vpnclient/internal/database"
	"vpnclient/internal/logging"
	"vpnclient/internal/notifications"
	"vpnclient/mocks"
	"vpnclient/src/core"
)

// TestSubscriptionManagerWithNotifications tests subscription manager with notifications integration
func TestSubscriptionManagerWithNotifications(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	serverManager := NewServerManager(mockStore)
	sm := NewSubscriptionManager(serverManager, mockStore)
	notificationManager := notifications.NewNotificationManager(10)
	sm.SetNotificationManager(notificationManager)

	// Mock expectations for successful subscription
	mockStore.On("Add", mock.Anything, mock.AnythingOfType("*core.Subscription")).Return(nil)
	mockStore.On("GetAll", mock.Anything).Return([]interface{}{}, nil)

	// Test adding subscription
	sub, err := sm.AddSubscription("http://example.com/subscription")
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, "http://example.com/subscription", sub.URL)

	// Check notifications
	notifications := notificationManager.GetNotifications()
	assert.NotEmpty(t, notifications)

	// Look for subscription added notification
	found := false
	for _, notif := range notifications {
		if notif.Title == "Subscription Added" {
			found = true
			assert.Contains(t, notif.Message, "Successfully added subscription")
			assert.Equal(t, notifications.Success, notif.Type)
			break
		}
	}
	assert.True(t, found)
}

// TestSubscriptionManagerWithLogger tests subscription manager with logger integration
func TestSubscriptionManagerWithLogger(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	serverManager := NewServerManager(mockStore)
	sm := NewSubscriptionManager(serverManager, mockStore)
	logger := &logging.Logger{}
	sm.SetLogger(logger)

	// Mock expectations
	subscription := &core.Subscription{
		ID:          "test-sub",
		URL:         "http://example.com/subscription",
		ServerCount: 0,
	}
	
	mockStore.On("Add", mock.Anything, mock.AnythingOfType("*core.Subscription")).Return(nil)
	mockStore.On("Get", mock.Anything, "test-sub").Return(subscription, nil)
	mockStore.On("Update", mock.Anything, mock.AnythingOfType("*core.Subscription")).Return(nil)
	mockStore.On("Delete", mock.Anything, "test-sub").Return(nil)
	mockStore.On("GetAll", mock.Anything).Return([]interface{}{subscription}, nil)

	// Test adding subscription
	sub, err := sm.AddSubscription("http://example.com/subscription")
	assert.NoError(t, err)
	assert.NotNil(t, sub)

	// Test getting subscription
	retrievedSub, err := sm.GetSubscription("test-sub")
	assert.NoError(t, err)
	assert.Equal(t, subscription, retrievedSub)

	// Test updating subscription
	subscription.ServerCount = 5
	err = sm.UpdateSubscription(subscription)
	assert.NoError(t, err)

	// Test getting all subscriptions
	allSubs, err := sm.GetAllSubscriptions()
	assert.NoError(t, err)
	assert.Len(t, allSubs, 1)
	assert.Equal(t, subscription, allSubs[0])

	// Test deleting subscription
	err = sm.DeleteSubscription("test-sub")
	assert.NoError(t, err)
}

// TestSubscriptionManagerUpdateSubscriptionServers tests updating servers from subscription
func TestSubscriptionManagerUpdateSubscriptionServers(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	serverManager := NewServerManager(mockStore)
	sm := NewSubscriptionManager(serverManager, mockStore)

	// Mock expectations
	subscription := &core.Subscription{
		ID:          "test-sub",
		URL:         "http://example.com/subscription",
		ServerCount: 0,
	}
	
	// For subscription store
	mockStore.On("Get", mock.Anything, "test-sub").Return(subscription, nil)
	mockStore.On("Update", mock.Anything, mock.AnythingOfType("*core.Subscription")).Return(nil)

	// Test updating subscription servers
	err := sm.UpdateSubscriptionServers("test-sub")
	assert.NoError(t, err)
}

// TestSubscriptionManagerUpdateSubscriptionServersError tests updating servers from subscription with error
func TestSubscriptionManagerUpdateSubscriptionServersError(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	serverManager := NewServerManager(mockStore)
	sm := NewSubscriptionManager(serverManager, mockStore)

	// Mock expectations for non-existent subscription
	mockStore.On("Get", mock.Anything, "non-existent").Return((*core.Subscription)(nil), database.ErrNotFound)

	// Test updating subscription servers with non-existent subscription
	err := sm.UpdateSubscriptionServers("non-existent")
	assert.Error(t, err)
	assert.Equal(t, database.ErrNotFound, err)
}

// TestSubscriptionManagerAddSubscriptionError tests adding subscription with error
func TestSubscriptionManagerAddSubscriptionError(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	serverManager := NewServerManager(mockStore)
	sm := NewSubscriptionManager(serverManager, mockStore)

	// Mock expectations for store error
	mockStore.On("Add", mock.Anything, mock.AnythingOfType("*core.Subscription")).Return(assert.AnError)

	// Test adding subscription with store error
	sub, err := sm.AddSubscription("http://example.com/subscription")
	assert.Error(t, err)
	assert.Nil(t, sub)
	assert.Equal(t, assert.AnError, err)
}

// TestSubscriptionManagerGetAllSubscriptionsError tests getting all subscriptions with error
func TestSubscriptionManagerGetAllSubscriptionsError(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	serverManager := NewServerManager(mockStore)
	sm := NewSubscriptionManager(serverManager, mockStore)

	// Mock expectations for store error
	mockStore.On("GetAll", mock.Anything).Return([]interface{}(nil), assert.AnError)

	// Test getting all subscriptions with store error
	subs, err := sm.GetAllSubscriptions()
	assert.Error(t, err)
	assert.Nil(t, subs)
	assert.Equal(t, assert.AnError, err)
}

// TestSubscriptionManagerGetSubscriptionError tests getting subscription with error
func TestSubscriptionManagerGetSubscriptionError(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	serverManager := NewServerManager(mockStore)
	sm := NewSubscriptionManager(serverManager, mockStore)

	// Mock expectations for store error
	mockStore.On("Get", mock.Anything, "test-sub").Return((*core.Subscription)(nil), assert.AnError)

	// Test getting subscription with store error
	sub, err := sm.GetSubscription("test-sub")
	assert.Error(t, err)
	assert.Nil(t, sub)
	assert.Equal(t, assert.AnError, err)
}

// TestSubscriptionManagerDeleteSubscriptionError tests deleting subscription with error
func TestSubscriptionManagerDeleteSubscriptionError(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	serverManager := NewServerManager(mockStore)
	sm := NewSubscriptionManager(serverManager, mockStore)

	// Mock expectations for store error
	mockStore.On("Get", mock.Anything, "test-sub").Return(&core.Subscription{ID: "test-sub"}, nil)
	mockStore.On("Delete", mock.Anything, "test-sub").Return(assert.AnError)

	// Test deleting subscription with store error
	err := sm.DeleteSubscription("test-sub")
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}